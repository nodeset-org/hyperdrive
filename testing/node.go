package testing

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/nodeset-org/hyperdrive-daemon/client"
	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/server"
	"github.com/nodeset-org/hyperdrive-daemon/shared/auth"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
)

const (
	apiAuthKey string = "test-key"
)

// A complete Hyperdrive node instance
type HyperdriveNode struct {
	// The daemon's service provider
	sp common.IHyperdriveServiceProvider

	// The daemon's HTTP API server
	serverMgr *server.ServerManager

	// An HTTP API client for the daemon
	client *client.ApiClient

	// The client logger
	logger *slog.Logger

	// Wait group for graceful shutdown
	wg *sync.WaitGroup
}

// Create a new Hyperdrive node, including its folder structure, service provider, server manager, and API client.
func newHyperdriveNode(sp common.IHyperdriveServiceProvider, address string, clientLogger *slog.Logger) (*HyperdriveNode, error) {
	// Create the server
	wg := &sync.WaitGroup{}
	cfg := sp.GetConfig()
	authMgr := auth.NewAuthorizationManager("")
	authMgr.SetKey([]byte(apiAuthKey))
	serverMgr, err := server.NewServerManager(sp, address, cfg.ApiPort.Value, wg, authMgr)
	if err != nil {
		return nil, fmt.Errorf("error creating hyperdrive server: %v", err)
	}

	// Create the client
	urlString := fmt.Sprintf("http://%s:%d/%s", address, serverMgr.GetPort(), config.HyperdriveApiClientRoute)
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("error parsing client URL [%s]: %v", urlString, err)
	}
	apiClient := client.NewApiClient(url, clientLogger, nil, authMgr)

	return &HyperdriveNode{
		sp:        sp,
		serverMgr: serverMgr,
		client:    apiClient,
		logger:    clientLogger,
		wg:        wg,
	}, nil
}

// Closes the Constellation node and its Hyperdrive parent.
func (n *HyperdriveNode) Close() error {
	if n.serverMgr != nil {
		n.serverMgr.Stop()
		n.wg.Wait()
		n.serverMgr = nil
		n.logger.Info("Stopped Hyperdrive daemon API server")
	}
	return nil
}

// Get the daemon's service provider
func (n *HyperdriveNode) GetServiceProvider() common.IHyperdriveServiceProvider {
	return n.sp
}

// Get the HTTP API server for the node's daemon
func (n *HyperdriveNode) GetServerManager() *server.ServerManager {
	return n.serverMgr
}

// Get the HTTP API client for interacting with the node's daemon server
func (n *HyperdriveNode) GetApiClient() *client.ApiClient {
	return n.client
}

// Create a new Hyperdrive node based on this one's configuration, but with a custom folder, address, and port.
func (n *HyperdriveNode) CreateSubNode(folder string, address string, port uint16) (*HyperdriveNode, error) {
	// Make a new config
	parentSp := n.sp
	parentHdCfg := parentSp.GetConfig()
	hdNetSettings := parentHdCfg.GetNetworkSettings()
	cfg, err := hdconfig.NewHyperdriveConfigForNetwork(folder, hdNetSettings, parentHdCfg.Network.Value)
	if err != nil {
		return nil, fmt.Errorf("error creating Hyperdrive config: %v", err)
	}
	cfg.UserDataPath.Value = filepath.Join(folder, "data")
	cfg.ApiPort.Value = port

	// Make sure the data and modules directories exist
	dataDir := cfg.UserDataPath.Value
	moduleDir := filepath.Join(dataDir, hdconfig.ModulesName)
	err = os.MkdirAll(moduleDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating data and modules directories [%s]: %v", moduleDir, err)
	}

	// Make a new service provider
	sp, err := common.NewHyperdriveServiceProviderFromCustomServices(
		cfg,
		parentSp.GetResources(),
		parentSp.GetEthClient(),
		parentSp.GetBeaconClient(),
		parentSp.GetDocker(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Hyperdrive service provider: %v", err)
	}

	return newHyperdriveNode(sp, address, n.logger)
}

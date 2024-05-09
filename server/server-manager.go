package server

import (
	"fmt"
	"sync"

	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/service"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/tx"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/utils"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/wallet"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/rocket-pool/node-manager-core/api/server"
)

// ServerManager manages the API server run by the daemon
type ServerManager struct {
	// The server for clients to interact with
	apiServer *server.NetworkSocketApiServer
}

// Creates a new server manager
func NewServerManager(sp *common.ServiceProvider, ip string, port uint16, stopWg *sync.WaitGroup) (*ServerManager, error) {
	// Start the API server
	apiServer, err := createServer(sp, ip, port)
	if err != nil {
		return nil, fmt.Errorf("error creating API server: %w", err)
	}
	err = apiServer.Start(stopWg)
	if err != nil {
		return nil, fmt.Errorf("error starting API server: %w", err)
	}
	fmt.Printf("API server started on %s:%d\n", ip, port)

	// Create the manager
	mgr := &ServerManager{
		apiServer: apiServer,
	}
	return mgr, nil
}

// Stops and shuts down the servers
func (m *ServerManager) Stop() {
	err := m.apiServer.Stop()
	if err != nil {
		fmt.Printf("WARNING: API server didn't shutdown cleanly: %s\n", err.Error())
	}
}

// Creates a new Hyperdrive API server
func createServer(sp *common.ServiceProvider, ip string, port uint16) (*server.NetworkSocketApiServer, error) {
	apiLogger := sp.GetApiLogger()
	ctx := apiLogger.CreateContextWithLogger(sp.GetBaseContext())

	handlers := []server.IHandler{
		service.NewServiceHandler(apiLogger, ctx, sp),
		tx.NewTxHandler(apiLogger, ctx, sp),
		utils.NewUtilsHandler(apiLogger, ctx, sp),
		wallet.NewWalletHandler(apiLogger, ctx, sp),
	}

	server, err := server.NewNetworkSocketApiServer(apiLogger.Logger, ip, port, handlers, config.HyperdriveDaemonRoute, config.HyperdriveApiVersion)
	if err != nil {
		return nil, err
	}
	return server, nil
}

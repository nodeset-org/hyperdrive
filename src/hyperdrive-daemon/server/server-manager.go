package server

import (
	"fmt"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/service"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/tx"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/wallet"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/api/server"
)

const (
	cliOrigin string = "cli"
	webOrigin string = "net"
)

// ServerManager manages all of the daemon sockets and servers run by the main Hyperdrive daemon
type ServerManager struct {
	// The server for the CLI to interact with
	cliServer *server.ApiServer

	// The server for the Stakewise module
	stakewiseServer *server.ApiServer

	// The daemon's main closing waitgroup
	stopWg *sync.WaitGroup
}

// Creates a new server manager
func NewServerManager(sp *common.ServiceProvider, cfgPath string, stopWg *sync.WaitGroup, moduleNames []string) (*ServerManager, error) {
	mgr := &ServerManager{
		stopWg: stopWg,
	}

	// Get the owner of the config file
	var cfgFileStat syscall.Stat_t
	err := syscall.Stat(cfgPath, &cfgFileStat)
	if err != nil {
		return nil, fmt.Errorf("error getting config file [%s] info: %w", cfgPath, err)
	}

	// Start the CLI server
	cliSocketPath := filepath.Join(sp.GetUserDir(), config.HyperdriveCliSocketFilename)
	cliServer, err := createServer(cliOrigin, sp, cliSocketPath)
	if err != nil {
		return nil, fmt.Errorf("error creating CLI server: %w", err)
	}
	err = cliServer.Start(stopWg, cfgFileStat.Uid, cfgFileStat.Gid)
	if err != nil {
		return nil, fmt.Errorf("error starting CLI server: %w", err)
	}
	mgr.cliServer = cliServer
	fmt.Printf("CLI daemon started on %s\n", cliSocketPath)

	// Handle each module server
	for _, module := range moduleNames {
		modulesDir := filepath.Join(sp.GetConfig().UserDataPath.Value, config.ModulesName)
		moduleSocketPath := filepath.Join(modulesDir, module, config.HyperdriveCliSocketFilename)
		server, err := createServer(module, sp, moduleSocketPath)
		if err != nil {
			return nil, fmt.Errorf("error creating server for module [%s]: %w", module, err)
		}
		err = server.Start(stopWg, cfgFileStat.Uid, cfgFileStat.Gid)
		if err != nil {
			return nil, fmt.Errorf("error starting server for module [%s]: %w", module, err)
		}
		mgr.stakewiseServer = server
		fmt.Printf("Daemon started on %s\n", moduleSocketPath)
	}

	return mgr, nil
}

// Stops and shuts down the servers
func (m *ServerManager) Stop() {
	err := m.cliServer.Stop()
	if err != nil {
		fmt.Printf("WARNING: CLI server didn't shutdown cleanly: %s\n", err.Error())
		m.stopWg.Done()
	}

	if m.stakewiseServer != nil {
		err := m.stakewiseServer.Stop()
		if err != nil {
			fmt.Printf("WARNING: Stakewise server didn't shutdown cleanly: %s\n", err.Error())
			m.stopWg.Done()
		}
	}
}

// Creates a new Hyperdrive API server
func createServer(origin string, sp *common.ServiceProvider, socketPath string) (*server.ApiServer, error) {
	apiLogger := sp.GetApiLogger()
	subLogger := apiLogger.CreateSubLogger(origin)
	ctx := subLogger.CreateContextWithLogger(sp.GetBaseContext())

	handlers := []server.IHandler{
		service.NewServiceHandler(subLogger, ctx, sp),
		tx.NewTxHandler(subLogger, ctx, sp),
		utils.NewUtilsHandler(subLogger, ctx, sp),
		wallet.NewWalletHandler(subLogger, ctx, sp),
	}

	server, err := server.NewApiServer(subLogger.Logger, socketPath, handlers, config.HyperdriveDaemonRoute, config.HyperdriveApiVersion)
	if err != nil {
		return nil, err
	}
	return server, nil
}

package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/nodeset"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/service"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/tx"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/utils"
	"github.com/nodeset-org/hyperdrive-daemon/server/api/wallet"
	"github.com/nodeset-org/hyperdrive-daemon/shared/auth"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/log"
)

// ServerManager manages the API server run by the daemon
type ServerManager struct {
	// The server for clients to interact with
	apiServer *server.NetworkSocketApiServer
}

// Creates a new server manager
func NewServerManager(sp common.IHyperdriveServiceProvider, ip string, port uint16, stopWg *sync.WaitGroup, authMgr *auth.AuthorizationManager) (*ServerManager, error) {
	// Start the API server
	apiServer, err := createServer(sp, ip, port, authMgr)
	if err != nil {
		return nil, fmt.Errorf("error creating API server: %w", err)
	}
	err = apiServer.Start(stopWg)
	if err != nil {
		return nil, fmt.Errorf("error starting API server: %w", err)
	}
	port = apiServer.GetPort()
	fmt.Printf("API server started on %s:%d\n", ip, port)

	// Create the manager
	mgr := &ServerManager{
		apiServer: apiServer,
	}
	return mgr, nil
}

// Returns the port the server is running on
func (m *ServerManager) GetPort() uint16 {
	return m.apiServer.GetPort()
}

// Stops and shuts down the servers
func (m *ServerManager) Stop() {
	err := m.apiServer.Stop()
	if err != nil {
		fmt.Printf("WARNING: API server didn't shutdown cleanly: %s\n", err.Error())
	}
}

// Creates a new Hyperdrive API server
func createServer(sp common.IHyperdriveServiceProvider, ip string, port uint16, authMgr *auth.AuthorizationManager) (*server.NetworkSocketApiServer, error) {
	apiLogger := sp.GetApiLogger()
	ctx := apiLogger.CreateContextWithLogger(sp.GetBaseContext())

	// Create the API handlers
	handlers := []server.IHandler{
		nodeset.NewNodeSetHandler(apiLogger, ctx, sp),
		service.NewServiceHandler(apiLogger, ctx, sp),
		tx.NewTxHandler(apiLogger, ctx, sp),
		utils.NewUtilsHandler(apiLogger, ctx, sp),
		wallet.NewWalletHandler(apiLogger, ctx, sp),
	}

	// Create the API server
	server, err := server.NewNetworkSocketApiServer(apiLogger.Logger, ip, port, handlers, config.HyperdriveDaemonRoute, config.HyperdriveApiVersion)
	if err != nil {
		return nil, err
	}

	// Add the authorization middleware
	server.GetApiRouter().Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err = authMgr.ValidateRequest(r)
			if err != nil {
				apiLogger.Error("Error validating request authorization",
					log.Err(err),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				)
				http.Error(w, "Authorization failed", http.StatusUnauthorized)
				return
			}

			// Valid request
			next.ServeHTTP(w, r)
		})
	})
	return server, nil
}

package api

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/services"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

const (
	ApiLogColor color.Attribute = color.FgHiBlue
)

type IHandler interface {
	RegisterRoutes(router *mux.Router, debugMode bool)
}

type ApiManager struct {
	debugMode  bool
	log        log.ColorLogger
	handlers   []IHandler
	socketPath string
	socket     net.Listener
	server     http.Server
	router     *mux.Router
}

// parameter: sp *services.ServiceProvider
func NewApiManager(sp *services.ServiceProvider) *ApiManager {
	// Create the router
	router := mux.NewRouter()

	// Create the manager
	cfg := sp.GetConfig()
	mgr := &ApiManager{
		debugMode: cfg.DebugMode.Value,
		log:       log.NewColorLogger(ApiLogColor),
		handlers:  []IHandler{
			// example.NewNodeHandler(sp, debugMode),
		},
		socketPath: config.DaemonSocketPath,
		router:     router,
		server: http.Server{
			Handler: router,
		},
	}

	if mgr.debugMode {
		fmt.Println("Debug mode active; printing commands without execution.")
	}

	// Register each route
	hyperdriveRouter := router.Host("hyperdrive").Subrouter() // The host will be accessible at http://hyperdrive/...
	for _, handler := range mgr.handlers {
		handler.RegisterRoutes(hyperdriveRouter, mgr.debugMode)
	}

	return mgr
}

// Starts listening for incoming HTTP requests
func (m *ApiManager) Start(wg *sync.WaitGroup, socketOwnerUid uint32, socketOwnerGid uint32) error {
	// Remove the socket if it's already there
	_, err := os.Stat(m.socketPath)
	if !errors.Is(err, fs.ErrNotExist) {
		err = os.Remove(m.socketPath)
		if err != nil {
			return fmt.Errorf("error removing socket file: %w", err)
		}
	}

	// Create the socket
	socket, err := net.Listen("unix", m.socketPath)
	if err != nil {
		return fmt.Errorf("error creating socket: %w", err)
	}
	m.socket = socket

	// Set the socket owner to the config file user
	err = os.Chown(m.socketPath, int(socketOwnerUid), int(socketOwnerGid))
	if err != nil {
		return fmt.Errorf("error setting socket owner: %w", err)
	}

	// Make it so only the user can write to the socket
	err = os.Chmod(m.socketPath, 0600)
	if err != nil {
		return fmt.Errorf("error relaxing permissions on socket: %w", err)
	}

	// Start listening
	go func() {
		err := m.server.Serve(socket)
		if !errors.Is(err, http.ErrServerClosed) {
			m.log.Printlnf("error while listening for HTTP requests: %s", err.Error())
		}
		wg.Done()
	}()

	return nil
}

// Stops the HTTP listener
func (m *ApiManager) Stop() error {
	err := m.server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("error stopping listener: %w", err)
	}
	return nil
}

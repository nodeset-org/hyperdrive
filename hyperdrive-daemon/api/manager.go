package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	example "github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/node"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/services"
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
func NewApiManager(sp *services.ServiceProvider, socketPath string, debugMode bool) *ApiManager {
	// Create the router
	router := mux.NewRouter()

	// Create the manager
	// cfg := sp.GetConfig()
	mgr := &ApiManager{
		debugMode: debugMode,
		log:       log.NewColorLogger(ApiLogColor),
		handlers: []IHandler{
			example.NewNodeHandler(sp, debugMode),
		},
		socketPath: socketPath,
		router:     router,
		server: http.Server{
			Handler: router,
		},
	}

	if debugMode {
		fmt.Println("Debug mode active; printing commands without execution.")
	}

	// Register each route
	hyperdriveRouter := router.Host("hyperdrive").Subrouter() // The host will be accessible at http://hyperdrive/...
	for _, handler := range mgr.handlers {
		handler.RegisterRoutes(hyperdriveRouter, debugMode)
	}

	return mgr
}

// Starts listening for incoming HTTP requests
func (m *ApiManager) Start(wg *sync.WaitGroup) error {
	// Create the socket
	socket, err := net.Listen("unix", m.socketPath)
	if err != nil {
		return fmt.Errorf("error creating socket: %w", err)
	}
	m.socket = socket

	// Make it so anyone can write to the socket
	err = os.Chmod(m.socketPath, 0766)
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

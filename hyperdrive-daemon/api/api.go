package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/log"
)

const (
	ApiLogColor color.Attribute = color.FgHiBlue
)

type IHandler interface {
	RegisterRoutes(router *mux.Router)
}

type ApiManager struct {
	log        log.ColorLogger
	handlers   []IHandler
	socketPath string
	socket     net.Listener
	server     http.Server
	router     *mux.Router
}

// parameter: sp *services.ServiceProvider
func NewApiManager() *ApiManager {
	// Create the router
	router := mux.NewRouter()

	// Create the manager
	// cfg := sp.GetConfig()
	mgr := &ApiManager{
		log: log.NewColorLogger(ApiLogColor),
		// handlers: []IHandler{
		// 	example.NewExampleHandler(sp),
		// },
		socketPath: "some-path-for-now", //cfg.Smartnode.GetSocketPath(),
		router:     router,
		server: http.Server{
			Handler: router,
		},
	}

	// Register each route
	for _, handler := range mgr.handlers {
		handler.RegisterRoutes(mgr.router)
	}

	return mgr
}

// Starts listening for incoming HTTP requests
func (m *ApiManager) Start() error {
	// Create the socket
	socket, err := net.Listen("unix", m.socketPath)
	if err != nil {
		return fmt.Errorf("error creating socket: %w", err)
	}
	m.socket = socket

	// Start listening
	go func() {
		err := m.server.Serve(socket)
		if !errors.Is(err, http.ErrServerClosed) {
			m.log.Printlnf("error while listening for HTTP requests: %s", err.Error())
		}
	}()

	return nil
}

// Stops the HTTP listener
func (m *ApiManager) Stop() error {
	// Shutdown the listener
	err := m.server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("error stopping listener: %w", err)
	}

	// Remove the socket file
	err = os.Remove(m.socketPath)
	if err != nil {
		return fmt.Errorf("error removing socket file: %w", err)
	}

	return nil
}

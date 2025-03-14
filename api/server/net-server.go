package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/logging"
)

type NetworkSocketApiServer struct {
	logger    *slog.Logger
	handlers  []IHandler
	ip        string
	port      uint16
	socket    net.Listener
	server    http.Server
	router    *mux.Router
	apiRouter *mux.Router
}

func NewNetworkSocketApiServer(logger *slog.Logger, ip string, port uint16, handlers []IHandler, baseRoute string, apiVersion string) (*NetworkSocketApiServer, error) {
	// Create the router
	router := mux.NewRouter()

	// Create the manager
	server := &NetworkSocketApiServer{
		logger:   logger,
		handlers: handlers,
		ip:       ip,
		port:     port,
		router:   router,
		server: http.Server{
			Handler: router,
		},
	}

	// Register each route
	nmcRouter := router.PathPrefix("/" + baseRoute + "/api/v" + apiVersion).Subrouter()
	for _, handler := range server.handlers {
		handler.RegisterRoutes(nmcRouter)
	}
	server.apiRouter = nmcRouter

	return server, nil
}

// Starts listening for incoming HTTP requests
func (s *NetworkSocketApiServer) Start(wg *sync.WaitGroup) error {
	// Create the socket
	socket, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		return fmt.Errorf("error creating socket: %w", err)
	}
	s.socket = socket

	// Get the port if random
	if s.port == 0 {
		s.port = uint16(socket.Addr().(*net.TCPAddr).Port)
	}

	// Start listening
	wg.Add(1)
	go func() {
		err := s.server.Serve(socket)
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("error while listening for HTTP requests", logging.Err(err))
		}
		wg.Done()
	}()

	return nil
}

// Stops the HTTP listener
func (s *NetworkSocketApiServer) Stop() error {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("error stopping listener: %w", err)
	}
	return nil
}

// Get the port the server is running on - useful if the port was automatically assigned
func (s *NetworkSocketApiServer) GetPort() uint16 {
	return s.port
}

// Get the API router for the server
func (s *NetworkSocketApiServer) GetApiRouter() *mux.Router {
	return s.apiRouter
}

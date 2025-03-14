package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
)

type UnixSocketApiServer struct {
	logger     *slog.Logger
	handlers   []IHandler
	socketPath string
	socket     net.Listener
	server     http.Server
	router     *mux.Router
}

func NewUnixSocketApiServer(logger *slog.Logger, socketPath string, handlers []IHandler, baseRoute string, apiVersion string) (*UnixSocketApiServer, error) {
	// Create the router
	router := mux.NewRouter()

	// Create the manager
	server := &UnixSocketApiServer{
		logger:     logger,
		handlers:   handlers,
		socketPath: socketPath,
		router:     router,
		server: http.Server{
			Handler: router,
		},
	}

	// Register each route
	nmcRouter := router.Host(baseRoute).PathPrefix("/api/v" + apiVersion).Subrouter()
	for _, handler := range server.handlers {
		handler.RegisterRoutes(nmcRouter)
	}

	// Create the socket directory
	socketDir := filepath.Dir(socketPath)
	err := os.MkdirAll(socketDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("error creating socket directory [%s]: %w", socketDir, err)
	}

	return server, nil
}

// Starts listening for incoming HTTP requests
func (s *UnixSocketApiServer) Start(wg *sync.WaitGroup, socketOwnerUid uint32, socketOwnerGid uint32) error {
	// Remove the socket if it's already there
	_, err := os.Stat(s.socketPath)
	if !errors.Is(err, fs.ErrNotExist) {
		err = os.Remove(s.socketPath)
		if err != nil {
			return fmt.Errorf("error removing socket file: %w", err)
		}
	}

	// Create the socket
	socket, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("error creating socket: %w", err)
	}
	s.socket = socket

	// Make it so only the user can write to the socket
	err = os.Chmod(s.socketPath, 0600)
	if err != nil {
		return fmt.Errorf("error setting permissions on socket: %w", err)
	}

	// Set the socket owner to the config file user
	err = os.Chown(s.socketPath, int(socketOwnerUid), int(socketOwnerGid))
	if err != nil {
		return fmt.Errorf("error setting socket owner: %w", err)
	}

	// Start listening
	wg.Add(1)
	go func() {
		err := s.server.Serve(socket)
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("error while listening for HTTP requests", "error", err)
		}
		wg.Done()
	}()

	return nil
}

// Stops the HTTP listener
func (s *UnixSocketApiServer) Stop() error {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("error stopping listener: %w", err)
	}
	return nil
}

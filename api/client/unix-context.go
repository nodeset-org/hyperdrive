package client

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
)

// The context passed into a requester
type UnixRequesterContext struct {
	// The path to the socket to send requests to
	socketPath string

	// An HTTP client for sending requests
	client *http.Client

	// Logger to print debug messages to
	logger *slog.Logger

	// The base route for the client to send requests to (<http://<base>/<route>/<method>)
	base string
}

// Creates a new API client requester context
func NewUnixRequesterContext(baseRoute string, socketPath string, log *slog.Logger) *UnixRequesterContext {
	requesterContext := &UnixRequesterContext{
		socketPath: socketPath,
		base:       baseRoute,
		logger:     log,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
		},
	}

	return requesterContext
}

// Get the base of the address used for submitting server requests
func (r *UnixRequesterContext) GetAddressBase() string {
	return fmt.Sprintf("http://%s", r.base)
}

// Get the logger for the context
func (r *UnixRequesterContext) GetLogger() *slog.Logger {
	return r.logger
}

// Set the logger for the context
func (r *UnixRequesterContext) SetLogger(logger *slog.Logger) {
	r.logger = logger
}

// Send an HTTP request to the server
func (r *UnixRequesterContext) SendRequest(request *http.Request) (*http.Response, error) {
	// Make sure the socket exists
	_, err := os.Stat(r.socketPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("the socket at [%s] does not exist - please start the service and try again", r.socketPath)
	}

	return r.client.Do(request)
}

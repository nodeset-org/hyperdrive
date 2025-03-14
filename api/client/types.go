package client

import (
	"log/slog"
	"net/http"
)

// IRequester is an interface for making HTTP requests to a specific subroute on the NMC server
type IRequester interface {
	// The human-readable requester name (for logging)
	GetName() string

	// The name of the subroute to send requests to
	GetRoute() string

	// Context to hold settings and utilities the requester should use
	GetContext() IRequesterContext
}

// IRequester is an interface for making HTTP requests to a specific subroute on the NMC server
type IRequesterContext interface {
	// Get the base of the address used for submitting server requests
	GetAddressBase() string

	// Get the logger for the context
	GetLogger() *slog.Logger

	// Set the logger for the context
	SetLogger(*slog.Logger)

	// Send an HTTP request to the server
	SendRequest(request *http.Request) (*http.Response, error)
}

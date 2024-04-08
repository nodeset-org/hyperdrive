package constclient

import (
	"log/slog"

	"github.com/rocket-pool/node-manager-core/api/client"
)

// Binder for the Hyperdrive daemon API server
type ApiClient struct {
	context *client.RequesterContext
}

// Creates a new API client instance
func NewApiClient(baseRoute string, socketPath string, logger *slog.Logger) *ApiClient {
	context := client.NewRequesterContext(baseRoute, socketPath, logger)

	client := &ApiClient{
		context: context,
	}
	return client
}

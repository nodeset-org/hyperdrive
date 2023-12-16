package client

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client/routes"
)

// Rocket Pool client
type Client struct {
	Api *routes.ApiRequester

	configPath string
	debugPrint bool
}

// Create new Rocket Pool client from CLI context without checking for sync status
// Only use this function from commands that may work if the Daemon service doesn't exist
// Most users should call NewClientFromCtx(c).WithStatus() or NewClientFromCtx(c).WithReady()
func NewClient(configPath string, debug bool) (*Client, error) {
	client := &Client{
		configPath: configPath,
		debugPrint: debug,
	}


	client.Api = routes.NewApiRequester(socketPath),
	return client
}

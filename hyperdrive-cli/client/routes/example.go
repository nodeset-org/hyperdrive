package routes

import (
	"net/http"

	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

type ExampleRequester struct {
	client *http.Client
}

func NewExampleRequester(client *http.Client) *ExampleRequester {
	return &ExampleRequester{
		client: client,
	}
}

func (r *ExampleRequester) GetName() string {
	return "Example"
}
func (r *ExampleRequester) GetRoute() string {
	return "example"
}
func (r *ExampleRequester) GetClient() *http.Client {
	return r.client
}

// Get the response from a subset of the Rocket Pool daemon's `network` commands
func (r *ExampleRequester) CallDaemon(command string) (*api.ApiResponse[api.CallDaemonData], error) {
	args := map[string]string{
		"cmd": command,
	}
	return sendGetRequest[api.CallDaemonData](r, "call-daemon", "CallDaemon", args)
}

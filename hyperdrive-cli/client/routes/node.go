package routes

import (
	"net/http"

	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/types/api"
)

type NodeRequester struct {
	client *http.Client
}

func NewNodeRequester(client *http.Client) *NodeRequester {
	return &NodeRequester{
		client: client,
	}
}

func (r *NodeRequester) GetName() string {
	return "Node"
}
func (r *NodeRequester) GetRoute() string {
	return "node"
}
func (r *NodeRequester) GetClient() *http.Client {
	return r.client
}

// Get the response from a subset of the Rocket Pool daemon's `network` commands
func (r *NodeRequester) UploadDepositData(command string) (*api.ApiResponse[api.UploadDepositDataData], error) {
	args := map[string]string{
		"cmd": command,
	}
	return sendGetRequest[api.UploadDepositDataData](r, "upload-deposit-data", "UploadDepositData", args)
}

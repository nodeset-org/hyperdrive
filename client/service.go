package client

import (
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/client"
	"github.com/rocket-pool/node-manager-core/api/types"
)

type ServiceRequester struct {
	context client.IRequesterContext
}

func NewServiceRequester(context client.IRequesterContext) *ServiceRequester {
	return &ServiceRequester{
		context: context,
	}
}

func (r *ServiceRequester) GetName() string {
	return "Service"
}
func (r *ServiceRequester) GetRoute() string {
	return "service"
}
func (r *ServiceRequester) GetContext() client.IRequesterContext {
	return r.context
}

// Gets the status of the configured Execution and Beacon clients
func (r *ServiceRequester) ClientStatus() (*types.ApiResponse[api.ServiceClientStatusData], error) {
	return client.SendGetRequest[api.ServiceClientStatusData](r, "client-status", "ClientStatus", nil)
}

// Gets the Hyperdrive configuration
func (r *ServiceRequester) GetConfig() (*types.ApiResponse[api.ServiceGetConfigData], error) {
	return client.SendGetRequest[api.ServiceGetConfigData](r, "get-config", "GetConfig", nil)
}

// Restarts a Docker container
func (r *ServiceRequester) RestartContainer(container string) (*types.ApiResponse[types.SuccessData], error) {
	args := map[string]string{
		"container": container,
	}
	return client.SendGetRequest[types.SuccessData](r, "restart-container", "RestartContainer", args)
}

// Deletes the data folder including the wallet file, password file, and all validator keys.
// Don't use this unless you have a very good reason to do it (such as switching from Prater to Mainnet).
func (r *ServiceRequester) TerminateDataFolder() (*types.ApiResponse[api.ServiceTerminateDataFolderData], error) {
	return client.SendGetRequest[api.ServiceTerminateDataFolderData](r, "terminate-data-folder", "TerminateDataFolder", nil)
}

// Gets the version of the daemon
func (r *ServiceRequester) Version() (*types.ApiResponse[api.ServiceVersionData], error) {
	return client.SendGetRequest[api.ServiceVersionData](r, "version", "Version", nil)
}

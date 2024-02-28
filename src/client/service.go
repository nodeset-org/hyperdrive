package client

import (
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

type ServiceRequester struct {
	context *RequesterContext
}

func NewServiceRequester(context *RequesterContext) *ServiceRequester {
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
func (r *ServiceRequester) GetContext() *RequesterContext {
	return r.context
}

// Gets the status of the configured Execution and Beacon clients
func (r *ServiceRequester) ClientStatus() (*api.ApiResponse[api.ServiceClientStatusData], error) {
	return SendGetRequest[api.ServiceClientStatusData](r, "client-status", "ClientStatus", nil)
}

// Gets the Hyperdrive configuration
func (r *ServiceRequester) GetConfig() (*api.ApiResponse[api.ServiceGetConfigData], error) {
	return SendGetRequest[api.ServiceGetConfigData](r, "get-config", "GetConfig", nil)
}

// Restarts a Docker container
func (r *ServiceRequester) RestartContainer(container string) (*api.ApiResponse[api.SuccessData], error) {
	args := map[string]string{
		"container": container,
	}
	return SendGetRequest[api.SuccessData](r, "restart-container", "RestartContainer", args)
}

// Deletes the data folder including the wallet file, password file, and all validator keys.
// Don't use this unless you have a very good reason to do it (such as switching from Prater to Mainnet).
func (r *ServiceRequester) TerminateDataFolder() (*api.ApiResponse[api.ServiceTerminateDataFolderData], error) {
	return SendGetRequest[api.ServiceTerminateDataFolderData](r, "terminate-data-folder", "TerminateDataFolder", nil)
}

// Gets the version of the daemon
func (r *ServiceRequester) Version() (*api.ApiResponse[api.ServiceVersionData], error) {
	return SendGetRequest[api.ServiceVersionData](r, "version", "Version", nil)
}

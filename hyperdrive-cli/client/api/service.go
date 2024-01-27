package api

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
	return sendGetRequest[api.ServiceClientStatusData](r, "client-status", "ClientStatus", nil)
}

// Restarts the Validator client
func (r *ServiceRequester) RestartVc() (*api.ApiResponse[api.SuccessData], error) {
	return sendGetRequest[api.SuccessData](r, "restart-vc", "RestartVc", nil)
}

// Deletes the data folder including the wallet file, password file, and all validator keys.
// Don't use this unless you have a very good reason to do it (such as switching from Prater to Mainnet).
func (r *ServiceRequester) TerminateDataFolder() (*api.ApiResponse[api.ServiceTerminateDataFolderData], error) {
	return sendGetRequest[api.ServiceTerminateDataFolderData](r, "terminate-data-folder", "TerminateDataFolder", nil)
}

// Gets the version of the daemon
func (r *ServiceRequester) Version() (*api.ApiResponse[api.ServiceVersionData], error) {
	return sendGetRequest[api.ServiceVersionData](r, "version", "Version", nil)
}

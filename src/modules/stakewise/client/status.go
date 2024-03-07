package swclient

import (
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	nmc_client "github.com/rocket-pool/node-manager-core/api/client"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
)

type StatusRequester struct {
	context *nmc_client.RequesterContext
}

func NewStatusRequester(context *nmc_client.RequesterContext) *StatusRequester {
	return &StatusRequester{
		context: context,
	}
}

func (r *StatusRequester) GetName() string {
	return "Status"
}

func (r *StatusRequester) GetRoute() string {
	return "status"
}

func (r *StatusRequester) GetContext() *nmc_client.RequesterContext {
	return r.context
}

func (r *StatusRequester) GetActiveValidators() (*nmc_types.ApiResponse[swapi.ActiveValidatorsData], error) {
	return nmc_client.SendGetRequest[swapi.ActiveValidatorsData](r, "status", "Status", nil)
}

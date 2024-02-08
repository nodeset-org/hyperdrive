package swclient

import (
	"github.com/nodeset-org/hyperdrive/client"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

type StatusRequester struct {
	context *client.RequesterContext
}

func NewStatusRequester(context *client.RequesterContext) *StatusRequester {
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

func (r *StatusRequester) GetContext() *client.RequesterContext {
	return r.context
}

func (r *StatusRequester) GetActiveWallets() (*api.ApiResponse[swapi.ActiveValidatorsData], error) {
	return client.SendGetRequest[swapi.ActiveValidatorsData](r, "status", "Status", nil)
}

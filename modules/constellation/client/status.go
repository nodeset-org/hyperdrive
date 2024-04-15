package constclient

import (
	constapi "github.com/nodeset-org/hyperdrive/modules/constellation/shared/api"
	"github.com/rocket-pool/node-manager-core/api/client"
	"github.com/rocket-pool/node-manager-core/api/types"
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

func (r *StatusRequester) GetValidatorStatuses() (*types.ApiResponse[constapi.ConstellationStatusData], error) {
	return client.SendGetRequest[constapi.ConstellationStatusData](r, "status", "Status", nil)
}

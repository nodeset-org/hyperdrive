package client

import (
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

type StatusRequester struct {
	context *RequesterContext
}

func NewStatusRequester(context *RequesterContext) *StatusRequester {
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
func (r *StatusRequester) GetContext() *RequesterContext {
	return r.context
}

func (r *StatusRequester) GetActiveValidators() (*api.ApiResponse[api.SuccessData], error) {
	return SendGetRequest[api.SuccessData](r, "get-active-validators", "GetActiveValidators", nil)
}

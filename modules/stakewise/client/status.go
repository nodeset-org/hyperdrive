package swclient

import (
	"github.com/nodeset-org/hyperdrive/client"
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

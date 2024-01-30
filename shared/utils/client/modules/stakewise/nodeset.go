package swclient

import (
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	swapi "github.com/nodeset-org/hyperdrive/shared/types/api/modules/stakewise"
	"github.com/nodeset-org/hyperdrive/shared/utils/client"
)

type NodesetRequester struct {
	context *client.RequesterContext
}

func NewNodesetRequester(context *client.RequesterContext) *NodesetRequester {
	return &NodesetRequester{
		context: context,
	}
}

func (r *NodesetRequester) GetName() string {
	return "Nodeset"
}
func (r *NodesetRequester) GetRoute() string {
	return "nodeset"
}
func (r *NodesetRequester) GetContext() *client.RequesterContext {
	return r.context
}

// Upload the aggregated deposit data file to NodeSet's servers
func (r *NodesetRequester) UploadDepositData() (*api.ApiResponse[swapi.NodesetUploadDepositDataData], error) {
	return client.SendGetRequest[swapi.NodesetUploadDepositDataData](r, "upload-deposit-data", "UploadDepositData", nil)
}

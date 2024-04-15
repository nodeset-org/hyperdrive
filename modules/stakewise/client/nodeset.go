package swclient

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/rocket-pool/node-manager-core/api/client"
	"github.com/rocket-pool/node-manager-core/api/types"
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

// Set the validators root for the NodeSet vault
func (r *NodesetRequester) SetValidatorsRoot(root common.Hash) (*types.ApiResponse[types.TxInfoData], error) {
	args := map[string]string{
		"root": root.Hex(),
	}
	return client.SendGetRequest[types.TxInfoData](r, "set-validators-root", "SetValidatorsRoot", args)
}

// Upload the aggregated deposit data file to NodeSet's servers
func (r *NodesetRequester) UploadDepositData(bypassBalanceCheck bool) (*types.ApiResponse[swapi.NodesetUploadDepositDataData], error) {
	args := map[string]string{
		"bypassBalanceCheck": strconv.FormatBool(bypassBalanceCheck),
	}
	return client.SendGetRequest[swapi.NodesetUploadDepositDataData](r, "upload-deposit-data", "UploadDepositData", args)
}

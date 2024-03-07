package swclient

import (
	"github.com/ethereum/go-ethereum/common"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	nmc_client "github.com/rocket-pool/node-manager-core/api/client"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
)

type NodesetRequester struct {
	context *nmc_client.RequesterContext
}

func NewNodesetRequester(context *nmc_client.RequesterContext) *NodesetRequester {
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
func (r *NodesetRequester) GetContext() *nmc_client.RequesterContext {
	return r.context
}

// Set the validators root for the NodeSet vault
func (r *NodesetRequester) SetValidatorsRoot(root common.Hash) (*nmc_types.ApiResponse[nmc_types.TxInfoData], error) {
	args := map[string]string{
		"root": root.Hex(),
	}
	return nmc_client.SendGetRequest[nmc_types.TxInfoData](r, "set-validators-root", "SetValidatorsRoot", args)
}

// Upload the aggregated deposit data file to NodeSet's servers
func (r *NodesetRequester) UploadDepositData() (*nmc_types.ApiResponse[swapi.NodesetUploadDepositDataData], error) {
	return nmc_client.SendGetRequest[swapi.NodesetUploadDepositDataData](r, "upload-deposit-data", "UploadDepositData", nil)
}

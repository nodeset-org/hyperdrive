package swclient

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/client"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
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
func (r *NodesetRequester) SetValidatorsRoot(root common.Hash) (*api.ApiResponse[api.TxInfoData], error) {
	args := map[string]string{
		"root": root.Hex(),
	}
	return client.SendGetRequest[api.TxInfoData](r, "set-validators-root", "SetValidatorsRoot", args)
}

// Upload the aggregated deposit data file to NodeSet's servers
func (r *NodesetRequester) UploadDepositData(bypassBalanceCheck bool) (*api.ApiResponse[swapi.NodesetUploadDepositDataData], error) {
	fmt.Printf("!!bypassBalanceCheck: %v\n", bypassBalanceCheck)
	args := map[string]string{
		"bypassBalanceCheck": strconv.FormatBool(bypassBalanceCheck),
	}
	return client.SendGetRequest[swapi.NodesetUploadDepositDataData](r, "upload-deposit-data", "UploadDepositData", args)
}

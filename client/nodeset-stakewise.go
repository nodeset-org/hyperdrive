package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	apiv1 "github.com/nodeset-org/nodeset-client-go/api-v1"
	"github.com/rocket-pool/node-manager-core/api/client"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/beacon"
)

// Requester for StakeWise module calls to the nodeset.io service
type NodeSetStakeWiseRequester struct {
	context client.IRequesterContext
}

func NewNodeSetStakeWiseRequester(context client.IRequesterContext) *NodeSetStakeWiseRequester {
	return &NodeSetStakeWiseRequester{
		context: context,
	}
}

func (r *NodeSetStakeWiseRequester) GetName() string {
	return "NodeSet-StakeWise"
}
func (r *NodeSetStakeWiseRequester) GetRoute() string {
	return "nodeset/stakewise"
}
func (r *NodeSetStakeWiseRequester) GetContext() client.IRequesterContext {
	return r.context
}

// Gets the list of validators that the node has registered with the provided vault
func (r *NodeSetStakeWiseRequester) GetRegisteredValidators(vault common.Address) (*types.ApiResponse[api.NodeSetStakeWise_GetRegisteredValidatorsData], error) {
	args := map[string]string{
		"vault": vault.Hex(),
	}
	return client.SendGetRequest[api.NodeSetStakeWise_GetRegisteredValidatorsData](r, "get-registered-validators", "GetRegisteredValidators", args)
}

// Gets the version of the latest deposit data set on the server for the provided vault
func (r *NodeSetStakeWiseRequester) GetDepositDataSetVersion(vault common.Address) (*types.ApiResponse[api.NodeSetStakeWise_GetDepositDataSetData], error) {
	args := map[string]string{
		"vault": vault.Hex(),
	}
	return client.SendGetRequest[api.NodeSetStakeWise_GetDepositDataSetData](r, "get-deposit-data-set/version", "GetDepositDataSetVersion", args)
}

// Gets the latest deposit data set on the server for the provided vault
func (r *NodeSetStakeWiseRequester) GetDepositDataSet(vault common.Address) (*types.ApiResponse[api.NodeSetStakeWise_GetDepositDataSetData], error) {
	args := map[string]string{
		"vault": vault.Hex(),
	}
	return client.SendGetRequest[api.NodeSetStakeWise_GetDepositDataSetData](r, "get-deposit-data-set", "GetDepositDataSet", args)
}

// Uploads new validator deposit data to the NodeSet service
func (r *NodeSetStakeWiseRequester) UploadDepositData(data []beacon.ExtendedDepositData) (*types.ApiResponse[api.NodeSetStakeWise_UploadDepositDataData], error) {
	return client.SendPostRequest[api.NodeSetStakeWise_UploadDepositDataData](r, "upload-deposit-data", "UploadDepositData", data)
}

// Uploads signed exit messages to the NodeSet service
func (r *NodeSetStakeWiseRequester) UploadSignedExits(data []apiv1.ExitData) (*types.ApiResponse[api.NodeSetStakeWise_UploadSignedExitsData], error) {
	return client.SendPostRequest[api.NodeSetStakeWise_UploadSignedExitsData](r, "upload-signed-exits", "UploadSignedExits", data)
}

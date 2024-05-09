package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/client"
	"github.com/rocket-pool/node-manager-core/api/types"
)

type UtilsRequester struct {
	context client.IRequesterContext
}

func NewUtilsRequester(context client.IRequesterContext) *UtilsRequester {
	return &UtilsRequester{
		context: context,
	}
}

func (r *UtilsRequester) GetName() string {
	return "Utils"
}
func (r *UtilsRequester) GetRoute() string {
	return "utils"
}
func (r *UtilsRequester) GetContext() client.IRequesterContext {
	return r.context
}

// Resolves an ENS name or reserve resolves an address
func (r *UtilsRequester) ResolveEns(address common.Address, name string) (*types.ApiResponse[api.UtilsResolveEnsData], error) {
	args := map[string]string{
		"address": address.Hex(),
		"name":    name,
	}
	return client.SendGetRequest[api.UtilsResolveEnsData](r, "resolve-ens", "ResolveEns", args)
}

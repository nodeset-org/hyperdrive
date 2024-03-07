package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/client"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
)

type UtilsRequester struct {
	context *client.RequesterContext
}

func NewUtilsRequester(context *client.RequesterContext) *UtilsRequester {
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
func (r *UtilsRequester) GetContext() *client.RequesterContext {
	return r.context
}

// Resolves an ENS name or reserve resolves an address
func (r *UtilsRequester) ResolveEns(address common.Address, name string) (*nmc_types.ApiResponse[api.UtilsResolveEnsData], error) {
	args := map[string]string{
		"address": address.Hex(),
		"name":    name,
	}
	return client.SendGetRequest[api.UtilsResolveEnsData](r, "resolve-ens", "ResolveEns", args)
}

// Get the node's ETH balance
func (r *UtilsRequester) Balance() (*nmc_types.ApiResponse[api.UtilsBalanceData], error) {
	return client.SendGetRequest[api.UtilsBalanceData](r, "balance", "Balance", nil)
}

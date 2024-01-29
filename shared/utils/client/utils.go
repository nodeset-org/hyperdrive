package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

type UtilsRequester struct {
	context *requesterContext
}

func NewUtilsRequester(context *requesterContext) *UtilsRequester {
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
func (r *UtilsRequester) GetContext() *requesterContext {
	return r.context
}

// Resolves an ENS name or reserve resolves an address
func (r *UtilsRequester) ResolveEns(address common.Address, name string) (*api.ApiResponse[api.UtilsResolveEnsData], error) {
	args := map[string]string{
		"address": address.Hex(),
		"name":    name,
	}
	return sendGetRequest[api.UtilsResolveEnsData](r, "resolve-ens", "ResolveEns", args)
}

// Get the node's ETH balance
func (r *UtilsRequester) Balance() (*api.ApiResponse[api.UtilsBalanceData], error) {
	return sendGetRequest[api.UtilsBalanceData](r, "balance", "Balance", nil)
}

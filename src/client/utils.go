package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

type UtilsRequester struct {
	context *RequesterContext
}

func NewUtilsRequester(context *RequesterContext) *UtilsRequester {
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
func (r *UtilsRequester) GetContext() *RequesterContext {
	return r.context
}

// Resolves an ENS name or reserve resolves an address
func (r *UtilsRequester) ResolveEns(address common.Address, name string) (*api.ApiResponse[api.UtilsResolveEnsData], error) {
	args := map[string]string{
		"address": address.Hex(),
		"name":    name,
	}
	return SendGetRequest[api.UtilsResolveEnsData](r, "resolve-ens", "ResolveEns", args)
}

// Get the node's ETH balance
func (r *UtilsRequester) Balance() (*api.ApiResponse[api.UtilsBalanceData], error) {
	return SendGetRequest[api.UtilsBalanceData](r, "balance", "Balance", nil)
}

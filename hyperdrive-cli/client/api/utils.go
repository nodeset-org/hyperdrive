package api

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

type UtilsRequester struct {
	client *http.Client
}

func NewUtilsRequester(client *http.Client) *UtilsRequester {
	return &UtilsRequester{
		client: client,
	}
}

func (r *UtilsRequester) GetName() string {
	return "Utils"
}
func (r *UtilsRequester) GetRoute() string {
	return "utils"
}
func (r *UtilsRequester) GetClient() *http.Client {
	return r.client
}

// Resolves an ENS name or reserve resolves an address
func (r *UtilsRequester) ResolveEns(address common.Address, name string) (*api.ApiResponse[api.UtilsResolveEnsData], error) {
	args := map[string]string{
		"address": address.Hex(),
		"name":    name,
	}
	return sendGetRequest[api.UtilsResolveEnsData](r, "resolve-ens", "ResolveEns", args)
}

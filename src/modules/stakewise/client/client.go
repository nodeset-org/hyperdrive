package swclient

import (
	nmc_client "github.com/rocket-pool/node-manager-core/api/client"
)

// Binder for the Hyperdrive daemon API server
type ApiClient struct {
	context   *nmc_client.RequesterContext
	Nodeset   *NodesetRequester
	Validator *ValidatorRequester
	Wallet    *WalletRequester
	Status    *StatusRequester
}

// Creates a new API client instance
func NewApiClient(baseRoute string, socketPath string, debugMode bool) *ApiClient {
	context := nmc_client.NewRequesterContext(baseRoute, socketPath, debugMode)

	client := &ApiClient{
		context:   context,
		Nodeset:   NewNodesetRequester(context),
		Validator: NewValidatorRequester(context),
		Wallet:    NewWalletRequester(context),
		Status:    NewStatusRequester(context),
	}
	return client
}

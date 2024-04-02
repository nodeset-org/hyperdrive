package swclient

import (
	"log/slog"

	"github.com/rocket-pool/node-manager-core/api/client"
)

// Binder for the Hyperdrive daemon API server
type ApiClient struct {
	context   *client.RequesterContext
	Nodeset   *NodesetRequester
	Validator *ValidatorRequester
	Wallet    *WalletRequester
	Status    *StatusRequester
}

// Creates a new API client instance
func NewApiClient(baseRoute string, socketPath string, logger *slog.Logger) *ApiClient {
	context := client.NewRequesterContext(baseRoute, socketPath, logger)

	client := &ApiClient{
		context:   context,
		Nodeset:   NewNodesetRequester(context),
		Validator: NewValidatorRequester(context),
		Wallet:    NewWalletRequester(context),
		Status:    NewStatusRequester(context),
	}
	return client
}

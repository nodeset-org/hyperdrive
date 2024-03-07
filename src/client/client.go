package client

import (
	"github.com/rocket-pool/node-manager-core/api/client"
)

// Binder for the Hyperdrive daemon API server
type ApiClient struct {
	context *client.RequesterContext
	Service *ServiceRequester
	Tx      *TxRequester
	Utils   *UtilsRequester
	Wallet  *WalletRequester
}

// Creates a new API client instance
func NewApiClient(baseRoute string, socketPath string, debugMode bool) *ApiClient {
	context := client.NewRequesterContext(baseRoute, socketPath, debugMode)

	client := &ApiClient{
		context: context,
		Service: NewServiceRequester(context),
		Tx:      NewTxRequester(context),
		Utils:   NewUtilsRequester(context),
		Wallet:  NewWalletRequester(context),
	}

	return client
}

// Set debug mode
func (c *ApiClient) SetDebug(debug bool) {
	c.context.DebugMode = debug
}

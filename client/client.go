package client

import (
	"log/slog"
	"net/http/httptrace"
	"net/url"

	"github.com/rocket-pool/node-manager-core/api/client"
)

// Binder for the Hyperdrive daemon API server
type ApiClient struct {
	context               client.IRequesterContext
	NodeSet               *NodeSetRequester
	NodeSet_StakeWise     *NodeSetStakeWiseRequester
	NodeSet_Constellation *NodeSetConstellationRequester
	Service               *ServiceRequester
	Tx                    *TxRequester
	Utils                 *UtilsRequester
	Wallet                *WalletRequester
}

// Creates a new API client instance
func NewApiClient(apiUrl *url.URL, logger *slog.Logger, tracer *httptrace.ClientTrace) *ApiClient {
	context := client.NewNetworkRequesterContext(apiUrl, logger, tracer)

	client := &ApiClient{
		context:               context,
		NodeSet:               NewNodeSetRequester(context),
		NodeSet_StakeWise:     NewNodeSetStakeWiseRequester(context),
		NodeSet_Constellation: NewNodeSetConstellationRequester(context),
		Service:               NewServiceRequester(context),
		Tx:                    NewTxRequester(context),
		Utils:                 NewUtilsRequester(context),
		Wallet:                NewWalletRequester(context),
	}
	return client
}

// Set debug mode
func (c *ApiClient) SetLogger(logger *slog.Logger) {
	c.context.SetLogger(logger)
}

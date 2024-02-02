package swclient

import (
	"context"
	"net"
	"net/http"

	"github.com/fatih/color"
	"github.com/nodeset-org/hyperdrive/client"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

const (
	baseUrl         string          = "http://" + swconfig.ModuleName + "/%s"
	jsonContentType string          = "application/json"
	apiColor        color.Attribute = color.FgHiCyan
)

// Binder for the Hyperdrive daemon API server
type ApiClient struct {
	Nodeset *NodesetRequester
	Wallet  *WalletRequester

	context *client.RequesterContext
}

// Creates a new API client instance
func NewApiClient(baseRoute string, socketPath string, debugMode bool) *ApiClient {
	apiRequester := &ApiClient{
		context: &client.RequesterContext{
			SocketPath: socketPath,
			DebugMode:  debugMode,
			Base:       baseRoute,
		},
	}

	apiRequester.context.Client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}

	log := log.NewColorLogger(apiColor)
	apiRequester.context.Log = &log

	apiRequester.Nodeset = NewNodesetRequester(apiRequester.context)
	apiRequester.Wallet = NewWalletRequester(apiRequester.context)
	return apiRequester
}

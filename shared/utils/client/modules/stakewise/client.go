package swclient

import (
	"context"
	"net"
	"net/http"

	"github.com/fatih/color"
	swconfig "github.com/nodeset-org/hyperdrive/shared/config/modules/stakewise"
	"github.com/nodeset-org/hyperdrive/shared/utils/client"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

const (
	baseUrl         string          = "http://" + swconfig.DaemonRoute + "/%s"
	jsonContentType string          = "application/json"
	apiColor        color.Attribute = color.FgHiCyan
)

// Binder for the Hyperdrive daemon API server
type ApiClient struct {
	Wallet *WalletRequester

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

	apiRequester.Wallet = NewWalletRequester(apiRequester.context)
	return apiRequester
}

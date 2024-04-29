package wallet

import (
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type walletStatusContextFactory struct {
	handler *WalletHandler
}

func (f *walletStatusContextFactory) Create(args url.Values) (*walletStatusContext, error) {
	c := &walletStatusContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletStatusContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletStatusContext, api.WalletStatusData](
		router, "status", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletStatusContext struct {
	handler *WalletHandler
}

func (c *walletStatusContext) PrepareData(data *api.WalletStatusData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	status, err := w.GetStatus()
	if err != nil {
		return types.ResponseStatus_Error, err
	}
	data.WalletStatus = status
	return types.ResponseStatus_Success, nil
}

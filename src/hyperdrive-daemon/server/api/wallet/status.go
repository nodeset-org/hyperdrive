package wallet

import (
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
)

// ===============
// === Factory ===
// ===============

type walletStatusFactory struct {
	handler *WalletHandler
}

func (f *walletStatusFactory) Create(args url.Values) (*walletStatusContext, error) {
	c := &walletStatusContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletStatusFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*walletStatusContext, api.WalletStatusData](
		router, "status", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletStatusContext struct {
	handler *WalletHandler
}

func (c *walletStatusContext) PrepareData(data *api.WalletStatusData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	status, err := w.GetStatus()
	if err != nil {
		return err
	}
	data.WalletStatus = status
	return nil
}

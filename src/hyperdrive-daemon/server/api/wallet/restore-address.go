package wallet

import (
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
)

// ===============
// === Factory ===
// ===============

type walletRestoreAddressContextFactory struct {
	handler *WalletHandler
}

func (f *walletRestoreAddressContextFactory) Create(args url.Values) (*walletRestoreAddressContext, error) {
	c := &walletRestoreAddressContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletRestoreAddressContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*walletRestoreAddressContext, api.SuccessData](
		router, "restore-address", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletRestoreAddressContext struct {
	handler *WalletHandler
}

func (c *walletRestoreAddressContext) PrepareData(data *api.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	return w.RestoreAddressToWallet()
}

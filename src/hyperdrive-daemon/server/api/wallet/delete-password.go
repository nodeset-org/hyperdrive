package wallet

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type walletDeletePasswordContextFactory struct {
	handler *WalletHandler
}

func (f *walletDeletePasswordContextFactory) Create(args url.Values) (*walletDeletePasswordContext, error) {
	c := &walletDeletePasswordContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletDeletePasswordContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*walletDeletePasswordContext, nmc_types.SuccessData](
		router, "delete-password", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletDeletePasswordContext struct {
	handler  *WalletHandler
	password []byte
	save     bool
}

func (c *walletDeletePasswordContext) PrepareData(data *nmc_types.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	err := w.DeletePassword()
	if err != nil {
		return fmt.Errorf("error deleting wallet password from disk: %w", err)
	}
	return nil
}

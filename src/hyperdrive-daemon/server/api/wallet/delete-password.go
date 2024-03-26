package wallet

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
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
	utils.RegisterQuerylessGet[*walletDeletePasswordContext, api.SuccessData](
		router, "delete-password", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletDeletePasswordContext struct {
	handler *WalletHandler
}

func (c *walletDeletePasswordContext) PrepareData(data *api.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	err := w.DeletePassword()
	if err != nil {
		return fmt.Errorf("error deleting wallet password from disk: %w", err)
	}
	return nil
}

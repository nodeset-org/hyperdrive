package wallet

import (
	"errors"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletMasqueradeContextFactory struct {
	handler *WalletHandler
}

func (f *walletMasqueradeContextFactory) Create(args url.Values) (*walletMasqueradeContext, error) {
	c := &walletMasqueradeContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("address", args, input.ValidateAddress, &c.address),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletMasqueradeContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletMasqueradeContext, api.SuccessData](
		router, "masquerade", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletMasqueradeContext struct {
	handler *WalletHandler
	address common.Address
}

func (c *walletMasqueradeContext) PrepareData(data *api.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	return w.MasqueradeAsAddress(c.address)
}

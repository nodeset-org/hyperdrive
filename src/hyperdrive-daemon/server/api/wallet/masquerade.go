package wallet

import (
	"errors"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/rocket-pool/node-manager-core/api/server"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
	nmc_input "github.com/rocket-pool/node-manager-core/utils/input"
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
		server.ValidateArg("address", args, nmc_input.ValidateAddress, &c.address),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletMasqueradeContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletMasqueradeContext, nmc_types.SuccessData](
		router, "masquerade", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletMasqueradeContext struct {
	handler *WalletHandler
	address common.Address
}

func (c *walletMasqueradeContext) PrepareData(data *nmc_types.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	return w.MasqueradeAsAddress(c.address)
}

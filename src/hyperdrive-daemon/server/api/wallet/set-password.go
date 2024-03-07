package wallet

import (
	"errors"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
	nmc_input "github.com/rocket-pool/node-manager-core/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletSetPasswordContextFactory struct {
	handler *WalletHandler
}

func (f *walletSetPasswordContextFactory) Create(args url.Values) (*walletSetPasswordContext, error) {
	c := &walletSetPasswordContext{
		handler: f.handler,
	}
	inputErrs := []error{
		nmc_server.ValidateArg("password", args, nmc_input.ValidateNodePassword, &c.password),
		nmc_server.ValidateArg("save", args, nmc_input.ValidateBool, &c.save),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletSetPasswordContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*walletSetPasswordContext, nmc_types.SuccessData](
		router, "set-password", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletSetPasswordContext struct {
	handler  *WalletHandler
	password string
	save     bool
}

func (c *walletSetPasswordContext) PrepareData(data *nmc_types.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	return w.SetPassword(c.password, c.save)
}

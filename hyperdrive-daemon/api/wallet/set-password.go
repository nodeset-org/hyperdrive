package wallet

import (
	"errors"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
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
		server.ValidateArg("password", args, input.ValidateNodePassword, &c.password),
		server.ValidateArg("save", args, input.ValidateBool, &c.save),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletSetPasswordContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletSetPasswordContext, api.SuccessData](
		router, "set-password", f, f.handler.serviceProvider,
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

func (c *walletSetPasswordContext) PrepareData(data *api.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	return w.SetPassword(c.password, c.save)
}

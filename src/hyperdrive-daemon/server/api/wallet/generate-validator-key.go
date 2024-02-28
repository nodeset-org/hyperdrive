package wallet

import (
	"errors"
	"fmt"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type walletGenerateValidatorKeyContextFactory struct {
	handler *WalletHandler
}

func (f *walletGenerateValidatorKeyContextFactory) Create(args url.Values) (*walletGenerateValidatorKeyContext, error) {
	c := &walletGenerateValidatorKeyContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.GetStringFromVars("path", args, &c.path),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletGenerateValidatorKeyContextFactory) RegisterRoute(router *mux.Router) {
	utils.RegisterQuerylessGet[*walletGenerateValidatorKeyContext, api.WalletGenerateValidatorKeyData](
		router, "generate-validator-key", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletGenerateValidatorKeyContext struct {
	handler *WalletHandler
	path    string
}

func (c *walletGenerateValidatorKeyContext) PrepareData(data *api.WalletGenerateValidatorKeyData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	key, err := w.GenerateValidatorKey(c.path)
	if err != nil {
		return fmt.Errorf("error generating validator key: %w", err)
	}

	data.PrivateKey = key
	return nil
}

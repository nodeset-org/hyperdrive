package wallet

import (
	"errors"
	"fmt"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
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
	server.RegisterQuerylessGet[*walletGenerateValidatorKeyContext, api.WalletGenerateValidatorKeyData](
		router, "generate-validator-key", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletGenerateValidatorKeyContext struct {
	handler *WalletHandler
	path    string
}

func (c *walletGenerateValidatorKeyContext) PrepareData(data *api.WalletGenerateValidatorKeyData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return types.ResponseStatus_WalletNotReady, err
	}

	key, err := w.GenerateValidatorKey(c.path)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error generating validator key: %w", err)
	}

	data.PrivateKey = key
	return types.ResponseStatus_Success, nil
}

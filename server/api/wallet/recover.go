package wallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/utils/input"
	"github.com/rocket-pool/node-manager-core/wallet"
)

// ===============
// === Factory ===
// ===============

type walletRecoverContextFactory struct {
	handler *WalletHandler
}

func (f *walletRecoverContextFactory) Create(args url.Values) (*walletRecoverContext, error) {
	c := &walletRecoverContext{
		handler: f.handler,
	}
	server.GetOptionalStringFromVars("derivation-path", args, &c.derivationPath)
	inputErrs := []error{
		server.ValidateArg("mnemonic", args, input.ValidateWalletMnemonic, &c.mnemonic),
		server.ValidateOptionalArg("index", args, input.ValidateUint, &c.index, nil),
		server.ValidateArg("password", args, input.ValidateNodePassword, &c.password),
		server.ValidateArg("save-password", args, input.ValidateBool, &c.savePassword),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletRecoverContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletRecoverContext, api.WalletRecoverData](
		router, "recover", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletRecoverContext struct {
	handler        *WalletHandler
	mnemonic       string
	derivationPath string
	index          uint64
	password       string
	savePassword   bool
}

func (c *walletRecoverContext) PrepareData(data *api.WalletRecoverData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	// Requirements
	status, err := w.GetStatus()
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting wallet status: %w", err)
	}
	if status.Wallet.IsOnDisk {
		return types.ResponseStatus_ResourceConflict, fmt.Errorf("a wallet is already present")
	}

	// Parse the derivation path
	path, err := wallet.GetDerivationPath(wallet.DerivationPath(c.derivationPath))
	if err != nil {
		return types.ResponseStatus_InvalidArguments, err
	}

	// Recover the wallet
	err = w.Recover(path, uint(c.index), c.mnemonic, c.password, c.savePassword, false)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error recovering wallet: %w", err)
	}
	data.AccountAddress, _ = w.GetAddress()
	return types.ResponseStatus_Success, nil
}

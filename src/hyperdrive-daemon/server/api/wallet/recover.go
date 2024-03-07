package wallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	nmc_nodewallet "github.com/rocket-pool/node-manager-core/node/wallet"
	nmc_input "github.com/rocket-pool/node-manager-core/utils/input"
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
		server.ValidateArg("mnemonic", args, nmc_input.ValidateWalletMnemonic, &c.mnemonic),
		server.ValidateOptionalArg("index", args, nmc_input.ValidateUint, &c.index, nil),
		server.ValidateArg("password", args, nmc_input.ValidateNodePassword, &c.password),
		server.ValidateArg("save-password", args, nmc_input.ValidateBool, &c.savePassword),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletRecoverContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletRecoverContext, api.WalletRecoverData](
		router, "recover", f, f.handler.serviceProvider.ServiceProvider,
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

func (c *walletRecoverContext) PrepareData(data *api.WalletRecoverData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	// Requirements
	status, err := w.GetStatus()
	if err != nil {
		return fmt.Errorf("error getting wallet status: %w", err)
	}
	if status.Wallet.IsOnDisk {
		return fmt.Errorf("a wallet is already present")
	}

	// Parse the derivation path
	path, err := nmc_nodewallet.GetDerivationPath(wallet.DerivationPath(c.derivationPath))
	if err != nil {
		return err
	}

	// Recover the wallet
	err = w.Recover(path, uint(c.index), c.mnemonic, c.password, c.savePassword, false)
	if err != nil {
		return fmt.Errorf("error recovering wallet: %w", err)
	}
	data.AccountAddress, _ = w.GetAddress()
	return nil
}

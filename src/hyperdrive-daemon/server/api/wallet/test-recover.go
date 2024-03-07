package wallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
	nmc_nodewallet "github.com/rocket-pool/node-manager-core/node/wallet"
	nmc_input "github.com/rocket-pool/node-manager-core/utils/input"
	nmc_wallet "github.com/rocket-pool/node-manager-core/wallet"
)

// ===============
// === Factory ===
// ===============

type walletTestRecoverContextFactory struct {
	handler *WalletHandler
}

func (f *walletTestRecoverContextFactory) Create(args url.Values) (*walletTestRecoverContext, error) {
	c := &walletTestRecoverContext{
		handler: f.handler,
	}
	nmc_server.GetOptionalStringFromVars("derivation-path", args, &c.derivationPath)
	inputErrs := []error{
		nmc_server.ValidateArg("mnemonic", args, nmc_input.ValidateWalletMnemonic, &c.mnemonic),
		nmc_server.ValidateOptionalArg("index", args, nmc_input.ValidateUint, &c.index, nil),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletTestRecoverContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*walletTestRecoverContext, api.WalletRecoverData](
		router, "test-recover", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletTestRecoverContext struct {
	handler        *WalletHandler
	mnemonic       string
	derivationPath string
	index          uint64
}

func (c *walletTestRecoverContext) PrepareData(data *api.WalletRecoverData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	rs := sp.GetResources()

	// Parse the derivation path
	path, err := nmc_nodewallet.GetDerivationPath(nmc_wallet.DerivationPath(c.derivationPath))
	if err != nil {
		return err
	}

	// Recover the wallet
	w, err := nmc_nodewallet.TestRecovery(path, uint(c.index), c.mnemonic, rs.ChainID)
	if err != nil {
		return fmt.Errorf("error recovering wallet: %w", err)
	}
	data.AccountAddress, _ = w.GetAddress()
	return nil
}

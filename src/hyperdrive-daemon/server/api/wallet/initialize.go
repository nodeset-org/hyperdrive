package wallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/wallet"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletInitializeContextFactory struct {
	handler *WalletHandler
}

func (f *walletInitializeContextFactory) Create(args url.Values) (*walletInitializeContext, error) {
	c := &walletInitializeContext{
		handler: f.handler,
	}
	server.GetOptionalStringFromVars("derivation-path", args, &c.derivationPath)
	inputErrs := []error{
		server.ValidateOptionalArg("index", args, input.ValidateUint, &c.index, nil),
		server.ValidateArg("password", args, input.ValidateNodePassword, &c.password),
		server.ValidateArg("save-wallet", args, input.ValidateBool, &c.saveWallet),
		server.ValidateArg("save-password", args, input.ValidateBool, &c.savePassword),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletInitializeContextFactory) RegisterRoute(router *mux.Router) {
	utils.RegisterQuerylessGet[*walletInitializeContext, api.WalletInitializeData](
		router, "initialize", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletInitializeContext struct {
	handler        *WalletHandler
	derivationPath string
	index          uint64
	password       string
	passwordExists bool
	savePassword   bool
	saveWallet     bool
}

func (c *walletInitializeContext) PrepareData(data *api.WalletInitializeData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider

	// Parse the derivation path
	path, err := GetDerivationPath(types.DerivationPath(c.derivationPath))
	if err != nil {
		return err
	}

	var w *wallet.Wallet
	var mnemonic string
	if !c.saveWallet {
		// Make a dummy wallet for the sake of creating a mnemonic and derived address
		mnemonic, err = wallet.GenerateNewMnemonic()
		if err != nil {
			return fmt.Errorf("error generating new mnemonic: %w", err)
		}

		w, err = wallet.TestRecovery(path, uint(c.index), mnemonic, 0)
		if err != nil {
			return fmt.Errorf("error generating wallet from new mnemonic: %w", err)
		}
	} else {
		// Initialize the daemon wallet
		w = sp.GetWallet()

		// Requirements
		status, err := w.GetStatus()
		if err != nil {
			return fmt.Errorf("error getting wallet status: %w", err)
		}
		if status.Wallet.IsOnDisk {
			return fmt.Errorf("a wallet is already present")
		}

		// Create the new wallet
		mnemonic, err = w.CreateNewLocalWallet(path, uint(c.index), c.password, c.savePassword)
		if err != nil {
			return fmt.Errorf("error initializing new wallet: %w", err)
		}
	}

	data.Mnemonic = mnemonic
	data.AccountAddress, _ = w.GetAddress()
	return nil
}

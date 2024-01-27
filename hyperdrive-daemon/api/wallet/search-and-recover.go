package wallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/wallet"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

const (
	findIterations uint = 100000
)

// ===============
// === Factory ===
// ===============

type walletSearchAndRecoverContextFactory struct {
	handler *WalletHandler
}

func (f *walletSearchAndRecoverContextFactory) Create(args url.Values) (*walletSearchAndRecoverContext, error) {
	c := &walletSearchAndRecoverContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("mnemonic", args, input.ValidateWalletMnemonic, &c.mnemonic),
		server.ValidateArg("address", args, input.ValidateAddress, &c.address),
		server.ValidateOptionalArg("password", args, input.ValidateNodePassword, &c.password, &c.passwordExists),
		server.ValidateOptionalArg("save-password", args, input.ValidateBool, &c.savePassword, nil),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletSearchAndRecoverContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletSearchAndRecoverContext, api.WalletSearchAndRecoverData](
		router, "search-and-recover", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletSearchAndRecoverContext struct {
	handler        *WalletHandler
	mnemonic       string
	address        common.Address
	password       []byte
	passwordExists bool
	savePassword   bool
}

func (c *walletSearchAndRecoverContext) PrepareData(data *api.WalletSearchAndRecoverData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()
	rs := sp.GetResources()

	// Requirements
	status := w.GetStatus()
	if status.HasKeystore {
		return fmt.Errorf("a wallet is already present")
	}

	_, hasPassword := w.GetPassword()
	if !hasPassword && !c.passwordExists {
		return fmt.Errorf("you must set a password before recovering a wallet, or provide one in this call")
	}
	w.RememberPassword(c.password)
	if c.savePassword {
		err := w.SavePassword()
		if err != nil {
			return fmt.Errorf("error saving wallet password to disk: %w", err)
		}
	}

	// Try each derivation path across all of the iterations
	paths := []string{
		wallet.DefaultNodeKeyPath,
		wallet.LedgerLiveNodeKeyPath,
		wallet.MyEtherWalletNodeKeyPath,
	}
	for i := uint(0); i < findIterations; i++ {
		for j := 0; j < len(paths); j++ {
			derivationPath := paths[j]
			recoveredWallet, err := wallet.TestRecovery(derivationPath, i, c.mnemonic, rs.ChainID)
			if err != nil {
				return fmt.Errorf("error recovering wallet with path [%s], index [%d]: %w", derivationPath, i, err)
			}

			// Get recovered account
			recoveredAddress, _ := recoveredWallet.GetAddress()
			if recoveredAddress == c.address {
				// We found the correct derivation path and index
				data.FoundWallet = true
				data.DerivationPath = derivationPath
				data.Index = i
				break
			}
		}
		if data.FoundWallet {
			break
		}
	}

	if !data.FoundWallet {
		return fmt.Errorf("exhausted all derivation paths and indices from 0 to %d, wallet not found", findIterations)
	}

	// Recover the wallet
	err := w.Recover(data.DerivationPath, data.Index, c.mnemonic)
	if err != nil {
		return fmt.Errorf("error recovering wallet: %w", err)
	}
	data.AccountAddress, _ = w.GetAddress()
	return nil
}

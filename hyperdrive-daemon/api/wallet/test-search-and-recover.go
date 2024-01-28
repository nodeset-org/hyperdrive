package wallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemons/common/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/wallet"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletTestSearchAndRecoverContextFactory struct {
	handler *WalletHandler
}

func (f *walletTestSearchAndRecoverContextFactory) Create(args url.Values) (*walletTestSearchAndRecoverContext, error) {
	c := &walletTestSearchAndRecoverContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("mnemonic", args, input.ValidateWalletMnemonic, &c.mnemonic),
		server.ValidateArg("address", args, input.ValidateAddress, &c.address),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletTestSearchAndRecoverContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletTestSearchAndRecoverContext, api.WalletSearchAndRecoverData](
		router, "test-search-and-recover", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletTestSearchAndRecoverContext struct {
	handler  *WalletHandler
	mnemonic string
	address  common.Address
}

func (c *walletTestSearchAndRecoverContext) PrepareData(data *api.WalletSearchAndRecoverData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	rs := sp.GetResources()

	// Try each derivation path across all of the iterations
	var recoveredWallet *wallet.Wallet
	paths := []string{
		wallet.DefaultNodeKeyPath,
		wallet.LedgerLiveNodeKeyPath,
		wallet.MyEtherWalletNodeKeyPath,
	}
	for i := uint(0); i < findIterations; i++ {
		for j := 0; j < len(paths); j++ {
			var err error
			derivationPath := paths[j]
			recoveredWallet, err = wallet.TestRecovery(derivationPath, i, c.mnemonic, rs.ChainID)
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
	data.AccountAddress, _ = recoveredWallet.GetAddress()
	return nil
}

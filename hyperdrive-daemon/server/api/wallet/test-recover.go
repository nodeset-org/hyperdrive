package wallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/wallet"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
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
	server.GetOptionalStringFromVars("derivation-path", args, &c.derivationPath)
	inputErrs := []error{
		server.ValidateArg("mnemonic", args, input.ValidateWalletMnemonic, &c.mnemonic),
		server.ValidateOptionalArg("index", args, input.ValidateUint, &c.index, nil),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletTestRecoverContextFactory) RegisterRoute(router *mux.Router) {
	utils.RegisterQuerylessGet[*walletTestRecoverContext, api.WalletRecoverData](
		router, "test-recover", f, f.handler.serviceProvider,
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
	pathType := types.DerivationPath(c.derivationPath)
	var path string
	switch pathType {
	case types.DerivationPath_Default:
		path = wallet.DefaultNodeKeyPath
	case types.DerivationPath_LedgerLive:
		path = wallet.LedgerLiveNodeKeyPath
	case types.DerivationPath_Mew:
		path = wallet.MyEtherWalletNodeKeyPath
	default:
		return fmt.Errorf("[%s] is not a valid derivation path type", c.derivationPath)
	}

	// Recover the wallet
	w, err := wallet.TestRecovery(path, uint(c.index), c.mnemonic, rs.ChainID)
	if err != nil {
		return fmt.Errorf("error recovering wallet: %w", err)
	}
	data.AccountAddress, _ = w.GetAddress()
	return nil
}

package wallet

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type walletExportContextFactory struct {
	handler *WalletHandler
}

func (f *walletExportContextFactory) Create(args url.Values) (*walletExportContext, error) {
	c := &walletExportContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletExportContextFactory) RegisterRoute(router *mux.Router) {
	utils.RegisterQuerylessGet[*walletExportContext, api.WalletExportData](
		router, "export", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletExportContext struct {
	handler *WalletHandler
}

func (c *walletExportContext) PrepareData(data *api.WalletExportData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return err
	}

	// Get password
	pw, isSet, err := w.GetPassword()
	if err != nil {
		return fmt.Errorf("error getting wallet password: %w", err)
	}
	if !isSet {
		return fmt.Errorf("password has not been set; cannot decrypt wallet keystore without it")
	}
	data.Password = pw

	// Serialize wallet
	walletString, err := w.SerializeData()
	if err != nil {
		return fmt.Errorf("error serializing wallet keystore: %w", err)
	}
	data.Wallet = walletString

	// Get account private key
	data.AccountPrivateKey, err = w.GetNodePrivateKeyBytes()
	if err != nil {
		return fmt.Errorf("error getting node wallet private key: %w", err)
	}

	return nil
}

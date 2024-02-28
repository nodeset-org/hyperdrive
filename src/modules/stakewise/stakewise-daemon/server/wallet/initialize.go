package swwallet

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	api "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
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
	return c, nil
}

func (f *walletInitializeContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletInitializeContext, api.WalletInitializeData](
		router, "initialize", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletInitializeContext struct {
	handler *WalletHandler
}

func (c *walletInitializeContext) PrepareData(data *api.WalletInitializeData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	client := sp.GetHyperdriveClient()

	// Get the wallet status
	response, err := client.Wallet.Status()
	if err != nil {
		return fmt.Errorf("error getting wallet status: %w", err)
	}
	status := response.Data.WalletStatus
	if !status.Wallet.IsLoaded {
		return fmt.Errorf("Hyperdrive does not currently have a wallet ready")
	}

	// Requirements
	/*
		err = sp.RequireWalletReady()
		if err != nil {
			return err
		}
	*/

	// Get the Geth keystore in JSON format
	ethkeyResponse, err := client.Wallet.ExportEthKey()
	if err != nil {
		return fmt.Errorf("error getting geth-style keystore: %w", err)
	}
	ethKey := ethkeyResponse.Data.EthKeyJson
	password := ethkeyResponse.Data.Password

	// Write the wallet to disk
	moduleDir := sp.GetModuleDir()
	walletPath := filepath.Join(moduleDir, swconfig.WalletFilename)
	err = os.WriteFile(walletPath, ethKey, 0600)
	if err != nil {
		return fmt.Errorf("error saving wallet keystore to disk: %w", err)
	}

	// Write the password to disk
	passwordPath := filepath.Join(moduleDir, swconfig.PasswordFilename)
	err = os.WriteFile(passwordPath, []byte(password), 0600)
	if err != nil {
		return fmt.Errorf("error saving wallet password to disk: %w", err)
	}

	data.AccountAddress = status.Wallet.WalletAddress
	return nil
}

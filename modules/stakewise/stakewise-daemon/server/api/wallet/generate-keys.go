package swwallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	api "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletGenerateKeysContextFactory struct {
	handler *WalletHandler
}

func (f *walletGenerateKeysContextFactory) Create(args url.Values) (*walletGenerateKeysContext, error) {
	c := &walletGenerateKeysContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("count", args, input.ValidateUint, &c.count),
		server.ValidateArg("restart-vc", args, input.ValidateBool, &c.restartVc),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletGenerateKeysContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletGenerateKeysContext, api.WalletGenerateKeysData](
		router, "generate-keys", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletGenerateKeysContext struct {
	handler   *WalletHandler
	count     uint64
	restartVc bool
}

func (c *walletGenerateKeysContext) PrepareData(data *api.WalletGenerateKeysData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	client := sp.GetHyperdriveClient()
	wallet := sp.GetWallet()

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

	// Generate and save the keys
	pubkeys := make([]beacon.ValidatorPubkey, c.count)
	for i := 0; i < int(c.count); i++ {
		key, err := wallet.GenerateNewValidatorKey()
		if err != nil {
			return fmt.Errorf("error generating validator key: %w", err)
		}
		pubkeys[i] = beacon.ValidatorPubkey(key.PublicKey().Marshal())
	}
	data.Pubkeys = pubkeys

	// Restart the VC
	if c.restartVc {
		_, err = client.Service.RestartContainer(swconfig.VcContainerSuffix)
		if err != nil {
			return err
		}
	}
	return nil
}

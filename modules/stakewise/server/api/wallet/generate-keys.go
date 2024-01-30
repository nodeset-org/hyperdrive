package swwallet

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
	api "github.com/nodeset-org/hyperdrive/shared/types/api/modules/stakewise"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletwalletGenerateKeysContextFactory struct {
	handler *WalletHandler
}

func (f *walletwalletGenerateKeysContextFactory) Create(args url.Values) (*walletGenerateKeysContext, error) {
	c := &walletGenerateKeysContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("count", args, input.ValidateUint, &c.count),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletwalletGenerateKeysContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletGenerateKeysContext, api.WalletGenerateKeysData](
		router, "generate-keys", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletGenerateKeysContext struct {
	handler *WalletHandler
	count   uint64
}

func (c *walletGenerateKeysContext) PrepareData(data *api.WalletGenerateKeysData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	client := sp.GetClient()
	wallet := sp.GetWallet()
	ddMgr := sp.GetDepositDataManager()

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

	// Regen the deposit data file
	totalCount, err := ddMgr.RegenerateDepositData()
	if err != nil {
		return fmt.Errorf("error regenerating deposit data: %w", err)
	}

	data.Pubkeys = pubkeys
	data.TotalCount = totalCount
	return nil
}

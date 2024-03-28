package wallet

import (
	"context"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type walletBalanceContextFactory struct {
	handler *WalletHandler
}

func (f *walletBalanceContextFactory) Create(args url.Values) (*walletBalanceContext, error) {
	c := &walletBalanceContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletBalanceContextFactory) RegisterRoute(router *mux.Router) {
	utils.RegisterQuerylessGet[*walletBalanceContext, api.WalletBalanceData](
		router, "balance", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletBalanceContext struct {
	handler *WalletHandler
}

func (c *walletBalanceContext) PrepareData(data *api.WalletBalanceData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()
	ec := sp.GetEthClient()
	nodeAddress, _ := w.GetAddress()

	// Requirements
	err := sp.RequireEthClientSynced(context.Background())
	if err != nil {
		return err
	}

	balance, err := ec.BalanceAt(context.Background(), nodeAddress, nil)
	if err != nil {
		return err
	}

	data.Balance = balance
	return nil
}

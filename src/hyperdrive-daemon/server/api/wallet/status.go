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

type walletStatusFactory struct {
	handler *WalletHandler
}

func (f *walletStatusFactory) Create(args url.Values) (*walletStatusContext, error) {
	c := &walletStatusContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletStatusFactory) RegisterRoute(router *mux.Router) {
	utils.RegisterQuerylessGet[*walletStatusContext, api.WalletStatusData](
		router, "status", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletStatusContext struct {
	handler *WalletHandler
}

func (c *walletStatusContext) PrepareData(data *api.WalletStatusData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()
	ec := sp.GetEthClient()
	nodeAddress, _ := w.GetAddress()

	balance, err := ec.BalanceAt(context.Background(), nodeAddress, nil)
	if err != nil {
		return err
	}

	status, err := w.GetStatus()
	if err != nil {
		return err
	}

	status.Wallet.WalletBalance = *balance
	data.WalletStatus = status
	return nil
}

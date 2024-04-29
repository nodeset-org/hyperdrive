package wallet

import (
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
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
	server.RegisterQuerylessGet[*walletBalanceContext, api.WalletBalanceData](
		router, "balance", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletBalanceContext struct {
	handler *WalletHandler
}

func (c *walletBalanceContext) PrepareData(data *api.WalletBalanceData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ctx := c.handler.ctx
	w := sp.GetWallet()
	ec := sp.GetEthClient()
	nodeAddress, _ := w.GetAddress()

	// Requirements
	err := sp.RequireNodeAddress()
	if err != nil {
		return types.ResponseStatus_AddressNotPresent, err
	}
	err = sp.RequireEthClientSynced(ctx)
	if err != nil {
		return types.ResponseStatus_ClientsNotSynced, err
	}

	data.Balance, err = ec.BalanceAt(ctx, nodeAddress, nil)
	if err != nil {
		return types.ResponseStatus_Error, err
	}
	return types.ResponseStatus_Success, nil
}

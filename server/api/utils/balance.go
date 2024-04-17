package utils

import (
	"fmt"
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

type utilsBalanceContextFactory struct {
	handler *UtilsHandler
}

func (f *utilsBalanceContextFactory) Create(args url.Values) (*utilsBalanceContext, error) {
	c := &utilsBalanceContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *utilsBalanceContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*utilsBalanceContext, api.UtilsBalanceData](
		router, "balance", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type utilsBalanceContext struct {
	handler *UtilsHandler
}

func (c *utilsBalanceContext) PrepareData(data *api.UtilsBalanceData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ec := sp.GetEthClient()
	ctx := c.handler.ctx
	nodeAddress, _ := sp.GetWallet().GetAddress()

	// Requirements
	err := sp.RequireNodeAddress()
	if err != nil {
		return types.ResponseStatus_AddressNotPresent, err
	}

	data.Balance, err = ec.BalanceAt(ctx, nodeAddress, nil)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting ETH balance of node %s: %w", nodeAddress.Hex(), err)
	}
	return types.ResponseStatus_Success, nil
}

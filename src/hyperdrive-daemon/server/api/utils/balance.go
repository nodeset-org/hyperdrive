package utils

import (
	"context"
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
	utils.RegisterQuerylessGet[*utilsBalanceContext, api.UtilsBalanceData](
		router, "balance", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type utilsBalanceContext struct {
	handler *UtilsHandler
}

func (c *utilsBalanceContext) PrepareData(data *api.UtilsBalanceData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	ec := sp.GetEthClient()
	nodeAddress, _ := sp.GetWallet().GetAddress()

	// Requirements
	err := sp.RequireNodeAddress()
	if err != nil {
		return err
	}

	data.Balance, err = ec.BalanceAt(context.Background(), nodeAddress, nil)
	if err != nil {
		return fmt.Errorf("error getting ETH balance of node %s: %w", nodeAddress.Hex(), err)
	}
	return nil
}

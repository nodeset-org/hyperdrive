package swwallet

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
	api "github.com/nodeset-org/hyperdrive/shared/types/api/modules/stakewise"
)

// ===============
// === Factory ===
// ===============

type walletRegenerateDepositDataContextFactory struct {
	handler *WalletHandler
}

func (f *walletRegenerateDepositDataContextFactory) Create(args url.Values) (*walletRegenerateDepositDataContext, error) {
	c := &walletRegenerateDepositDataContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletRegenerateDepositDataContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletRegenerateDepositDataContext, api.WalletRegenerateDepositDataData](
		router, "regen-deposit-data", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletRegenerateDepositDataContext struct {
	handler *WalletHandler
}

func (c *walletRegenerateDepositDataContext) PrepareData(data *api.WalletRegenerateDepositDataData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	ddMgr := sp.GetDepositDataManager()

	// Regen the deposit data file
	pubkeys, err := ddMgr.RegenerateDepositData()
	if err != nil {
		return fmt.Errorf("error regenerating deposit data: %w", err)
	}
	data.Pubkeys = pubkeys
	return nil
}

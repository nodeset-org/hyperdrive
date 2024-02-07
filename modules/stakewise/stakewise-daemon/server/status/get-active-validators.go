package swstatus

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type statusGetActiveValidatorsContextFactory struct {
	handler *StatusHandler
}

func (f *statusGetActiveValidatorsContextFactory) Create(args url.Values) (*statusGetActiveValidatorsContext, error) {
	c := &statusGetActiveValidatorsContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("root", args, input.ValidateHash, &c.root),
	}
	return c, errors.Join(inputErrs...)
}

func (f *statusGetActiveValidatorsContextFactory) RegisterRoute(router *mux.Router) {
	// TODO: Should I be using SuccessData here???
	server.RegisterQuerylessGet[*statusGetActiveValidatorsContext, api.SuccessData](
		router, "get-active-validators", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type statusGetActiveValidatorsContext struct {
	handler *StatusHandler
	root    common.Hash
}

func (c *statusGetActiveValidatorsContext) PrepareData(data *api.SuccessData, opts *bind.TransactOpts) error {
	fmt.Printf("statusGetActiveValidatorsContext.PrepareData data: %+v\n", data)
	// sp := c.handler.serviceProvider
	// ec := sp.GetEthClient()
	// res := sp.GetResources()
	// txMgr := sp.GetTransactionManager()

	// vault, err := swcommon.NewStakewiseVault(res.Vault, ec, txMgr)
	// if err != nil {
	// 	return fmt.Errorf("error creating Stakewise Vault binding: %w", err)
	// }

	// data.TxInfo, err = vault.SetDepositDataRoot(c.root, opts)
	// if err != nil {
	// 	return fmt.Errorf("error creating SetDepositDataRoot TX: %w", err)
	// }
	return nil
}

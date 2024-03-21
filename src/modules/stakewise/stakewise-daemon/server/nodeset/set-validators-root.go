package swnodeset

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swcontracts "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common/contracts"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type nodesetSetValidatorsRootContextFactory struct {
	handler *NodesetHandler
}

func (f *nodesetSetValidatorsRootContextFactory) Create(args url.Values) (*nodesetSetValidatorsRootContext, error) {
	c := &nodesetSetValidatorsRootContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("root", args, input.ValidateHash, &c.root),
	}
	return c, errors.Join(inputErrs...)
}

func (f *nodesetSetValidatorsRootContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*nodesetSetValidatorsRootContext, api.TxInfoData](
		router, "set-validators-root", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type nodesetSetValidatorsRootContext struct {
	handler *NodesetHandler
	root    common.Hash
}

func (c *nodesetSetValidatorsRootContext) PrepareData(data *api.TxInfoData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	ec := sp.GetEthClient()
	res := sp.GetResources()
	txMgr := sp.GetTransactionManager()

	vault, err := swcontracts.NewStakewiseVault(res.Vault, ec, txMgr)
	if err != nil {
		return fmt.Errorf("error creating Stakewise Vault binding: %w", err)
	}

	data.TxInfo, err = vault.SetDepositDataRoot(c.root, opts)
	if err != nil {
		return fmt.Errorf("error creating SetDepositDataRoot TX: %w", err)
	}
	return nil
}

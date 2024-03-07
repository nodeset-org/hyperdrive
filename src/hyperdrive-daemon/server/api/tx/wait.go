package tx

import (
	"errors"
	"fmt"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type txWaitContextFactory struct {
	handler *TxHandler
}

func (f *txWaitContextFactory) Create(args url.Values) (*txWaitContext, error) {
	c := &txWaitContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("hash", args, input.ValidateHash, &c.hash),
	}
	return c, errors.Join(inputErrs...)
}

func (f *txWaitContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*txWaitContext, nmc_types.SuccessData](
		router, "wait", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type txWaitContext struct {
	handler *TxHandler
	hash    common.Hash
}

func (c *txWaitContext) PrepareData(data *nmc_types.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	txMgr := sp.GetTransactionManager()

	err := txMgr.WaitForTransactionByHash(c.hash)
	if err != nil {
		return fmt.Errorf("error waiting for tx %s: %w", c.hash.Hex(), err)
	}
	return nil
}

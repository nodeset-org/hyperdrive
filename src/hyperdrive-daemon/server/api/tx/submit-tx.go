package tx

import (
	"errors"
	"fmt"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type txSubmitTxContextFactory struct {
	handler *TxHandler
}

func (f *txSubmitTxContextFactory) Create(body api.SubmitTxBody) (*txSubmitTxContext, error) {
	c := &txSubmitTxContext{
		handler: f.handler,
		body:    body,
	}
	// Validate the submission
	if body.Submission.TxInfo == nil {
		return nil, fmt.Errorf("submission TX info must be set")
	}
	if body.Submission.GasLimit == 0 {
		return nil, fmt.Errorf("submission gas limit must be set")
	}
	if body.MaxFee == nil {
		return nil, fmt.Errorf("submission max fee must be set")
	}
	if body.MaxPriorityFee == nil {
		return nil, fmt.Errorf("submission max priority fee must be set")
	}
	return c, nil
}

func (f *txSubmitTxContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessPost[*txSubmitTxContext, api.SubmitTxBody, api.TxData](
		router, "submit-tx", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type txSubmitTxContext struct {
	handler *TxHandler
	body    api.SubmitTxBody
}

func (c *txSubmitTxContext) PrepareData(data *api.TxData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	txMgr := sp.GetTransactionManager()

	err := errors.Join(
		sp.RequireWalletReady(),
	)
	if err != nil {
		return types.ResponseStatus_WalletNotReady, err
	}

	if c.body.Nonce != nil {
		opts.Nonce = c.body.Nonce
	}
	opts.GasLimit = c.body.Submission.GasLimit
	opts.GasFeeCap = c.body.MaxFee
	opts.GasTipCap = c.body.MaxPriorityFee

	tx, err := txMgr.ExecuteTransaction(c.body.Submission.TxInfo, opts)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error submitting transaction: %w", err)
	}
	data.TxHash = tx.Hash()
	return types.ResponseStatus_Success, nil
}

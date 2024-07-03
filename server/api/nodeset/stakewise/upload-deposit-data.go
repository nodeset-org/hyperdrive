package ns_stakewise

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/rocket-pool/node-manager-core/beacon"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type stakeWiseUploadDepositDataContextFactory struct {
	handler *StakeWiseHandler
}

func (f *stakeWiseUploadDepositDataContextFactory) Create(body []beacon.ExtendedDepositData) (*stakeWiseUploadDepositDataContext, error) {
	c := &stakeWiseUploadDepositDataContext{
		handler: f.handler,
		body:    body,
	}
	return c, nil
}

func (f *stakeWiseUploadDepositDataContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessPost[*stakeWiseUploadDepositDataContext, []beacon.ExtendedDepositData, types.SuccessData](
		router, "upload-deposit-data", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============
type stakeWiseUploadDepositDataContext struct {
	handler *StakeWiseHandler
	body    []beacon.ExtendedDepositData
}

func (c *stakeWiseUploadDepositDataContext) PrepareData(data *types.SuccessData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ctx := c.handler.ctx

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return types.ResponseStatus_WalletNotReady, err
	}
	err = sp.RequireRegisteredWithNodeSet(ctx)
	if err != nil {
		return types.ResponseStatus_Error, err
	}

	// Upload the deposit data
	ns := sp.GetNodeSetServiceManager()
	err = ns.StakeWise_UploadDepositData(ctx, c.body)
	if err != nil {
		return types.ResponseStatus_Error, err
	}
	return types.ResponseStatus_Success, nil
}

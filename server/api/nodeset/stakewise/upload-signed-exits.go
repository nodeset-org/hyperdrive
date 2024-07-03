package ns_stakewise

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type stakeWiseUploadSignedExitsContextFactory struct {
	handler *StakeWiseHandler
}

func (f *stakeWiseUploadSignedExitsContextFactory) Create(body []api.ExitData) (*stakeWiseUploadSignedExitsContext, error) {
	c := &stakeWiseUploadSignedExitsContext{
		handler: f.handler,
		body:    body,
	}
	return c, nil
}

func (f *stakeWiseUploadSignedExitsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessPost[*stakeWiseUploadSignedExitsContext, []api.ExitData, types.SuccessData](
		router, "upload-signed-exits", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============
type stakeWiseUploadSignedExitsContext struct {
	handler *StakeWiseHandler
	body    []api.ExitData
}

func (c *stakeWiseUploadSignedExitsContext) PrepareData(data *types.SuccessData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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
	err = ns.StakeWise_UploadSignedExitMessages(ctx, c.body)
	if err != nil {
		return types.ResponseStatus_Error, err
	}
	return types.ResponseStatus_Success, nil
}

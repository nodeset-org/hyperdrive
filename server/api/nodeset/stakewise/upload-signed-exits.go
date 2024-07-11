package ns_stakewise

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	apiv1 "github.com/nodeset-org/nodeset-client-go/api-v1"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type stakeWiseUploadSignedExitsContextFactory struct {
	handler *StakeWiseHandler
}

func (f *stakeWiseUploadSignedExitsContextFactory) Create(body []apiv1.ExitData) (*stakeWiseUploadSignedExitsContext, error) {
	c := &stakeWiseUploadSignedExitsContext{
		handler: f.handler,
		body:    body,
	}
	return c, nil
}

func (f *stakeWiseUploadSignedExitsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessPost[*stakeWiseUploadSignedExitsContext, []apiv1.ExitData, api.NodeSetStakeWise_UploadSignedExitsData](
		router, "upload-signed-exits", f, f.handler.logger.Logger, f.handler.serviceProvider.IServiceProvider,
	)
}

// ===============
// === Context ===
// ===============
type stakeWiseUploadSignedExitsContext struct {
	handler *StakeWiseHandler
	body    []apiv1.ExitData
}

func (c *stakeWiseUploadSignedExitsContext) PrepareData(data *api.NodeSetStakeWise_UploadSignedExitsData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ctx := c.handler.ctx

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return types.ResponseStatus_WalletNotReady, err
	}
	err = sp.RequireRegisteredWithNodeSet(ctx)
	if err != nil {
		if errors.Is(err, hdcommon.ErrNotRegisteredWithNodeSet) {
			data.NotRegistered = true
			return types.ResponseStatus_Success, nil
		}
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

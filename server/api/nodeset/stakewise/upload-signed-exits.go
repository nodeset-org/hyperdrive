package ns_stakewise

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
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

func (f *stakeWiseUploadSignedExitsContextFactory) Create(body api.NodeSetStakeWise_UploadSignedExitsRequestBody) (*stakeWiseUploadSignedExitsContext, error) {
	c := &stakeWiseUploadSignedExitsContext{
		handler: f.handler,
		body:    body,
	}
	return c, nil
}

func (f *stakeWiseUploadSignedExitsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessPost[*stakeWiseUploadSignedExitsContext, api.NodeSetStakeWise_UploadSignedExitsRequestBody, api.NodeSetStakeWise_UploadSignedExitsData](
		router, "upload-signed-exits", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============
type stakeWiseUploadSignedExitsContext struct {
	handler *StakeWiseHandler
	body    api.NodeSetStakeWise_UploadSignedExitsRequestBody
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
	err = ns.StakeWise_UploadSignedExitMessages(ctx, c.body.Deployment, c.body.Vault, c.body.ExitData)
	if err != nil {
		return types.ResponseStatus_Error, err
	}
	return types.ResponseStatus_Success, nil
}

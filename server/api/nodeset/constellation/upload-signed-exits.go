package ns_constellation

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	v2constellation "github.com/nodeset-org/nodeset-client-go/api-v2/constellation"
	"github.com/nodeset-org/nodeset-client-go/common"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type constellationUploadSignedExitsContextFactory struct {
	handler *ConstellationHandler
}

func (f *constellationUploadSignedExitsContextFactory) Create(body api.NodeSetConstellation_UploadSignedExitsRequestBody) (*constellationUploadSignedExitsContext, error) {
	c := &constellationUploadSignedExitsContext{
		handler: f.handler,
		body:    body,
	}
	return c, nil
}

func (f *constellationUploadSignedExitsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessPost[*constellationUploadSignedExitsContext, api.NodeSetConstellation_UploadSignedExitsRequestBody, api.NodeSetConstellation_UploadSignedExitsData](
		router, "upload-signed-exits", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type constellationUploadSignedExitsContext struct {
	handler *ConstellationHandler
	body    api.NodeSetConstellation_UploadSignedExitsRequestBody
}

func (c *constellationUploadSignedExitsContext) PrepareData(data *api.NodeSetConstellation_UploadSignedExitsData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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
	err = ns.Constellation_UploadSignedExitMessages(ctx, c.body.Deployment, c.body.ExitMessages)
	if err != nil {
		if errors.Is(err, v2constellation.ErrMissingWhitelistedNodeAddress) {
			data.NotWhitelisted = true
			return types.ResponseStatus_Success, nil
		}
		if errors.Is(err, v2constellation.ErrIncorrectNodeAddress) {
			data.IncorrectNodeAddress = true
			return types.ResponseStatus_Success, nil
		}
		if errors.Is(err, common.ErrInvalidValidatorOwner) {
			data.InvalidValidatorOwner = true
			return types.ResponseStatus_Success, nil
		}
		if errors.Is(err, v2constellation.ErrExitMessageExists) {
			data.ExitMessageAlreadyExists = true
			return types.ResponseStatus_Success, nil
		}
		if errors.Is(err, common.ErrInvalidExitMessage) {
			data.InvalidExitMessage = true
			return types.ResponseStatus_Success, nil
		}
		if errors.Is(err, common.ErrInvalidPermissions) {
			data.InvalidPermissions = true
			return types.ResponseStatus_Success, nil
		}
		return types.ResponseStatus_Error, err
	}
	return types.ResponseStatus_Success, nil
}

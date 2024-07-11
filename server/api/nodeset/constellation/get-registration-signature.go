package ns_constellation

import (
	"errors"
	"net/url"

	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	apiv2 "github.com/nodeset-org/nodeset-client-go/api-v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type constellationGetRegistrationSignatureContextFactory struct {
	handler *ConstellationHandler
}

func (f *constellationGetRegistrationSignatureContextFactory) Create(args url.Values) (*constellationGetRegistrationSignatureContext, error) {
	c := &constellationGetRegistrationSignatureContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *constellationGetRegistrationSignatureContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*constellationGetRegistrationSignatureContext, api.NodeSetConstellation_GetRegistrationSignatureData](
		router, "get-registration-signature", f, f.handler.logger.Logger, f.handler.serviceProvider.IServiceProvider,
	)
}

// ===============
// === Context ===
// ===============
type constellationGetRegistrationSignatureContext struct {
	handler *ConstellationHandler
}

func (c *constellationGetRegistrationSignatureContext) PrepareData(data *api.NodeSetConstellation_GetRegistrationSignatureData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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

	// Get the registration signature
	ns := sp.GetNodeSetServiceManager()
	signature, err := ns.Constellation_GetRegistrationSignature(ctx)
	if err != nil {
		if errors.Is(err, apiv2.ErrNotAuthorized) {
			data.NotAuthorized = true
			return types.ResponseStatus_Success, nil
		}
		return types.ResponseStatus_Error, err
	}

	data.Signature = signature
	return types.ResponseStatus_Success, nil
}

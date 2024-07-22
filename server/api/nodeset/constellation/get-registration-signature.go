package ns_constellation

import (
	"errors"
	"net/url"

	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	apiv2 "github.com/nodeset-org/nodeset-client-go/api-v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/utils/input"
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
	inputErrs := []error{
		server.ValidateArg("whitelistAddress", args, input.ValidateAddress, &c.whitelistAddress),
	}
	return c, errors.Join(inputErrs...)
}

func (f *constellationGetRegistrationSignatureContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*constellationGetRegistrationSignatureContext, api.NodeSetConstellation_GetRegistrationSignatureData](
		router, "get-registration-signature", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============
type constellationGetRegistrationSignatureContext struct {
	handler *ConstellationHandler

	whitelistAddress common.Address
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
	timestamp, signature, err := ns.Constellation_GetRegistrationSignatureAndTime(ctx, c.whitelistAddress)
	if err != nil {
		if errors.Is(err, apiv2.ErrNotAuthorized) {
			data.NotAuthorized = true
			return types.ResponseStatus_Success, nil
		}
		return types.ResponseStatus_Error, err
	}

	data.Signature = signature
	data.Time = timestamp
	return types.ResponseStatus_Success, nil
}

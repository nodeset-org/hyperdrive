package ns_constellation

import (
	"errors"
	"net/url"

	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/nodeset-org/nodeset-client-go/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type constellationGetRegisteredAddressContextFactory struct {
	handler *ConstellationHandler
}

func (f *constellationGetRegisteredAddressContextFactory) Create(args url.Values) (*constellationGetRegisteredAddressContext, error) {
	c := &constellationGetRegisteredAddressContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.GetStringFromVars("deployment", args, &c.deployment),
	}
	return c, errors.Join(inputErrs...)
}

func (f *constellationGetRegisteredAddressContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*constellationGetRegisteredAddressContext, api.NodeSetConstellation_GetRegisteredAddressData](
		router, "get-registered-address", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type constellationGetRegisteredAddressContext struct {
	handler *ConstellationHandler

	deployment string
}

func (c *constellationGetRegisteredAddressContext) PrepareData(data *api.NodeSetConstellation_GetRegisteredAddressData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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
			data.NotRegisteredWithNodeSet = true
			return types.ResponseStatus_Success, nil
		}
		return types.ResponseStatus_Error, err
	}

	// Get the registered address
	ns := sp.GetNodeSetServiceManager()
	address, err := ns.Constellation_GetRegisteredAddress(ctx, c.deployment)
	if err != nil {
		if errors.Is(err, common.ErrInvalidPermissions) {
			data.InvalidPermissions = true
			return types.ResponseStatus_Success, nil
		}
		return types.ResponseStatus_Error, err
	}

	if address == nil {
		data.NotRegisteredWithConstellation = true
	} else {
		data.RegisteredAddress = *address
	}
	return types.ResponseStatus_Success, nil
}

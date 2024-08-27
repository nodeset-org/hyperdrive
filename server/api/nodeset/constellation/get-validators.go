package ns_constellation

import (
	"errors"
	"net/url"

	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	v2constellation "github.com/nodeset-org/nodeset-client-go/api-v2/constellation"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type constellationGetValidatorsContextFactory struct {
	handler *ConstellationHandler
}

func (f *constellationGetValidatorsContextFactory) Create(args url.Values) (*constellationGetValidatorsContext, error) {
	c := &constellationGetValidatorsContext{
		handler: f.handler,
	}
	inputErrs := []error{}
	return c, errors.Join(inputErrs...)
}

func (f *constellationGetValidatorsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*constellationGetValidatorsContext, api.NodeSetConstellation_GetValidatorsData](
		router, "get-validators", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type constellationGetValidatorsContext struct {
	handler *ConstellationHandler
}

func (c *constellationGetValidatorsContext) PrepareData(data *api.NodeSetConstellation_GetValidatorsData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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

	// Get the registered validators
	ns := sp.GetNodeSetServiceManager()
	validators, err := ns.Constellation_GetValidators(ctx)
	if err != nil {
		if errors.Is(err, v2constellation.ErrNotAuthorized) {
			data.NotAuthorized = true
			return types.ResponseStatus_Success, nil
		}
		return types.ResponseStatus_Error, err
	}

	data.Validators = validators
	return types.ResponseStatus_Success, nil
}

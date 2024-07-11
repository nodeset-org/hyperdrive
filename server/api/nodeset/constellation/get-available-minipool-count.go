package ns_constellation

import (
	"errors"
	"net/url"

	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type constellationGetAvailableMinipoolCountContextFactory struct {
	handler *ConstellationHandler
}

func (f *constellationGetAvailableMinipoolCountContextFactory) Create(args url.Values) (*constellationGetAvailableMinipoolCountContext, error) {
	c := &constellationGetAvailableMinipoolCountContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *constellationGetAvailableMinipoolCountContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*constellationGetAvailableMinipoolCountContext, api.NodeSetConstellation_GetAvailableMinipoolCount](
		router, "get-available-minipool-count", f, f.handler.logger.Logger, f.handler.serviceProvider.IServiceProvider,
	)
}

// ===============
// === Context ===
// ===============
type constellationGetAvailableMinipoolCountContext struct {
	handler *ConstellationHandler
}

func (c *constellationGetAvailableMinipoolCountContext) PrepareData(data *api.NodeSetConstellation_GetAvailableMinipoolCount, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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

	// Get the set version
	ns := sp.GetNodeSetServiceManager()
	count, err := ns.Constellation_GetAvailableMinipoolCount(ctx)
	if err != nil {
		return types.ResponseStatus_Error, err
	}

	data.Count = count
	return types.ResponseStatus_Success, nil
}

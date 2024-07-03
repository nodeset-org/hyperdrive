package nodeset

import (
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type nodeSetGetRegistrationStatusContextFactory struct {
	handler *NodeSetHandler
}

func (f *nodeSetGetRegistrationStatusContextFactory) Create(args url.Values) (*nodeSetGetRegistrationStatusContext, error) {
	c := &nodeSetGetRegistrationStatusContext{
		handler: f.handler,
	}

	return c, nil
}

func (f *nodeSetGetRegistrationStatusContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*nodeSetGetRegistrationStatusContext, api.NodeSetGetRegistrationStatusData](
		router, "get-registration-status", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type nodeSetGetRegistrationStatusContext struct {
	handler *NodeSetHandler
}

func (c *nodeSetGetRegistrationStatusContext) PrepareData(data *api.NodeSetGetRegistrationStatusData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ctx := c.handler.ctx

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return types.ResponseStatus_WalletNotReady, err
	}

	// Register the node
	ns := sp.GetNodeSetServiceManager()
	data.Status, err = ns.GetRegistrationStatus(ctx)
	if err != nil {
		return types.ResponseStatus_Error, err
	}

	return types.ResponseStatus_Success, nil
}

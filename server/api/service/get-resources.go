package service

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

type serviceGetResourcesContextFactory struct {
	handler *ServiceHandler
}

func (f *serviceGetResourcesContextFactory) Create(args url.Values) (*serviceGetResourcesContext, error) {
	c := &serviceGetResourcesContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *serviceGetResourcesContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*serviceGetResourcesContext, api.ServiceGetResourcesData](
		router, "get-resources", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type serviceGetResourcesContext struct {
	handler *ServiceHandler
}

func (c *serviceGetResourcesContext) PrepareData(data *api.ServiceGetResourcesData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	data.Resources = sp.GetResources()
	return types.ResponseStatus_Success, nil
}

package service

import (
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type serviceVersionContextFactory struct {
	handler *ServiceHandler
}

func (f *serviceVersionContextFactory) Create(args url.Values) (*serviceVersionContext, error) {
	c := &serviceVersionContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *serviceVersionContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*serviceVersionContext, api.ServiceVersionData](
		router, "version", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type serviceVersionContext struct {
	handler *ServiceHandler
}

func (c *serviceVersionContext) PrepareData(data *api.ServiceVersionData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	data.Version = shared.HyperdriveVersion
	return types.ResponseStatus_Success, nil
}

package service

import (
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
)

// ===============
// === Factory ===
// ===============

type serviceGetConfigContextFactory struct {
	handler *ServiceHandler
}

func (f *serviceGetConfigContextFactory) Create(args url.Values) (*serviceGetConfigContext, error) {
	c := &serviceGetConfigContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *serviceGetConfigContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*serviceGetConfigContext, api.ServiceGetConfigData](
		router, "get-config", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type serviceGetConfigContext struct {
	handler *ServiceHandler
}

func (c *serviceGetConfigContext) PrepareData(data *api.ServiceGetConfigData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	cfg := sp.GetConfig()

	data.Config = cfg.Serialize(nil)
	return nil
}

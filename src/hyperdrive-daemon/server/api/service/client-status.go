package service

import (
	"context"
	"net/url"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
)

// ===============
// === Factory ===
// ===============

type serviceClientStatusContextFactory struct {
	handler *ServiceHandler
}

func (f *serviceClientStatusContextFactory) Create(args url.Values) (*serviceClientStatusContext, error) {
	c := &serviceClientStatusContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *serviceClientStatusContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*serviceClientStatusContext, api.ServiceClientStatusData](
		router, "client-status", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type serviceClientStatusContext struct {
	handler *ServiceHandler
}

func (c *serviceClientStatusContext) PrepareData(data *api.ServiceClientStatusData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	ec := sp.GetEthClient()
	bc := sp.GetBeaconClient()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Get the EC manager status
	go func() {
		ecMgrStatus := ec.CheckStatus(context.Background())
		data.EcManagerStatus = *ecMgrStatus
		wg.Done()
	}()

	// Get the BC manager status
	go func() {
		bcMgrStatus := bc.CheckStatus(context.Background())
		data.BcManagerStatus = *bcMgrStatus
		wg.Done()
	}()

	wg.Wait()
	return nil
}

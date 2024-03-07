package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/docker/docker/api/types/container"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type serviceRestartContainerContextFactory struct {
	handler *ServiceHandler
}

func (f *serviceRestartContainerContextFactory) Create(args url.Values) (*serviceRestartContainerContext, error) {
	c := &serviceRestartContainerContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.GetStringFromVars("container", args, &c.container),
	}
	return c, errors.Join(inputErrs...)
}

func (f *serviceRestartContainerContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*serviceRestartContainerContext, nmc_types.SuccessData](
		router, "restart-container", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type serviceRestartContainerContext struct {
	handler   *ServiceHandler
	container string
}

func (c *serviceRestartContainerContext) PrepareData(data *nmc_types.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	cfg := sp.GetConfig()
	d := sp.GetDocker()

	id := cfg.GetDockerArtifactName(c.container)
	return d.ContainerRestart(context.Background(), id, container.StopOptions{})
}

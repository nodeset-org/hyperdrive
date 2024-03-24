package service

import (
	"errors"
	"net/url"

	"github.com/docker/docker/api/types/container"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
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
	server.RegisterQuerylessGet[*serviceRestartContainerContext, types.SuccessData](
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

func (c *serviceRestartContainerContext) PrepareData(data *types.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	cfg := sp.GetConfig()
	d := sp.GetDocker()
	ctx := sp.GetContext()

	id := cfg.GetDockerArtifactName(c.container)
	return d.ContainerRestart(ctx, id, container.StopOptions{})
}

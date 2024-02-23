package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/docker/docker/api/types/container"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
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
	utils.RegisterQuerylessGet[*serviceRestartContainerContext, api.SuccessData](
		router, "restart-container", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type serviceRestartContainerContext struct {
	handler   *ServiceHandler
	container string
}

func (c *serviceRestartContainerContext) PrepareData(data *api.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	cfg := sp.GetConfig()
	d := sp.GetDocker()

	id := cfg.GetDockerArtifactName(c.container)
	return d.ContainerRestart(context.Background(), id, container.StopOptions{})
}

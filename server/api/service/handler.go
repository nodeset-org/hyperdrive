package service

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/log"
)

type ServiceHandler struct {
	logger          *log.Logger
	ctx             context.Context
	serviceProvider *common.ServiceProvider
	factories       []server.IContextFactory
}

func NewServiceHandler(logger *log.Logger, ctx context.Context, serviceProvider *common.ServiceProvider) *ServiceHandler {
	h := &ServiceHandler{
		logger:          logger,
		ctx:             ctx,
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&serviceClientStatusContextFactory{h},
		&serviceGetConfigContextFactory{h},
		&serviceRestartContainerContextFactory{h},
		&serviceRotateLogsContextFactory{h},
		&serviceVersionContextFactory{h},
	}
	return h
}

func (h *ServiceHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/service").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

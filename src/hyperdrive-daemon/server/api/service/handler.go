package service

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
)

type ServiceHandler struct {
	serviceProvider *common.ServiceProvider
	factories       []nmc_server.IContextFactory
}

func NewServiceHandler(serviceProvider *common.ServiceProvider) *ServiceHandler {
	h := &ServiceHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []nmc_server.IContextFactory{
		&serviceClientStatusContextFactory{h},
		&serviceGetConfigContextFactory{h},
		&serviceRestartContainerContextFactory{h},
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

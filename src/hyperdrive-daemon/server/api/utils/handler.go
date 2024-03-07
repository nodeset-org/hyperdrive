package utils

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
)

type UtilsHandler struct {
	serviceProvider *common.ServiceProvider
	factories       []nmc_server.IContextFactory
}

func NewUtilsHandler(serviceProvider *common.ServiceProvider) *UtilsHandler {
	h := &UtilsHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []nmc_server.IContextFactory{
		&utilsBalanceContextFactory{h},
		&utilsResolveEnsContextFactory{h},
	}
	return h
}

func (h *UtilsHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/utils").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

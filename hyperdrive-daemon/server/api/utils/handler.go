package utils

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
)

type UtilsHandler struct {
	serviceProvider *common.ServiceProvider
	factories       []server.IContextFactory
}

func NewUtilsHandler(serviceProvider *common.ServiceProvider) *UtilsHandler {
	h := &UtilsHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
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

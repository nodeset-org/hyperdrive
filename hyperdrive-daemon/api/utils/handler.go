package utils

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/services"
)

type UtilsHandler struct {
	serviceProvider *services.ServiceProvider
	factories       []server.IContextFactory
}

func NewUtilsHandler(serviceProvider *services.ServiceProvider) *UtilsHandler {
	h := &UtilsHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
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

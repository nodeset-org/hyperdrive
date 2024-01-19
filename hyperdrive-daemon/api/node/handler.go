package example

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/common/services"
)

type NodeHandler struct {
	serviceProvider *services.ServiceProvider
	factories       []server.IContextFactory
	isDebug         bool
}

func NewNodeHandler(serviceProvider *services.ServiceProvider, isDebug bool) *NodeHandler {
	h := &NodeHandler{
		serviceProvider: serviceProvider,
		isDebug:         isDebug,
	}
	h.factories = []server.IContextFactory{
		&uploadDepositDataContextFactory{h},
	}
	return h
}

func (h *NodeHandler) RegisterRoutes(router *mux.Router, debugMode bool) {
	subrouter := router.PathPrefix("/node").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

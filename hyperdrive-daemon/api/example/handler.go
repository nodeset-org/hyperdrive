package example

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/services"
)

type ExampleHandler struct {
	serviceProvider *services.ServiceProvider
	factories       []server.IContextFactory
	isDebug         bool
}

func NewExampleHandler(serviceProvider *services.ServiceProvider, isDebug bool) *ExampleHandler {
	h := &ExampleHandler{
		serviceProvider: serviceProvider,
		isDebug:         isDebug,
	}
	h.factories = []server.IContextFactory{
		&callDaemonContextFactory{h},
	}
	return h
}

func (h *ExampleHandler) RegisterRoutes(router *mux.Router, debugMode bool) {
	subrouter := router.PathPrefix("/example").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

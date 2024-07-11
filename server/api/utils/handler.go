package utils

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/log"
)

type UtilsHandler struct {
	logger          *log.Logger
	ctx             context.Context
	serviceProvider *common.ServiceProvider
	factories       []server.IContextFactory
}

func NewUtilsHandler(logger *log.Logger, ctx context.Context, serviceProvider *common.ServiceProvider) *UtilsHandler {
	h := &UtilsHandler{
		logger:          logger,
		ctx:             ctx,
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

package ns_constellation

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/log"
)

type ConstellationHandler struct {
	logger          *log.Logger
	ctx             context.Context
	serviceProvider common.IHyperdriveServiceProvider
	factories       []server.IContextFactory
}

func NewConstellationHandler(logger *log.Logger, ctx context.Context, serviceProvider common.IHyperdriveServiceProvider) *ConstellationHandler {
	h := &ConstellationHandler{
		logger:          logger,
		ctx:             ctx,
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&constellationGetAvailableMinipoolCountContextFactory{h},
		&constellationGetDepositSignatureContextFactory{h},
		&constellationGetRegistrationSignatureContextFactory{h},
	}
	return h
}

func (h *ConstellationHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/constellation").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

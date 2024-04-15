package swnodeset

import (
	"context"

	"github.com/gorilla/mux"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/log"
)

type NodesetHandler struct {
	logger          *log.Logger
	ctx             context.Context
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewNodesetHandler(logger *log.Logger, ctx context.Context, serviceProvider *swcommon.StakewiseServiceProvider) *NodesetHandler {
	h := &NodesetHandler{
		logger:          logger,
		ctx:             ctx,
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&nodesetSetValidatorsRootContextFactory{h},
		&nodesetUploadDepositDataContextFactory{h},
	}
	return h
}

func (h *NodesetHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/nodeset").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

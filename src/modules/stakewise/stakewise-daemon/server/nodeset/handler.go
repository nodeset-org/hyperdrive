package swnodeset

import (
	"github.com/gorilla/mux"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
)

type NodesetHandler struct {
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewNodesetHandler(serviceProvider *swcommon.StakewiseServiceProvider) *NodesetHandler {
	h := &NodesetHandler{
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

package swnodeset

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
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

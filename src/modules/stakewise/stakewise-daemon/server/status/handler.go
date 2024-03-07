package swstatus

import (
	"github.com/gorilla/mux"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
)

type StatusHandler struct {
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewStatusHandler(serviceProvider *swcommon.StakewiseServiceProvider) *StatusHandler {
	h := &StatusHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&statusGetActiveValidatorsContextFactory{h},
	}
	return h
}

func (h *StatusHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/status").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

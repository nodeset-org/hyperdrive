package swvalidator

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
)

type ValidatorHandler struct {
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewValidatorHandler(serviceProvider *swcommon.StakewiseServiceProvider) *ValidatorHandler {
	h := &ValidatorHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&validatorExitContextFactory{h},
	}
	return h
}

func (h *ValidatorHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/validator").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

package swvalidator

import (
	"github.com/gorilla/mux"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
)

type ValidatorHandler struct {
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []nmc_server.IContextFactory
}

func NewValidatorHandler(serviceProvider *swcommon.StakewiseServiceProvider) *ValidatorHandler {
	h := &ValidatorHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []nmc_server.IContextFactory{
		&validatorGetSignedExitMessagesContextFactory{h},
	}
	return h
}

func (h *ValidatorHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/validator").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

package swvalidator

import (
	"context"

	"github.com/gorilla/mux"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/log"
)

type ValidatorHandler struct {
	logger          *log.Logger
	ctx             context.Context
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewValidatorHandler(logger *log.Logger, ctx context.Context, serviceProvider *swcommon.StakewiseServiceProvider) *ValidatorHandler {
	h := &ValidatorHandler{
		logger:          logger,
		ctx:             ctx,
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

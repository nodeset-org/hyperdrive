package swstatus

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
)

type StatusHandler struct {
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewStatusHandler(serviceProvider *swcommon.StakewiseServiceProvider) *StatusHandler {
	h := &StatusHandler{
		serviceProvider: serviceProvider,
	}
	// h.factories = []server.IContextFactory{
	// 	&nodesetSetValidatorsRootContextFactory{h},
	// 	&nodesetUploadDepositDataContextFactory{h},
	// }
	return h
}

func (h *StatusHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/status").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

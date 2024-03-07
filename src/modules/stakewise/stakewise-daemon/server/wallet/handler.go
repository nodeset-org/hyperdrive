package swwallet

import (
	"github.com/gorilla/mux"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
)

type WalletHandler struct {
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []nmc_server.IContextFactory
}

func NewWalletHandler(serviceProvider *swcommon.StakewiseServiceProvider) *WalletHandler {
	h := &WalletHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []nmc_server.IContextFactory{
		&walletGenerateKeysContextFactory{h},
		&walletInitializeContextFactory{h},
	}
	return h
}

func (h *WalletHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/wallet").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

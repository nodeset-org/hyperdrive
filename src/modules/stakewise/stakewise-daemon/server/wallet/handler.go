package swwallet

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
)

type WalletHandler struct {
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewWalletHandler(serviceProvider *swcommon.StakewiseServiceProvider) *WalletHandler {
	h := &WalletHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&walletClaimRewardsContextFactory{h},
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

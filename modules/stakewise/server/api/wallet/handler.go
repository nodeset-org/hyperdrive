package swwallet

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
	"github.com/nodeset-org/hyperdrive/modules/stakewise/common"
)

type WalletHandler struct {
	serviceProvider *common.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewWalletHandler(serviceProvider *common.StakewiseServiceProvider) *WalletHandler {
	h := &WalletHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
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

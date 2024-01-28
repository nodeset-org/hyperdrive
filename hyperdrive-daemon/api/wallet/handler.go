package wallet

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/services"
)

type WalletHandler struct {
	serviceProvider *services.ServiceProvider
	factories       []server.IContextFactory
}

func NewWalletHandler(serviceProvider *services.ServiceProvider) *WalletHandler {
	h := &WalletHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&walletDeletePasswordContextFactory{h},
		&walletExportContextFactory{h},
		&walletExportEthKeyContextFactory{h},
		&walletInitializeContextFactory{h},
		&walletMasqueradeContextFactory{h},
		&walletRecoverContextFactory{h},
		&walletRestoreAddressContextFactory{h},
		&walletSearchAndRecoverContextFactory{h},
		&walletSendMessageContextFactory{h},
		&walletSetEnsNameContextFactory{h},
		&walletSetPasswordContextFactory{h},
		&walletSignMessageContextFactory{h},
		&walletStatusFactory{h},
		&walletTestRecoverContextFactory{h},
		&walletTestSearchAndRecoverContextFactory{h},
	}
	return h
}

func (h *WalletHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/wallet").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

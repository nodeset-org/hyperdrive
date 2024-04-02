package swwallet

import (
	"context"

	"github.com/gorilla/mux"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/log"
)

type WalletHandler struct {
	logger          *log.Logger
	ctx             context.Context
	serviceProvider *swcommon.StakewiseServiceProvider
	factories       []server.IContextFactory
}

func NewWalletHandler(logger *log.Logger, ctx context.Context, serviceProvider *swcommon.StakewiseServiceProvider) *WalletHandler {
	h := &WalletHandler{
		logger:          logger,
		ctx:             ctx,
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

package tx

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/log"
)

type TxHandler struct {
	logger          *log.Logger
	ctx             context.Context
	serviceProvider common.IHyperdriveServiceProvider
	factories       []server.IContextFactory
}

func NewTxHandler(logger *log.Logger, ctx context.Context, serviceProvider common.IHyperdriveServiceProvider) *TxHandler {
	h := &TxHandler{
		logger:          logger,
		ctx:             ctx,
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&txBatchSignTxsContextFactory{h},
		&txBatchSubmitTxsContextFactory{h},
		&txSignTxContextFactory{h},
		&txSubmitTxContextFactory{h},
		&txWaitContextFactory{h},
	}
	return h
}

func (h *TxHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/tx").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

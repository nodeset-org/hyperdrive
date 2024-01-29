package tx

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
)

type TxHandler struct {
	serviceProvider *common.ServiceProvider
	factories       []server.IContextFactory
}

func NewTxHandler(serviceProvider *common.ServiceProvider) *TxHandler {
	h := &TxHandler{
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

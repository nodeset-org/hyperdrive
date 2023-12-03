package example

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/services"
)

// Context factories can implement this generally so they can register themselves with an HTTP router.
type IContextFactory interface {
	RegisterRoute(router *mux.Router)
}

type ExampleHandler struct {
	serviceProvider *services.ServiceProvider
	factories       []IContextFactory
}

func NewAuctionHandler(serviceProvider *services.ServiceProvider) *ExampleHandler {
	h := &ExampleHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []IContextFactory{
		// &auctionBidContextFactory{h},
		// &auctionClaimContextFactory{h},
		// &auctionCreateContextFactory{h},
		// &auctionLotContextFactory{h},
		// &auctionRecoverContextFactory{h},
		// &auctionStatusContextFactory{h},
	}
	return h
}

func (h *ExampleHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/auction").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

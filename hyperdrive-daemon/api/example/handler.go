package example

import (
	"fmt"
	"net/http"

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

// ===============
// === Factory ===
// ===============
type exampleContext struct {
	handler *ExampleHandler
	// rp      *rocketpool.RocketPool

	// lotIndex  uint64
	// amountWei *big.Int
	// lot       *auction.AuctionLot
	// pSettings *protocol.ProtocolDaoSettings
}

type exampleContextFactory struct {
	handler *ExampleHandler
}

func (f *exampleContextFactory) Create(vars map[string]string) (*exampleContext, error) {
	c := &exampleContext{
		handler: f.handler,
	}
	return c, nil

	// inputErrs := []error{
	// 	server.ValidateArg("index", vars, input.ValidateUint, &c.lotIndex),
	// 	server.ValidateArg("amount", vars, input.ValidatePositiveWeiAmount, &c.amountWei),
	// }
	// return c, errors.Join(inputErrs...)
}

func (f *exampleContextFactory) RegisterRoute(router *mux.Router) {
	RegisterSingleStageRoute(
		router, "example-function-name",
	)
}

// ===============

func NewExampleHandler(serviceProvider *services.ServiceProvider) *ExampleHandler {
	h := &ExampleHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []IContextFactory{
		&exampleContextFactory{h},
	}
	return h
}

func (h *ExampleHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/example").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}

// Registers a new route with the router, which will invoke the provided factory to create and execute the context
// for the route when it's called; use this for typical general-purpose calls
func RegisterSingleStageRoute(
	router *mux.Router,
	functionName string,
) {
	router.HandleFunc(fmt.Sprintf("/%s", functionName), func(w http.ResponseWriter, r *http.Request) {
		// Run the context's processing routine
		fmt.Print("Running context's processing routine\n")
	})
}

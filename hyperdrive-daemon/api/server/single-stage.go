package server

// import (
// 	"fmt"
// 	"net/http"
// 	"net/url"

// 	"github.com/gorilla/mux"
// 	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/services"
// 	batch "github.com/rocket-pool/batch-query"
// )

// // Wrapper for callbacks used by call runners that follow a common single-stage pattern:
// // Create bindings, query the chain, and then do whatever else they want.
// // Structs implementing this will handle the caller-specific functionality.
// type ISingleStageCallContext[DataType any] interface {
// 	// Initialize the context with any bootstrapping, requirements checks, or bindings it needs to set up
// 	Initialize() error

// 	// Used to get any supplemental state required during initialization - anything in here will be fed into an rp.Query() multicall
// 	GetState(mc *batch.MultiCaller)

// 	// Prepare the response data in whatever way the context needs to do
// 	//PrepareData(data *DataType, opts *bind.TransactOpts) error
// 	PrepareData(data *DataType) error
// }

// // Interface for single-stage call context factories - these will be invoked during route handling to create the
// // unique context for the route
// type ISingleStageCallContextFactory[ContextType ISingleStageCallContext[DataType], DataType any] interface {
// 	// Create the context for the route
// 	Create(args url.Values) (ContextType, error)
// }

// // Registers a new route with the router, which will invoke the provided factory to create and execute the context
// // for the route when it's called; use this for typical general-purpose calls
// func RegisterSingleStageRoute[ContextType ISingleStageCallContext[DataType], DataType any](
// 	router *mux.Router,
// 	functionName string,
// 	factory ISingleStageCallContextFactory[ContextType, DataType],
// 	serviceProvider *services.ServiceProvider,
// ) {
// 	router.HandleFunc(fmt.Sprintf("/%s", functionName), func(w http.ResponseWriter, r *http.Request) {
// 		// Create the handler and deal with any input validation errors
// 		args := r.URL.Query()
// 		context, err := factory.Create(args)
// 		if err != nil {
// 			handleInputError(w, err)
// 			return
// 		}

// 		// Run the context's processing routine
// 		response, err := runSingleStageRoute[DataType](context, serviceProvider)
// 		handleResponse(w, response, err)
// 	})
// }

// // TODO: Talk to Joe about merging branches on Rocketpool go
// func Query() error {

// 	return nil
// }

// // Run a route registered with the common single-stage querying pattern
// func runSingleStageRoute[DataType any](ctx ISingleStageCallContext[DataType], serviceProvider *services.ServiceProvider) (*ApiResponse[DataType], error) {
// 	// Get the services
// 	rp := serviceProvider.GetRocketPool()
// 	fmt.Println(rp)
// 	// Initialize the context with any bootstrapping, requirements checks, or bindings it needs to set up
// 	err := ctx.Initialize()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get the context-specific contract state
// 	err = Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting contract state: %w", err)
// 	}

// 	// Get the transact opts if this node is ready for transaction
// 	// TODO: WALLET WITH v2 LATER
// 	//var opts *bind.TransactOpts

// 	// Create the response and data
// 	data := new(DataType)
// 	response := &ApiResponse[DataType]{
// 		Data: data,
// 	}

// 	// Prep the data with the context-specific behavior
// 	err = ctx.PrepareData(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Return
// 	return response, nil
// }

package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	batch "github.com/rocket-pool/batch-query"
)

// Wrapper for callbacks used by call runners that follow a common single-stage pattern:
// Create bindings, query the chain, and then do whatever else they want.
// Structs implementing this will handle the caller-specific functionality.
type ISingleStageCallContext[DataType any] interface {
	// Initialize the context with any bootstrapping, requirements checks, or bindings it needs to set up
	Initialize() error

	// Used to get any supplemental state required during initialization - anything in here will be fed into an hd.Query() multicall
	GetState(mc *batch.MultiCaller)

	// Prepare the response data in whatever way the context needs to do
	PrepareData(data *DataType, opts *bind.TransactOpts) error
}

// Interface for single-stage call context factories - these will be invoked during route handling to create the
// unique context for the route
type ISingleStageGetContextFactory[ContextType ISingleStageCallContext[DataType], DataType any] interface {
	// Create the context for the route
	Create(args url.Values) (ContextType, error)
}

// Interface for queryless call context factories that handle POST requests.
// These will be invoked during route handling to create the unique context for the route
type ISingleStagePostContextFactory[ContextType ISingleStageCallContext[DataType], BodyType any, DataType any] interface {
	// Create the context for the route
	Create(body BodyType) (ContextType, error)
}

// Registers a new route with the router, which will invoke the provided factory to create and execute the context
// for the route when it's called; use this for typical general-purpose calls
func RegisterSingleStageRoute[ContextType ISingleStageCallContext[DataType], DataType any](
	router *mux.Router,
	functionName string,
	factory ISingleStageGetContextFactory[ContextType, DataType],
	serviceProvider *common.ServiceProvider,
) {
	router.HandleFunc(fmt.Sprintf("/%s", functionName), func(w http.ResponseWriter, r *http.Request) {
		// Log
		args := r.URL.Query()
		log := serviceProvider.GetApiLogger()
		isDebug := serviceProvider.IsDebugMode()
		if isDebug {
			log.Printlnf("[%s] => %s", r.Method, r.URL.String())
		} else {
			log.Printlnf("[%s] => %s", r.Method, r.URL.Path)
		}

		// Check the method
		if r.Method != http.MethodGet {
			server.HandleInvalidMethod(log, w)
			return
		}

		// Create the handler and deal with any input validation errors
		context, err := factory.Create(args)
		if err != nil {
			server.HandleInputError(log, w, err)
			return
		}

		// Run the context's processing routine
		response, err := runSingleStageRoute[DataType](context, serviceProvider)
		server.HandleResponse(log, w, response, err, isDebug)
	})
}

// Registers a new route with the router, which will invoke the provided factory to create and execute the context
// for the route when it's called via POST; use this for typical general-purpose calls
func RegisterSingleStagePost[ContextType ISingleStageCallContext[DataType], BodyType any, DataType any](
	router *mux.Router,
	functionName string,
	factory ISingleStagePostContextFactory[ContextType, BodyType, DataType],
	serviceProvider *common.ServiceProvider,
) {
	router.HandleFunc(fmt.Sprintf("/%s", functionName), func(w http.ResponseWriter, r *http.Request) {
		// Log
		log := serviceProvider.GetApiLogger()
		isDebug := serviceProvider.IsDebugMode()
		log.Printlnf("[%s] => %s", r.Method, r.URL.Path)

		// Check the method
		if r.Method != http.MethodPost {
			server.HandleInvalidMethod(log, w)
			return
		}

		// Read the body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			server.HandleInputError(log, w, fmt.Errorf("error reading request body: %w", err))
			return
		}
		if isDebug {
			log.Println(string(bodyBytes))
		}

		// Deserialize the body
		var body BodyType
		err = json.Unmarshal(bodyBytes, &body)
		if err != nil {
			server.HandleInputError(log, w, fmt.Errorf("error deserializing request body: %w", err))
			return
		}

		// Create the handler and deal with any input validation errors
		context, err := factory.Create(body)
		if err != nil {
			server.HandleInputError(log, w, err)
			return
		}

		// Run the context's processing routine
		response, err := runSingleStageRoute[DataType](context, serviceProvider)
		server.HandleResponse(log, w, response, err, isDebug)
	})
}

// Run a route registered with the common single-stage querying pattern
func runSingleStageRoute[DataType any](ctx ISingleStageCallContext[DataType], serviceProvider *common.ServiceProvider) (*api.ApiResponse[DataType], error) {
	// Get the services
	w := serviceProvider.GetWallet()
	q := serviceProvider.GetQueryManager()

	// Initialize the context with any bootstrapping, requirements checks, or bindings it needs to set up
	err := ctx.Initialize()
	if err != nil {
		return nil, err
	}

	// Get the context-specific contract state
	err = q.Query(func(mc *batch.MultiCaller) error {
		ctx.GetState(mc)
		return nil
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting contract state: %w", err)
	}

	// Get the transact opts if this node is ready for transaction
	var opts *bind.TransactOpts
	walletStatus, err := w.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("error getting wallet status: %w", err)
	}
	if utils.IsWalletReady(walletStatus) {
		var err error
		opts, err = w.GetTransactor()
		if err != nil {
			return nil, fmt.Errorf("error getting node account transactor: %w", err)
		}
	}

	// Create the response and data
	data := new(DataType)
	response := &api.ApiResponse[DataType]{
		Data: data,
	}

	// Prep the data with the context-specific behavior
	err = ctx.PrepareData(data, opts)
	if err != nil {
		return nil, err
	}

	// Return
	return response, nil
}

package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

const (
	baseUrl         string          = "http://hyperdrive/%s"
	jsonContentType string          = "application/json"
	apiColor        color.Attribute = color.FgHiCyan
)

type RequesterContext struct {
	socketPath string
	client     *http.Client
	debugMode  bool
	log        *log.ColorLogger
}

type IRequester interface {
	GetName() string
	GetRoute() string
	GetContext() *RequesterContext
}

// Binder for the Hyperdrive daemon API server
type ApiRequester struct {
	Service *ServiceRequester
	Tx      *TxRequester
	Utils   *UtilsRequester
	Wallet  *WalletRequester

	context *RequesterContext
}

// Creates a new API requester instance
func NewApiRequester(socketPath string, debugMode bool) *ApiRequester {
	apiRequester := &ApiRequester{
		context: &RequesterContext{
			socketPath: socketPath,
			debugMode:  debugMode,
		},
	}

	apiRequester.context.client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}

	log := log.NewColorLogger(apiColor)
	apiRequester.context.log = &log

	apiRequester.Service = NewServiceRequester(apiRequester.context)
	apiRequester.Tx = NewTxRequester(apiRequester.context)
	apiRequester.Utils = NewUtilsRequester(apiRequester.context)
	apiRequester.Wallet = NewWalletRequester(apiRequester.context)
	return apiRequester
}

// Submit a GET request to the API server
func sendGetRequest[DataType any](r IRequester, method string, requestName string, args map[string]string) (*api.ApiResponse[DataType], error) {
	if args == nil {
		args = map[string]string{}
	}
	response, err := rawGetRequest[DataType](r.GetContext(), fmt.Sprintf("%s/%s", r.GetRoute(), method), args)
	if err != nil {
		return nil, fmt.Errorf("error during %s %s request: %w", r.GetName(), requestName, err)
	}
	return response, nil
}

// Submit a GET request to the API server
func rawGetRequest[DataType any](context *RequesterContext, path string, params map[string]string) (*api.ApiResponse[DataType], error) {
	// Make sure the socket exists
	_, err := os.Stat(context.socketPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("the socket at [%s] does not exist - please start the Hyperdrive daemon and try again", context.socketPath)
	}

	// Create the request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(baseUrl, path), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Encode the params into a query string
	values := url.Values{}
	for name, value := range params {
		values.Add(name, value)
	}
	req.URL.RawQuery = values.Encode()

	// Debug log
	if context.debugMode {
		context.log.Printlnf("[DEBUG] Query: GET %s", req.URL.String())
	}

	// Run the request
	resp, err := context.client.Do(req)
	return handleResponse[DataType](context, resp, path, err)
}

// Submit a POST request to the API server
func sendPostRequest[DataType any](r IRequester, method string, requestName string, body any) (*api.ApiResponse[DataType], error) {
	// Serialize the body
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error serializing request body for %s %s: %w", r.GetName(), requestName, err)
	}

	response, err := rawPostRequest[DataType](r.GetContext(), fmt.Sprintf("%s/%s", r.GetRoute(), method), string(bytes))
	if err != nil {
		return nil, fmt.Errorf("error during %s %s request: %w", r.GetName(), requestName, err)
	}
	return response, nil
}

// Submit a POST request to the API server
func rawPostRequest[DataType any](context *RequesterContext, path string, body string) (*api.ApiResponse[DataType], error) {
	// Make sure the socket exists
	_, err := os.Stat(context.socketPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("the socket at [%s] does not exist - please start the Hyperdrive daemon and try again", context.socketPath)
	}

	// Debug log
	if context.debugMode {
		context.log.Printlnf("[DEBUG] Query: POST %s", path)
		context.log.Printlnf("[DEBUG] Body: %s", body)
	}

	resp, err := context.client.Post(fmt.Sprintf(baseUrl, path), jsonContentType, strings.NewReader(body))
	return handleResponse[DataType](context, resp, path, err)
}

// Processes a response to a request
func handleResponse[DataType any](context *RequesterContext, resp *http.Response, path string, err error) (*api.ApiResponse[DataType], error) {
	if err != nil {
		return nil, fmt.Errorf("error requesting %s: %w", path, err)
	}

	// Read the body
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)

	// Check if the request failed
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return nil, fmt.Errorf("server responded to %s with code %s but reading the response body failed: %w", path, resp.Status, err)
		}
		msg := string(bytes)
		return nil, fmt.Errorf("server responded to %s with code %s: [%s]", path, resp.Status, msg)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading the response body for %s: %w", path, err)
	}

	// Debug log
	if context.debugMode {
		context.log.Printlnf("[DEBUG] Response: %s", string(bytes))
	}

	// Deserialize the response into the provided type
	var parsedResponse api.ApiResponse[DataType]
	err = json.Unmarshal(bytes, &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("error deserializing response to %s: %w; original body: [%s]", path, err, string(bytes))
	}

	return &parsedResponse, nil
}

// Types that can be batched into a comma-delmited string
type BatchInputType interface {
	uint64 | common.Address
}

// Converts an array of inputs into a comma-delimited string
func makeBatchArg[DataType BatchInputType](input []DataType) string {
	results := make([]string, len(input))

	// Figure out how to stringify the input
	switch typedInput := any(&input).(type) {
	case *[]uint64:
		for i, index := range *typedInput {
			results[i] = strconv.FormatUint(index, 10)
		}
	case *[]common.Address:
		for i, address := range *typedInput {
			results[i] = address.Hex()
		}
	}
	return strings.Join(results, ",")
}

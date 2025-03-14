package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/api"
)

// Submit a GET request to the API server
func SendGetRequest[DataType any](r IRequester, method string, requestName string, args map[string]string) (*api.ApiResponse[DataType], error) {
	if args == nil {
		args = map[string]string{}
	}
	response, err := RawGetRequest[DataType](r.GetContext(), fmt.Sprintf("%s/%s", r.GetRoute(), method), args)
	if err != nil {
		return nil, fmt.Errorf("error during %s %s request: %w", r.GetName(), requestName, err)
	}
	return response, nil
}

// Submit a GET request to the API server
func RawGetRequest[DataType any](context IRequesterContext, path string, params map[string]string) (*api.ApiResponse[DataType], error) {
	// Create the request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", context.GetAddressBase(), path), nil)
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
	context.GetLogger().Debug(
		"API Request",
		"method", http.MethodGet,
		"query", req.URL.String(),
	)

	// Run the request
	resp, err := context.SendRequest(req)
	return HandleResponse[DataType](context, resp, path, err)
}

// Submit a POST request to the API server
func SendPostRequest[DataType any](r IRequester, method string, requestName string, body any) (*api.ApiResponse[DataType], error) {
	// Serialize the body
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error serializing request body for %s %s: %w", r.GetName(), requestName, err)
	}

	response, err := RawPostRequest[DataType](r.GetContext(), fmt.Sprintf("%s/%s", r.GetRoute(), method), string(bytes))
	if err != nil {
		return nil, fmt.Errorf("error during %s %s request: %w", r.GetName(), requestName, err)
	}
	return response, nil
}

// Submit a POST request to the API server
func RawPostRequest[DataType any](context IRequesterContext, path string, body string) (*api.ApiResponse[DataType], error) {
	// Create the request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", context.GetAddressBase(), path), strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", jsonContentType)

	// Debug log
	context.GetLogger().Debug(
		"API Request",
		"method", http.MethodPost,
		"path", path,
		"body", body,
	)

	// Run the request
	resp, err := context.SendRequest(req)
	return HandleResponse[DataType](context, resp, path, err)
}

// Processes a response to a request
func HandleResponse[DataType any](context IRequesterContext, resp *http.Response, path string, err error) (*api.ApiResponse[DataType], error) {
	if err != nil {
		return nil, fmt.Errorf("error requesting %s: %w", path, err)
	}
	logger := context.GetLogger()

	// Read the body
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response body for %s: %w", path, err)
	}

	// Handle 404s specially since they won't have a JSON body
	if resp.StatusCode == http.StatusNotFound {
		logger.Debug(
			"API Response (raw)",
			"code", resp.Status,
			"body", string(bytes),
		)
		return nil, fmt.Errorf("route '%s' not found", path)
	}

	// Deserialize the response into the provided type
	var parsedResponse api.ApiResponse[DataType]
	err = json.Unmarshal(bytes, &parsedResponse)
	if err != nil {
		logger.Debug(
			"API Response (raw)",
			"code", resp.Status,
			"body", string(bytes),
		)
		return nil, fmt.Errorf("error deserializing response to %s: %w", path, err)
	}

	// Check if the request failed
	if resp.StatusCode != http.StatusOK {
		logger.Debug(
			"API Response",
			"path", path,
			"code", resp.Status,
			"error", parsedResponse.Error,
		)
		return nil, errors.New(parsedResponse.Error)
	}

	// Debug log
	logger.Debug(
		"API Response",
		"body", string(bytes),
	)

	return &parsedResponse, nil
}

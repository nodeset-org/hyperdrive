package api_v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/nodeset-org/hyperdrive-daemon/nodeset"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive-daemon/shared/keys"
	"github.com/rocket-pool/node-manager-core/log"
)

const (
	// API version to use
	apiVersion string = "v2"

	// Header to include when sending messages that have been logged in
	authHeader string = "Authorization"

	// Format for the authorization header
	authHeaderFormat string = "Bearer %s"
)

var (
	// A login session hasn't been established yet
	ErrNoSession error = errors.New("not logged in yet")
)

// Client for interacting with the NodeSet server
type NodeSetClient struct {
	baseUrl     string
	session     *nodeset.Session
	networkName string
	client      *http.Client
}

// Creates a new NodeSet client
func NewNodeSetClient(resources *config.HyperdriveResources, timeout time.Duration) *NodeSetClient {
	return &NodeSetClient{
		baseUrl:     fmt.Sprintf("%s/%s", resources.NodeSetApiUrl, apiVersion),
		networkName: resources.EthNetworkName,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Set the session for the client after logging in
func (c *NodeSetClient) SetSession(session *nodeset.Session) {
	c.session = session
}

// Send a request to the server and read the response
func SubmitRequest[DataType any](c *NodeSetClient, ctx context.Context, requireAuth bool, method string, body io.Reader, queryParams map[string]string, subroutes ...string) (int, NodeSetResponse[DataType], error) {
	var defaultVal NodeSetResponse[DataType]

	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Make the request
	path, err := url.JoinPath(c.baseUrl, subroutes...)
	if err != nil {
		return 0, defaultVal, fmt.Errorf("error joining path [%v]: %w", subroutes, err)
	}
	request, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return 0, defaultVal, fmt.Errorf("error generating request to [%s]: %w", path, err)
	}
	query := request.URL.Query()
	for name, value := range queryParams {
		query.Add(name, value)
	}
	request.URL.RawQuery = query.Encode()

	// Set the headers
	if requireAuth {
		if c.session == nil {
			return 0, defaultVal, ErrNoSession
		}
		request.Header.Set(authHeader, fmt.Sprintf(authHeaderFormat, c.session.Token))
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Upload it to the server
	logger.Debug("Sending NodeSet server request", slog.String(log.QueryKey, request.URL.String()))

	resp, err := c.client.Do(request)
	if err != nil {
		return 0, defaultVal, fmt.Errorf("error submitting request to nodeset server: %w", err)
	}

	// Read the body
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, defaultVal, fmt.Errorf("nodeset server responded to request with code %s but reading the response body failed: %w", resp.Status, err)
	}

	// Unmarshal the response
	var response NodeSetResponse[DataType]
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return 0, defaultVal, fmt.Errorf("nodeset server responded to request with code %s and unmarshalling the response failed: [%w]... original body: [%s]", resp.Status, err, string(bytes))
	}

	// Debug log
	logger.Debug("NodeSet response:", slog.String(log.CodeKey, resp.Status), slog.String(keys.MessageKey, response.Message))
	return resp.StatusCode, response, nil
}

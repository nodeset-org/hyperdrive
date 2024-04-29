package swcommon

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/goccy/go-json"

	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/utils"
)

const (
	// The message to sign with the node wallet when uploading deposit data
	nodesetAuthMessage string = "nodesetdev"

	// Header used for the wallet signature during a deposit data upload
	authHeader string = "Authorization"

	// API paths
	devPath         string = "dev"
	depositDataPath string = "deposit-data"
	metaPath        string = "meta"
	validatorsPath  string = "validators"
)

// =================
// === Requests  ===
// =================
type ExitMessageDetails struct {
	Epoch          string `json:"epoch"`
	ValidatorIndex string `json:"validator_index"`
}

type ExitMessage struct {
	Message   ExitMessageDetails `json:"message"`
	Signature string             `json:"signature"`
}

type ExitData struct {
	Pubkey      string      `json:"pubkey"`
	ExitMessage ExitMessage `json:"exit_message"`
}

// =================
// === Responses ===
// =================

// api/deposit-data/meta
type DepositDataMetaResponse struct {
	Version int `json:"version"`
}

// api/deposit-data
type DepositDataResponse struct {
	Version int                         `json:"version"`
	Data    []types.ExtendedDepositData `json:"data"`
}

// api/validators
type ValidatorStatus struct {
	Pubkey              beacon.ValidatorPubkey `json:"pubkey"`
	Status              string                 `json:"status"`
	ExitMessageUploaded bool                   `json:"exitMessage"`
}

// api/dev/validators
type ValidatorsResponse struct {
	Data []ValidatorStatus `json:"data"`
}

// ==============
// === Client ===
// ==============

// Client for interacting with the Nodeset server
type NodesetClient struct {
	sp            *StakewiseServiceProvider
	res           *swconfig.StakewiseResources
	authSignature []byte
}

// Creates a new Nodeset client
func NewNodesetClient(sp *StakewiseServiceProvider) *NodesetClient {
	return &NodesetClient{
		sp:  sp,
		res: sp.GetResources(),
	}
}

// Checks if the NodeSet authorization signature has been set, and if not, creates it by getting a signed message from the node wallet
func (c *NodesetClient) EnsureAuthSignatureExists() error {
	if c.authSignature != nil {
		return nil
	}

	// Sign the auth message
	hd := c.sp.GetHyperdriveClient()
	signResponse, err := hd.Wallet.SignMessage([]byte(nodesetAuthMessage))
	if err != nil {
		return fmt.Errorf("error signing authorization message: %w", err)
	}
	c.authSignature = signResponse.Data.SignedMessage
	return nil
}

// Uploads deposit data to Nodeset
func (c *NodesetClient) UploadDepositData(ctx context.Context, depositData []byte) ([]byte, error) {
	response, err := c.submitRequest(ctx, http.MethodPost, bytes.NewBuffer(depositData), nil, depositDataPath)
	if err != nil {
		return nil, fmt.Errorf("error uploading deposit data: %w", err)
	}
	return response, nil
}

// Submit signed exit data to Nodeset
func (c *NodesetClient) UploadSignedExitData(ctx context.Context, exitData []ExitData) ([]byte, error) {
	// Serialize the exit data into JSON
	jsonData, err := json.Marshal(exitData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling exit data to JSON: %w", err)
	}
	params := map[string]string{
		"network": c.res.EthNetworkName,
	}
	// Submit the PATCH request with the serialized JSON
	response, err := c.submitRequest(ctx, http.MethodPatch, bytes.NewBuffer(jsonData), params, devPath, validatorsPath)
	if err != nil {
		return nil, fmt.Errorf("error submitting exit data: %w", err)
	}

	return response, nil
}

// Get the current version of the aggregated deposit data on the server
func (c *NodesetClient) GetServerDepositDataVersion(ctx context.Context) (int, error) {
	vault := utils.RemovePrefix(strings.ToLower(c.res.Vault.Hex()))
	params := map[string]string{
		"vault":   vault,
		"network": c.res.EthNetworkName,
	}
	response, err := c.submitRequest(ctx, http.MethodGet, nil, params, depositDataPath, metaPath)
	if err != nil {
		return 0, fmt.Errorf("error getting deposit data version: %w", err)
	}

	var body DepositDataMetaResponse
	err = json.Unmarshal(response, &body)
	if err != nil {
		return 0, fmt.Errorf("error deserializing deposit data version response: %w", err)
	}
	return body.Version, nil
}

// Get the aggregated deposit data from the server
func (c *NodesetClient) GetServerDepositData(ctx context.Context) (int, []types.ExtendedDepositData, error) {
	vault := utils.RemovePrefix(strings.ToLower(c.res.Vault.Hex()))
	params := map[string]string{
		"vault":   vault,
		"network": c.res.EthNetworkName,
	}
	response, err := c.submitRequest(ctx, http.MethodGet, nil, params, depositDataPath)
	if err != nil {
		return 0, nil, fmt.Errorf("error getting deposit data: %w", err)
	}

	var body DepositDataResponse
	err = json.Unmarshal(response, &body)
	if err != nil {
		return 0, nil, fmt.Errorf("error deserializing deposit data response: %w", err)
	}
	return body.Version, body.Data, nil
}

// Get a list of all of the pubkeys that have already been registered with NodeSet for this node
func (c *NodesetClient) GetRegisteredValidators(ctx context.Context) ([]ValidatorStatus, error) {
	queryParams := map[string]string{
		"network": c.res.EthNetworkName,
	}
	response, err := c.submitRequest(ctx, http.MethodGet, nil, queryParams, devPath, validatorsPath)
	if err != nil {
		return nil, fmt.Errorf("error getting registered validators: %w", err)
	}

	var body ValidatorsResponse
	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, fmt.Errorf("error deserializing registered validators response: %w", err)
	}
	return body.Data, nil
}

// Send a request to the server and read the response
func (c *NodesetClient) submitRequest(ctx context.Context, method string, body io.Reader, queryParams map[string]string, subroutes ...string) ([]byte, error) {
	// Get the logger
	logger, exists := log.FromContext(ctx)
	if !exists {
		panic("context didn't have a logger!")
	}

	// Make the request
	path, err := url.JoinPath(c.res.NodesetApiUrl, subroutes...)
	if err != nil {
		return nil, fmt.Errorf("error joining path [%v]: %w", subroutes, err)
	}
	request, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, fmt.Errorf("error generating request to [%s]: %w", path, err)
	}
	query := request.URL.Query()
	for name, value := range queryParams {
		query.Add(name, value)
	}
	request.URL.RawQuery = query.Encode()

	// Set the headers
	err = c.EnsureAuthSignatureExists()
	if err != nil {
		return nil, fmt.Errorf("initializing authorization signature failed: %w", err)
	}
	request.Header.Set(authHeader, utils.EncodeHexWithPrefix(c.authSignature))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Upload it to the server
	logger.Debug("Sending NodeSet server request", slog.String(log.QueryKey, request.URL.String()))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error submitting request to nodeset server: %w", err)
	}

	// Read the body
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)

	// Check if the request failed
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return nil, fmt.Errorf("nodeset server responded to request with code %s but reading the response body failed: %w", resp.Status, err)
		}
		msg := string(bytes)
		return nil, fmt.Errorf("nodeset server responded to request with code %s: [%s]", resp.Status, msg)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading the response body for nodeset request: %w", err)
	}

	// Debug log
	logger.Debug("NodeSet response:", slog.String(log.CodeKey, resp.Status), slog.String(log.BodyKey, string(bytes)))
	return bytes, nil
}

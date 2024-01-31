package common

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/eth-utils/common"
)

const (
	// The message to sign with the node wallet when uploading deposit data
	NodesetAuthMessage string = "nodesetdev"

	// Header used for the wallet signature during a deposit data upload
	nodesetAuthHeader string = "Authorization"
)

// =================
// === Responses ===
// =================

// api/deposit-data/meta
type DepositDataMetaResponse struct {
	Version uint64 `json:"version"`
}

// api/validators
type ValidatorsResponse struct {
	Data []beacon.ValidatorPubkey `json:"data"`
}

// ==============
// === Client ===
// ==============

// Client for interacting with the Nodeset server
type NodesetClient struct {
	sp    *StakewiseServiceProvider
	debug bool
}

// Creates a new Nodeset client
func NewNodesetClient(sp *StakewiseServiceProvider) *NodesetClient {
	cfg := sp.GetConfig()
	return &NodesetClient{
		sp:    sp,
		debug: cfg.DebugMode.Value,
	}
}

// Uploads deposit data to Nodeset
func (c *NodesetClient) UploadDepositData(signature []byte, depositData []byte) ([]byte, error) {
	res := c.sp.GetResources()

	// Create a new POST request
	request, err := http.NewRequest(http.MethodPost, res.NodesetDepositUrl, bytes.NewBuffer(depositData))
	if err != nil {
		return nil, fmt.Errorf("error generating request: %w", err)
	}
	request.Header.Set(nodesetAuthHeader, common.EncodeHexWithPrefix(signature))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	/*
		This is way too big - maybe turn it on later?
		if debug {
			buffer := bytes.Buffer{}
			request.Write(&buffer)
			fmt.Printf("[%s] => [%s]\n", res.NodesetDepositUrl, buffer.String())
		}
	*/
	beacon.ValidatorPubkey{}.Hex()
	// Upload it to the server
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
			return nil, fmt.Errorf("nodeset server responded to upload request with code %s but reading the response body failed: %w", resp.Status, err)
		}
		msg := string(bytes)
		return nil, fmt.Errorf("nodeset server responded to upload request with code %s: [%s]", resp.Status, msg)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading the response body for nodeset upload request: %w", err)
	}

	// Debug log
	if c.debug {
		fmt.Printf("NODESET UPLOAD RESPONSE: %s\n", string(bytes))
	}
	return bytes, nil
}

// Downloads complete merged deposit data from the Nodeset server
func (c *NodesetClient) DownloadDepositData() error {
	return fmt.Errorf("NYI")
}

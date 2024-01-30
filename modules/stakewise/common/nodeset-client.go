package common

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

const (
	// The message to sign with the node wallet when uploading deposit data
	NodesetAuthMessage string = "nodesetdev"

	// Header used for the wallet signature during a deposit data upload
	nodesetAuthHeader string = "Authorization"
)

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
func (c *NodesetClient) UploadDepositData(signature []byte, depositData []byte) error {
	res := c.sp.GetResources()

	// Create a new POST request
	request, err := http.NewRequest(http.MethodPost, res.NodesetDepositUrl, bytes.NewBuffer(depositData))
	if err != nil {
		return fmt.Errorf("error generating request: %w", err)
	}
	request.Header.Set(nodesetAuthHeader, hex.EncodeToString(signature))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	/*
		This is way too big - maybe turn it on later?
		if debug {
			buffer := bytes.Buffer{}
			request.Write(&buffer)
			fmt.Printf("[%s] => [%s]\n", res.NodesetDepositUrl, buffer.String())
		}
	*/

	// Upload it to the server
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error submitting request to nodeset server: %w", err)
	}

	// Read the body
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)

	// Check if the request failed
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return fmt.Errorf("nodeset server responded to upload request with code %s but reading the response body failed: %w", resp.Status, err)
		}
		msg := string(bytes)
		return fmt.Errorf("nodeset server responded to upload request with code %s: [%s]", resp.Status, msg)
	}
	if err != nil {
		return fmt.Errorf("error reading the response body for nodeset upload request: %w", err)
	}

	// Debug log
	if c.debug {
		fmt.Printf("NODESET UPLOAD RESPONSE: %s\n", string(bytes))
	}
	return nil
}

// Downloads complete merged deposit data from the Nodeset server
func (c *NodesetClient) DownloadDepositData() error {
	return fmt.Errorf("NYI")
}

package swnodeset

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

const (
	nodesetAuthMessage string = "nodesetdev"
	nodesetAuthHeader  string = "Authorization"
)

// ===============
// === Factory ===
// ===============

type nodesetUploadDepositDataContextFactory struct {
	handler *NodesetHandler
}

func (f *nodesetUploadDepositDataContextFactory) Create(args url.Values) (*nodesetUploadDepositDataContext, error) {
	c := &nodesetUploadDepositDataContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *nodesetUploadDepositDataContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*nodesetUploadDepositDataContext, api.SuccessData](
		router, "upload-deposit-data", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type nodesetUploadDepositDataContext struct {
	handler *NodesetHandler
}

func (c *nodesetUploadDepositDataContext) PrepareData(data *api.SuccessData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	ddMgr := sp.GetDepositDataManager()
	res := sp.GetResources()
	hd := sp.GetClient()

	// Read the deposit data
	depositData, err := ddMgr.GetDepositData()
	if err != nil {
		return err
	}

	// Sign a message
	signResponse, err := hd.Wallet.SignMessage([]byte(nodesetAuthMessage))
	if err != nil {
		return fmt.Errorf("error signing authorization message: %w", err)
	}
	signature := signResponse.Data.SignedMessage

	// Create a new POST request
	request, err := http.NewRequest(http.MethodPost, res.NodesetDepositUrl, bytes.NewBuffer(depositData))
	if err != nil {
		return fmt.Errorf("error generating request: %w", err)
	}
	request.Header.Set(nodesetAuthHeader, hex.EncodeToString(signature))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

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
	if sp.GetConfig().DebugMode.Value {
		fmt.Printf("NODESET UPLOAD RESPONSE: %s\n", string(bytes))
	}
	return nil
}

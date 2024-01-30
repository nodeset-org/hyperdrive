package swnodeset

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
	"github.com/nodeset-org/hyperdrive/modules/stakewise/common"
	swapi "github.com/nodeset-org/hyperdrive/shared/types/api/modules/stakewise"
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
	server.RegisterQuerylessGet[*nodesetUploadDepositDataContext, swapi.NodesetUploadDepositDataData](
		router, "upload-deposit-data", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type nodesetUploadDepositDataContext struct {
	handler *NodesetHandler
}

func (c *nodesetUploadDepositDataContext) PrepareData(data *swapi.NodesetUploadDepositDataData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	ddMgr := sp.GetDepositDataManager()
	hd := sp.GetHyperdriveClient()
	nc := sp.GetNodesetClient()

	// Read the deposit data
	depositData, err := ddMgr.GetDepositData()
	if err != nil {
		return err
	}

	// Sign a message
	signResponse, err := hd.Wallet.SignMessage([]byte(common.NodesetAuthMessage))
	if err != nil {
		return fmt.Errorf("error signing authorization message: %w", err)
	}
	signature := signResponse.Data.SignedMessage

	// Submit the upload
	response, err := nc.UploadDepositData(signature, depositData)
	if err != nil {
		return err
	}
	data.ServerResponse = response
	return nil
}

package swnodeset

import (
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
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
	nc := sp.GetNodesetClient()

	// Read the deposit data
	depositData, err := ddMgr.GetDepositData()
	if err != nil {
		return err
	}

	// Submit the upload
	response, err := nc.UploadDepositData(depositData)
	if err != nil {
		return err
	}
	data.ServerResponse = response
	return nil
}

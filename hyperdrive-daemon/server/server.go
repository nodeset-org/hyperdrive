package server

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/service"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/tx"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/wallet"
	"github.com/nodeset-org/hyperdrive/modules/common/server"
	"github.com/nodeset-org/hyperdrive/shared/config"
)

type HyperdriveServer struct {
	*server.ApiManager
}

func NewHyperdriveServer(sp *common.ServiceProvider, socketPath string) (*HyperdriveServer, error) {
	handlers := []server.IHandler{
		service.NewServiceHandler(sp),
		tx.NewTxHandler(sp),
		utils.NewUtilsHandler(sp),
		wallet.NewWalletHandler(sp),
	}

	mgr, err := server.NewApiServer(socketPath, handlers, config.HyperdriveDaemonRoute)
	if err != nil {
		return nil, err
	}

	return &HyperdriveServer{
		ApiManager: mgr,
	}, nil
}

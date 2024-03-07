package server

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/service"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/tx"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/wallet"
	"github.com/nodeset-org/hyperdrive/shared/config"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
)

type HyperdriveServer struct {
	*nmc_server.ApiServer
}

func NewHyperdriveServer(sp *common.ServiceProvider, socketPath string) (*HyperdriveServer, error) {
	handlers := []nmc_server.IHandler{
		service.NewServiceHandler(sp),
		tx.NewTxHandler(sp),
		utils.NewUtilsHandler(sp),
		wallet.NewWalletHandler(sp),
	}

	server, err := nmc_server.NewApiServer(socketPath, handlers, config.HyperdriveDaemonRoute)
	if err != nil {
		return nil, err
	}

	return &HyperdriveServer{
		ApiServer: server,
	}, nil
}

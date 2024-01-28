package api

import (
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/daemons/common/server"
	"github.com/nodeset-org/hyperdrive/daemons/common/services"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/service"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/tx"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/wallet"
	"github.com/nodeset-org/hyperdrive/shared/config"
)

type HyperdriveServer struct {
	*server.ApiManager
}

func NewHyperdriveServer(sp *services.ServiceProvider) *HyperdriveServer {
	socketPath := filepath.Join(sp.GetUserDir(), config.SocketFilename)
	handlers := []server.IHandler{
		service.NewServiceHandler(sp),
		tx.NewTxHandler(sp),
		utils.NewUtilsHandler(sp),
		wallet.NewWalletHandler(sp),
	}
	mgr := server.NewApiManager(sp, socketPath, handlers, config.HyperdriveDaemonRoute)

	return &HyperdriveServer{
		ApiManager: mgr,
	}
}

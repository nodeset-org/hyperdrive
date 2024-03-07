package server

import (
	"path/filepath"

	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	swnodeset "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/nodeset"
	swstatus "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/status"
	swvalidator "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/validator"
	swwallet "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/wallet"
	"github.com/rocket-pool/node-manager-core/api/server"
)

type StakewiseServer struct {
	*server.ApiServer
	socketPath string
}

func NewStakewiseServer(sp *swcommon.StakewiseServiceProvider) (*StakewiseServer, error) {
	socketPath := filepath.Join(sp.GetUserDir(), swconfig.SocketFilename)
	handlers := []server.IHandler{
		swnodeset.NewNodesetHandler(sp),
		swvalidator.NewValidatorHandler(sp),
		swwallet.NewWalletHandler(sp),
		swstatus.NewStatusHandler(sp),
	}
	server, err := server.NewApiServer(socketPath, handlers, swconfig.ModuleName)
	if err != nil {
		return nil, err
	}

	return &StakewiseServer{
		ApiServer:  server,
		socketPath: socketPath,
	}, nil
}

func (s *StakewiseServer) GetSocketPath() string {
	return s.socketPath
}

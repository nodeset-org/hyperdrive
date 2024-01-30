package server

import (
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	swnodeset "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/api/nodeset"
	swwallet "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/api/wallet"
)

type StakewiseServer struct {
	*server.ApiManager
	socketPath string
}

func NewStakewiseServer(sp *common.StakewiseServiceProvider) (*StakewiseServer, error) {
	socketPath := filepath.Join(sp.GetUserDir(), swconfig.SocketFilename)
	handlers := []server.IHandler{
		swnodeset.NewNodesetHandler(sp),
		swwallet.NewWalletHandler(sp),
	}
	mgr, err := server.NewApiServer(socketPath, handlers, swconfig.DaemonRoute)
	if err != nil {
		return nil, err
	}

	return &StakewiseServer{
		ApiManager: mgr,
		socketPath: socketPath,
	}, nil
}

func (s *StakewiseServer) GetSocketPath() string {
	return s.socketPath
}

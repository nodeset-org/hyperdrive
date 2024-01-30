package server

import (
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/modules/common/server"
	"github.com/nodeset-org/hyperdrive/modules/stakewise/common"
	swwallet "github.com/nodeset-org/hyperdrive/modules/stakewise/server/api/wallet"
	swconfig "github.com/nodeset-org/hyperdrive/shared/config/modules/stakewise"
)

type StakewiseServer struct {
	*server.ApiManager
	socketPath string
}

func NewStakewiseServer(sp *common.StakewiseServiceProvider) (*StakewiseServer, error) {
	socketPath := filepath.Join(sp.GetUserDir(), swconfig.SocketFilename)
	handlers := []server.IHandler{
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

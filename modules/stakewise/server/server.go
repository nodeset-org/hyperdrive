package server

import (
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/modules/common/server"
	"github.com/nodeset-org/hyperdrive/modules/common/services"
	swwallet "github.com/nodeset-org/hyperdrive/modules/stakewise/server/api/wallet"
	swconfig "github.com/nodeset-org/hyperdrive/shared/config/modules/stakewise"
)

type StakewiseServer struct {
	*server.ApiManager
	socketPath string
}

func NewStakewiseServer(sp *services.ServiceProvider) (*StakewiseServer, error) {
	socketPath := filepath.Join(sp.GetUserDir(), swconfig.StakewiseSocketFilename)
	handlers := []server.IHandler{
		swwallet.NewWalletHandler(sp),
	}
	mgr, err := server.NewApiServer(socketPath, handlers, swconfig.StakewiseDaemonRoute)
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

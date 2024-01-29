package api

import (
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/modules/common/server"
	"github.com/nodeset-org/hyperdrive/modules/common/services"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/modules/stakewise"
)

type StakewiseServer struct {
	*server.ApiManager
}

func NewStakewiseServer(sp *services.ServiceProvider) (*StakewiseServer, error) {
	socketPath := filepath.Join(sp.GetModuleDir(), config.HyperdriveSocketFilename)
	handlers := []server.IHandler{}
	mgr, err := server.NewApiServer(socketPath, handlers, stakewise.StakewiseDaemonRoute)
	if err != nil {
		return nil, err
	}

	return &StakewiseServer{
		ApiManager: mgr,
	}, nil
}

package api

import (
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/daemons/common/server"
	"github.com/nodeset-org/hyperdrive/daemons/common/services"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/stakewise"
)

type StakewiseServer struct {
	*server.ApiManager
}

func NewStakewiseServer(sp *services.ServiceProvider) *StakewiseServer {
	socketPath := filepath.Join(sp.GetUserDir(), config.ModulesDir, stakewise.DaemonRoute)
	handlers := []server.IHandler{}
	mgr := server.NewApiManager(sp, socketPath, handlers, stakewise.DaemonRoute)

	return &StakewiseServer{
		ApiManager: mgr,
	}
}

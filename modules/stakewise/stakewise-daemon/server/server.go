package server

import (
	"fmt"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	swnodeset "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/nodeset"
	swstatus "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/status"
	swvalidator "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/validator"
	swwallet "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server/wallet"
)

type StakewiseServer struct {
	*server.ApiManager
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
	fmt.Printf("!!! NewStakewiseServer handlers: %+v\n", handlers)
	mgr, err := server.NewApiServer(socketPath, handlers, swconfig.ModuleName)
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

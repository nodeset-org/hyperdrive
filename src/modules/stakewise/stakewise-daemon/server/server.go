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

const (
	CliOrigin string = "cli"
	WebOrigin string = "net"
)

type StakewiseServer struct {
	*server.ApiServer
	socketPath string
}

func NewStakewiseServer(origin string, sp *swcommon.StakewiseServiceProvider) (*StakewiseServer, error) {
	apiLogger := sp.GetApiLogger()
	subLogger := apiLogger.CreateSubLogger(origin)
	ctx := subLogger.CreateContextWithLogger(sp.GetBaseContext())

	socketPath := filepath.Join(sp.GetUserDir(), swconfig.CliSocketFilename)
	handlers := []server.IHandler{
		swnodeset.NewNodesetHandler(subLogger, ctx, sp),
		swvalidator.NewValidatorHandler(subLogger, ctx, sp),
		swwallet.NewWalletHandler(subLogger, ctx, sp),
		swstatus.NewStatusHandler(subLogger, ctx, sp),
	}
	server, err := server.NewApiServer(subLogger.Logger, socketPath, handlers, swconfig.DaemonBaseRoute, swconfig.ApiVersion)
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

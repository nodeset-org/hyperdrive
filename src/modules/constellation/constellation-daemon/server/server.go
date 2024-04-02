package server

import (
	"path/filepath"

	constcommon "github.com/nodeset-org/hyperdrive/modules/constellation/constellation-daemon/common"
	constconfig "github.com/nodeset-org/hyperdrive/modules/constellation/shared/config"
	"github.com/rocket-pool/node-manager-core/api/server"
)

const (
	CliOrigin string = "cli"
	WebOrigin string = "net"
)

type ConstellationServer struct {
	*server.ApiServer
	socketPath string
}

func NewConstellationServer(origin string, sp *constcommon.ConstellationServiceProvider) (*ConstellationServer, error) {
	apiLogger := sp.GetApiLogger()
	subLogger := apiLogger.CreateSubLogger(origin)
	// ctx := subLogger.CreateContextWithLogger(sp.GetBaseContext())

	socketPath := filepath.Join(sp.GetUserDir(), constconfig.CliSocketFilename)
	handlers := []server.IHandler{
		// constnodeset.NewNodesetHandler(subLogger, ctx, sp),
		// constvalidator.NewValidatorHandler(subLogger, ctx, sp),
		// constwallet.NewWalletHandler(subLogger, ctx, sp),
		// conststatus.NewStatusHandler(subLogger, ctx, sp),
	}
	server, err := server.NewApiServer(subLogger.Logger, socketPath, handlers, constconfig.DaemonBaseRoute, constconfig.ApiVersion)
	if err != nil {
		return nil, err
	}

	return &ConstellationServer{
		ApiServer:  server,
		socketPath: socketPath,
	}, nil
}

func (s *ConstellationServer) GetSocketPath() string {
	return s.socketPath
}

package client

import (
	"fmt"
	"log/slog"
	"path/filepath"

	docker "github.com/docker/docker/client"
	"github.com/fatih/color"
	csclient "github.com/nodeset-org/hyperdrive-constellation/client"
	"github.com/nodeset-org/hyperdrive-daemon/client"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
	swclient "github.com/nodeset-org/hyperdrive-stakewise/client"
	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/context"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/urfave/cli/v2"
)

// Config
const (
	InstallerName string = "install.sh"
	InstallerURL  string = "https://github.com/nodeset-org/hyperdrive/releases/download/%s/" + InstallerName

	SettingsFile       string = "user-settings.yml"
	BackupSettingsFile string = "user-settings-backup.yml"

	terminalLogColor color.Attribute = color.FgHiYellow
)

// Hyperdrive client
type HyperdriveClient struct {
	Api      *client.ApiClient
	Context  *context.HyperdriveContext
	Logger   *slog.Logger
	docker   *docker.Client
	cfg      *GlobalConfig
	isNewCfg bool
}

// Stakewise client
type StakewiseClient struct {
	Api     *swclient.ApiClient
	Context *context.HyperdriveContext
	Logger  *slog.Logger
}

// Constellation client
type ConstellationClient struct {
	Api     *csclient.ApiClient
	Context *context.HyperdriveContext
	Logger  *slog.Logger
}

// Create new Hyperdrive client from CLI context without checking for sync status
// Only use this function from commands that may work if the Daemon service doesn't exist
// Most users should call NewHyperdriveClientFromCtx(c).WithStatus() or NewHyperdriveClientFromCtx(c).WithReady()
func NewHyperdriveClientFromCtx(c *cli.Context) *HyperdriveClient {
	hdCtx := context.GetHyperdriveContext(c)
	socketPath := filepath.Join(hdCtx.ConfigPath, config.HyperdriveCliSocketFilename)

	// Make the client
	logger := log.NewTerminalLogger(hdCtx.DebugEnabled, terminalLogColor).With(slog.String(log.OriginKey, config.HyperdriveDaemonRoute))
	client := &HyperdriveClient{
		Api:     client.NewApiClient(config.HyperdriveApiClientRoute, socketPath, logger),
		Context: hdCtx,
		Logger:  logger,
	}
	return client
}

// Create new Stakewise client from CLI context without checking for sync status
// Only use this function from commands that may work if the Daemon service doesn't exist
func NewStakewiseClientFromCtx(c *cli.Context) *StakewiseClient {
	hdCtx := context.GetHyperdriveContext(c)
	socketPath := filepath.Join(hdCtx.ConfigPath, swconfig.CliSocketFilename)

	// Make the client
	logger := log.NewTerminalLogger(hdCtx.DebugEnabled, terminalLogColor).With(slog.String(log.OriginKey, swconfig.ModuleName))
	client := &StakewiseClient{
		Api:     swclient.NewApiClient(swconfig.ApiClientRoute, socketPath, logger),
		Context: hdCtx,
		Logger:  logger,
	}
	return client
}

// Get the Docker client
func (c *HyperdriveClient) GetDocker() (*docker.Client, error) {
	if c.docker == nil {
		var err error
		c.docker, err = docker.NewClientWithOpts(docker.WithAPIVersionNegotiation())
		if err != nil {
			return nil, fmt.Errorf("error creating Docker client: %w", err)
		}
	}

	return c.docker, nil
}

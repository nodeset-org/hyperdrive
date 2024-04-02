package client

import (
	"fmt"
	"log/slog"
	"path/filepath"

	docker "github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/nodeset-org/hyperdrive/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/context"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	swclient "github.com/nodeset-org/hyperdrive/modules/stakewise/client"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config"
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

// Check the status of a newly created client and return it
// Only use this function from commands that may work without the clients being synced-
// most users should use WithReady instead
func (c *HyperdriveClient) WithStatus() (*HyperdriveClient, bool, error) {
	ready, err := c.checkClientStatus()
	if err != nil {
		return nil, false, err
	}

	return c, ready, nil
}

// Check the status of a newly created client and ensure the eth clients are synced and ready
func (c *HyperdriveClient) WithReady() (*HyperdriveClient, error) {
	_, ready, err := c.WithStatus()
	if err != nil {
		return nil, err
	}

	if !ready {
		return nil, fmt.Errorf("clients not ready")
	}

	return c, nil
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

// Check the status of the Execution and Beacon Node(s) and provision the API with them
func (c *HyperdriveClient) checkClientStatus() (bool, error) {
	// Check if the primary clients are up, synced, and able to respond to requests - if not, forces the use of the fallbacks for this command
	response, err := c.Api.Service.ClientStatus()
	if err != nil {
		return false, fmt.Errorf("error checking client status: %w", err)
	}

	ecMgrStatus := response.Data.EcManagerStatus
	bcMgrStatus := response.Data.BcManagerStatus

	// Primary EC and CC are good
	if ecMgrStatus.PrimaryClientStatus.IsSynced && bcMgrStatus.PrimaryClientStatus.IsSynced {
		//c.SetClientStatusFlags(true, false)
		return true, nil
	}

	// Get the status messages
	primaryEcStatus := getClientStatusString(ecMgrStatus.PrimaryClientStatus)
	primaryBcStatus := getClientStatusString(bcMgrStatus.PrimaryClientStatus)
	fallbackEcStatus := getClientStatusString(ecMgrStatus.FallbackClientStatus)
	fallbackBcStatus := getClientStatusString(bcMgrStatus.FallbackClientStatus)

	// Check the fallbacks if enabled
	if ecMgrStatus.FallbackEnabled && bcMgrStatus.FallbackEnabled {
		// Fallback EC and CC are good
		if ecMgrStatus.FallbackClientStatus.IsSynced && bcMgrStatus.FallbackClientStatus.IsSynced {
			fmt.Printf("%sNOTE: primary clients are not ready, using fallback clients...\n\tPrimary EC status: %s\n\tPrimary CC status: %s%s\n\n", terminal.ColorYellow, primaryEcStatus, primaryBcStatus, terminal.ColorReset)
			//c.SetClientStatusFlags(true, true)
			return true, nil
		}

		// Both pairs aren't ready
		fmt.Printf("Error: neither primary nor fallback client pairs are ready.\n\tPrimary EC status: %s\n\tFallback EC status: %s\n\tPrimary CC status: %s\n\tFallback CC status: %s\n", primaryEcStatus, fallbackEcStatus, primaryBcStatus, fallbackBcStatus)
		return false, nil
	}

	// Primary isn't ready and fallback isn't enabled
	fmt.Printf("Error: primary client pair isn't ready and fallback clients aren't enabled.\n\tPrimary EC status: %s\n\tPrimary CC status: %s\n", primaryEcStatus, primaryBcStatus)
	return false, nil
}

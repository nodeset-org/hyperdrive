package client

import (
	"fmt"
	"log/slog"
	"net/http/httptrace"

	docker "github.com/docker/docker/client"
	"github.com/fatih/color"
	csclient "github.com/nodeset-org/hyperdrive-constellation/client"
	csconfig "github.com/nodeset-org/hyperdrive-constellation/shared/config"
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

// Create new Hyperdrive client from CLI context
func NewHyperdriveClientFromCtx(c *cli.Context) (*HyperdriveClient, error) {
	hdCtx := context.GetHyperdriveContext(c)
	logger := log.NewTerminalLogger(hdCtx.DebugEnabled, terminalLogColor).With(slog.String(log.OriginKey, config.HyperdriveDaemonRoute))

	// Create the tracer if required
	var tracer *httptrace.ClientTrace
	if hdCtx.HttpTraceFile != nil {
		var err error
		tracer, err = createTracer(hdCtx.HttpTraceFile, logger)
		if err != nil {
			logger.Error("Error creating HTTP trace", log.Err(err))
		}
	}

	// Make the client
	hdClient := &HyperdriveClient{
		Context: hdCtx,
		Logger:  logger,
	}

	// Get the API URL
	url := hdCtx.ApiUrl
	if url == nil {
		// Load the config to get the API port
		cfg, _, err := hdClient.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading config: %w", err)
		}

		url, err = url.Parse(fmt.Sprintf("http://localhost:%d/%s", cfg.Hyperdrive.ApiPort.Value, config.HyperdriveApiClientRoute))
		if err != nil {
			return nil, fmt.Errorf("error parsing Hyperdrive API URL: %w", err)
		}
	}
	hdClient.Api = client.NewApiClient(url, logger, tracer)
	return hdClient, nil
}

// Create new Stakewise client from CLI context
// Only use this function from commands that may work if the Daemon service doesn't exist
func NewStakewiseClientFromCtx(c *cli.Context, hdClient *HyperdriveClient) (*StakewiseClient, error) {
	hdCtx := context.GetHyperdriveContext(c)
	logger := log.NewTerminalLogger(hdCtx.DebugEnabled, terminalLogColor).With(slog.String(log.OriginKey, swconfig.ModuleName))

	// Create the tracer if required
	var tracer *httptrace.ClientTrace
	if hdCtx.HttpTraceFile != nil {
		var err error
		tracer, err = createTracer(hdCtx.HttpTraceFile, logger)
		if err != nil {
			logger.Error("Error creating HTTP trace", log.Err(err))
		}
	}

	// Make the client
	swClient := &StakewiseClient{
		Context: hdCtx,
		Logger:  logger,
	}

	// Get the API URL
	url := hdCtx.ApiUrl
	if url == nil {
		var err error
		url, err = url.Parse(fmt.Sprintf("http://localhost:%d/%s", hdClient.cfg.Stakewise.ApiPort.Value, swconfig.ApiClientRoute))
		if err != nil {
			return nil, fmt.Errorf("error parsing StakeWise API URL: %w", err)
		}
	} else {
		host := fmt.Sprintf("%s://%s:%d/%s", url.Scheme, url.Hostname(), hdClient.cfg.Stakewise.ApiPort.Value, swconfig.ApiClientRoute)
		var err error
		url, err = url.Parse(host)
		if err != nil {
			return nil, fmt.Errorf("error parsing StakeWise API URL: %w", err)
		}
	}

	swClient.Api = swclient.NewApiClient(url, logger, tracer)
	return swClient, nil
}

// Create new Constellation client from CLI context
// Only use this function from commands that may work if the Daemon service doesn't exist
func NewConstellationClientFromCtx(c *cli.Context, hdClient *HyperdriveClient) (*ConstellationClient, error) {
	hdCtx := context.GetHyperdriveContext(c)
	logger := log.NewTerminalLogger(hdCtx.DebugEnabled, terminalLogColor).With(slog.String(log.OriginKey, csconfig.ModuleName))

	// Create the tracer if required
	var tracer *httptrace.ClientTrace
	if hdCtx.HttpTraceFile != nil {
		var err error
		tracer, err = createTracer(hdCtx.HttpTraceFile, logger)
		if err != nil {
			logger.Error("Error creating HTTP trace", log.Err(err))
		}
	}

	// Make the client
	csClient := &ConstellationClient{
		Context: hdCtx,
		Logger:  logger,
	}

	// Get the API URL
	url := hdCtx.ApiUrl
	if url == nil {
		var err error
		url, err = url.Parse(fmt.Sprintf("http://localhost:%d/%s", hdClient.cfg.Constellation.ApiPort.Value, csconfig.ApiClientRoute))
		if err != nil {
			return nil, fmt.Errorf("error parsing Constellation API URL: %w", err)
		}
	} else {
		host := fmt.Sprintf("%s://%s:%d/%s", url.Scheme, url.Hostname(), hdClient.cfg.Constellation.ApiPort.Value, csconfig.ApiClientRoute)
		var err error
		url, err = url.Parse(host)
		if err != nil {
			return nil, fmt.Errorf("error parsing Constellation API URL: %w", err)
		}
	}

	csClient.Api = csclient.NewApiClient(url, logger, tracer)
	return csClient, nil
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

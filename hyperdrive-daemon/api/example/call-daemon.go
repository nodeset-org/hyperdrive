package example

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

const (
	apiContainerName string = "api"
	binaryPath       string = "/go/bin/rocketpool"
)

// ===============
// === Factory ===
// ===============

type callDaemonContextFactory struct {
	handler *ExampleHandler
}

func (f *callDaemonContextFactory) Create(args url.Values) (*callDaemonContext, error) {
	c := &callDaemonContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.GetStringFromValues("cmd", args, &c.cmd),
	}
	switch c.cmd {
	case "node-fee", "rpl-price", "stats", "timezone-map", "dao-proposals", "download-rewards-file", "latest-delegate":
		break
	default:
		return nil, fmt.Errorf("%s is not a valid command for the network subroutine", c.cmd)
	}
	return c, errors.Join(inputErrs...)
}

func (f *callDaemonContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*callDaemonContext, api.CallDaemonData](
		router, "call-daemon", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type callDaemonContext struct {
	handler *ExampleHandler
	cfg     *config.HyperdriveConfig
	// rp      *rocketpool.RocketPool

	cmd string
	// amountWei *big.Int
	// lot       *auction.AuctionLot
	// pSettings *protocol.ProtocolDaoSettings
}

func (c *callDaemonContext) PrepareData(data *api.CallDaemonData) error {
	// Get the container name for the Smartnode daemon
	sp := c.handler.serviceProvider
	sncfg := sp.GetConfig().SmartnodeConfig
	projectName := sncfg.Smartnode.ProjectName.Value.(string)
	containerName := fmt.Sprintf("%s_%s", projectName, apiContainerName)

	// Construct the command
	args := []string{
		"exec",
		containerName,
		binaryPath,
		"api",
		"network",
		c.cmd,
	}
	cmd := exec.Command("docker", args...)

	// If it's in debug mode, print it and exit
	if c.handler.isDebug {
		fmt.Printf("[Debug] Command: %s\n", cmd.String())
		return nil
	}

	// Run it
	response, err := cmd.Output()
	if err != nil {
		exitError, isExitError := err.(*exec.ExitError)
		if isExitError {
			data.Error = fmt.Sprintf("smartnode daemon exited with status code %d:\n%s", exitError.ExitCode(), string(exitError.Stderr))
		} else {
			return fmt.Errorf("error running command: %w", err)
		}
	} else {
		data.Response = string(response)
	}
	return nil
}

package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

// Terminate the Hyperdrive service
func terminateService(c *cli.Context) error {
	// Prompt for confirmation
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm(fmt.Sprintf("%sWARNING: Are you sure you want to terminate the Hyperdrive service? Any staking minipools will be penalized, your Execution client and Beacon node chain databases will be deleted, you will lose ALL of your sync progress, and you will lose your Prometheus metrics database!\nAfter doing this, you will have to **reinstall** Hyperdrive uses `hyperdrive service install -d` in order to use it again.%s", terminal.ColorRed, terminal.ColorReset))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Stop service
	return hd.TerminateService(getComposeFiles(c), hd.Context.ConfigPath)
}

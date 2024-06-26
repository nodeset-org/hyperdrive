package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

var (
	includeVolumesFlag = &cli.BoolFlag{
		Name:    "include-volumes",
		Aliases: []string{"v"},
		Usage:   fmt.Sprintf("Remove volumes as well. %sThis will delete your EC and BN chain data and you will have to resync from scratch if you include this!%s", terminal.ColorRed, terminal.ColorReset),
	}
)

// Stop the Hyperdrive service, removing the Docker artifacts
func downService(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Prompt for confirmation
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Are you sure you want to shut down the Hyperdrive service and remove the Docker artifacts?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Take it down
	includeVolumes := c.Bool(includeVolumesFlag.Name)
	return hd.DownService(getComposeFiles(c), includeVolumes)
}

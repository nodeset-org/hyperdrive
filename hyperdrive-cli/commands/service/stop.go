package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Pause the Hyperdrive service
func stopService(c *cli.Context) error {
	// Get RP client
	rp := client.NewClientFromCtx(c)

	// Prompt for confirmation
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Are you sure you want to pause the Hyperdrive service?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Pause service
	return rp.PauseService(getComposeFiles(c))
}

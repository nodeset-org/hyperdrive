package service

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

// View the Hyperdrive service stats
func serviceStats(c *cli.Context) error {
	// Get RP client
	hd := client.NewClientFromCtx(c)

	// Print service stats
	return hd.PrintServiceStats(getComposeFiles(c))
}

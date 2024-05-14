package service

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

// View the Hyperdrive service compose config
func serviceCompose(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Print service compose config
	return hd.PrintServiceCompose(getComposeFiles(c))
}

package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// View the Hyperdrive service status
func serviceStatus(c *cli.Context) error {
	// Get RP client
	hd := client.NewClientFromCtx(c)

	// Get the config
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("Error loading configuration: %w", err)
	}

	// Print what network we're on
	err = utils.PrintNetwork(cfg.Network.Value, isNew)
	if err != nil {
		return err
	}

	// Print service status
	return hd.PrintServiceStatus(getComposeFiles(c))
}

package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/urfave/cli/v2"
)

var (
	downIncludeVolumesFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "include-volumes",
		Aliases: []string{"v"},
		Usage:   "Include volumes in the down command, so all volumes (including module volumes) will be deleted as well",
		Value:   false,
	}
)

// Delete the Hyperdrive service, stopping and deleting the Docker containers for all modules
func downService(c *cli.Context) error {
	hd, err := utils.NewHyperdriveManagerFromCtx(c)
	if err != nil {
		return err
	}

	// Load the current settings
	settings, isNew, err := hd.LoadMainSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}
	if isNew {
		fmt.Println("Hyperdrive has not been configured yet. Please run 'hyperdrive service configure' first.")
		return nil
	}

	// Remove the service
	err = hd.DownService(settings, c.Bool(downIncludeVolumesFlag.Name))
	if err != nil {
		return fmt.Errorf("error deleting service: %w", err)
	}
	return nil
}

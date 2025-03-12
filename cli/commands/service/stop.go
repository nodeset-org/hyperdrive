package service

import (
	"fmt"

	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/urfave/cli/v2"
)

// Stop the Hyperdrive service, stopping the Docker containers for all modules
func stopService(c *cli.Context) error {
	hd, err := cliutils.NewHyperdriveManagerFromCtx(c)
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

	// Stop the service
	err = hd.StopService(settings)
	if err != nil {
		return fmt.Errorf("error stopping service: %w", err)
	}
	return nil
}

package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/cli/utils"
	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
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

	// Load the modules
	modLoadResults, err := hd.LoadModules()
	if err != nil {
		return fmt.Errorf("error loading modules: %w", err)
	}
	failures := false
	for _, result := range modLoadResults {
		if result.LoadError != nil {
			fmt.Printf("WARNING: Module %s failed to load: %s\n", result.Info.Descriptor.Name, result.LoadError.Error())
			failures = true
		}
	}
	if failures {
		fmt.Println("The above modules cannot be deleted (if they are running) until their problems are resolved.")
		if !(c.Bool(cliutils.YesFlag.Name) || cliutils.Confirm("Are you sure you want to continue?")) {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// No pending settings, so load the main settings
	settings, isNew, err := hd.LoadMainSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}

	// Check if the config is new (hasn't been installed before)
	if isNew {
		fmt.Println("Hyperdrive has not been configured yet. Please run 'hyperdrive service configure' first.")
		return nil
	}

	// Disable modules that failed to load
	for _, result := range modLoadResults {
		if result.LoadError == nil {
			continue
		}
		modInstance, exists := settings.Modules[result.Info.Descriptor.GetFullyQualifiedModuleName()]
		if !exists {
			continue
		}
		modInstance.Enabled = false
	}

	// Stop the service
	err = hd.DownService(settings, c.Bool(downIncludeVolumesFlag.Name))
	if err != nil {
		return fmt.Errorf("error deleting service: %w", err)
	}
	return nil
}

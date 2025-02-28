package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/cli/client"
	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/urfave/cli/v2"
)

// Start the Hyperdrive service, starting the Docker containers for all modules
func startService(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
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
		fmt.Println("The above modules will be disabled if you proceed until their problems are resolved.")
		if !(c.Bool(cliutils.YesFlag.Name) || cliutils.Confirm("Are you sure you want to continue?")) {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Load the settings from disk
	pendingSettings, noPendingSettings, err := hd.LoadPendingSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading pending settings: %w", err)
	}
	if noPendingSettings {
		pendingSettings = nil
	}
	currentSettings, noCurrentSettings, err := hd.LoadMainSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}

	// Check if the config is new (hasn't been installed before)
	if noCurrentSettings {
		fmt.Println("Hyperdrive has not been configured yet. Please run 'hyperdrive service configure' first.")
		return nil
	}

	// Disable modules that failed to load
	for _, result := range modLoadResults {
		if result.LoadError == nil {
			continue
		}
		modInstance, exists := currentSettings.Modules[result.Info.Descriptor.GetFullyQualifiedModuleName()]
		if !exists {
			continue
		}
		modInstance.Enabled = false

		if pendingSettings == nil {
			continue
		}
		modInstance, exists = pendingSettings.Modules[result.Info.Descriptor.GetFullyQualifiedModuleName()]
		if !exists {
			continue
		}
		modInstance.Enabled = false
	}

	// Start the service
	err = hd.StartService(currentSettings, pendingSettings)
	if err != nil {
		return fmt.Errorf("error starting service: %w", err)
	}
	return nil
}

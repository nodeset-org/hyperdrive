package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/cli/utils"
	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/nodeset-org/hyperdrive/management"
	"github.com/urfave/cli/v2"
)

// Start the Hyperdrive service, starting the Docker containers for all modules
func startService(c *cli.Context) error {
	hd, err := utils.NewHyperdriveManagerFromCtx(c)
	if err != nil {
		return err
	}

	// Load the modules
	err = hd.LoadModules()
	if err != nil {
		return fmt.Errorf("error loading modules: %w", err)
	}

	// Warn about broken modules
	for _, result := range hd.BrokenModules {
		if result.ConfigurationLoadError != nil {
			fmt.Printf("Skipping module %s because it failed to load: %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.ConfigurationLoadError)
		} else if result.GlobalAdapterContainerStatus != management.ContainerStatus_Running {
			fmt.Printf("Skipping module %s because its global adapter container could not start\n", result.Descriptor.GetFullyQualifiedModuleName())
		} else if result.GlobalAdapterRuntimeFileError != nil {
			fmt.Printf("Skipping module %s because its global adapter container file could not be instantiated: %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.GlobalAdapterRuntimeFileError)
		} else if result.DescriptorLoadError != nil {
			fmt.Printf("Skipping module %s because its descriptor could not be loaded: %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.DescriptorLoadError)
		} else {
			fmt.Printf("Skipping module %s because it could not be loaded for an unknown reason\n", result.Descriptor.GetFullyQualifiedModuleName())
		}
	}
	if len(hd.BrokenModules) > 0 {
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
	for _, result := range hd.BrokenModules {
		modInstance, exists := currentSettings.Modules[result.Descriptor.GetFullyQualifiedModuleName()]
		if !exists {
			continue
		}
		modInstance.Enabled = false

		if pendingSettings == nil {
			continue
		}
		modInstance, exists = pendingSettings.Modules[result.Descriptor.GetFullyQualifiedModuleName()]
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

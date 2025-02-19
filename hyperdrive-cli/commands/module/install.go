package module

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

// Install a module
func installModule(c *cli.Context, moduleFile string) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Check if we have permissions to install the module

	// Install the module
	mgr := hd.GetModuleManager()
	err = mgr.InstallModule(moduleFile)
	if err != nil {
		return fmt.Errorf("error installing module: %w", err)
	}
	fmt.Println("Module installed successfully.")

	// Start the global adapters
	results, err := mgr.LoadModuleInfo(true)
	if err != nil {
		return fmt.Errorf("error loading module info: %w", err)
	}
	for _, result := range results {
		if result.LoadError != nil {
			fmt.Printf("WARNING: Module %s failed to load: %s\n", result.Info.Descriptor.Name, result.LoadError.Error())
		}
	}
	return nil
}

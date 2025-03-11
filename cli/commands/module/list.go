package module

import (
	"fmt"

	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/nodeset-org/hyperdrive/management"
	"github.com/urfave/cli/v2"
)

func listModules(c *cli.Context) error {
	hd, err := cliutils.NewHyperdriveManagerFromCtx(c)
	if err != nil {
		return err
	}

	// Get the list of modules
	err = hd.LoadModules()
	if err != nil {
		return fmt.Errorf("error loading modules: %w", err)
	}
	if len(hd.BrokenModules) == 0 && len(hd.HealthyModules) == 0 {
		fmt.Println("No modules are currently installed.")
		return nil
	}

	// Print the successfully loaded modules
	if len(hd.HealthyModules) > 0 {
		fmt.Println("Successfully loaded modules:")
	}
	for _, result := range hd.HealthyModules {
		descriptor := result.Descriptor
		fmt.Printf("\t%s - %s (%s)\n", descriptor.Author, descriptor.Name, descriptor.Version)
	}
	fmt.Println()

	// Print the failed modules
	if len(hd.BrokenModules) > 0 {
		fmt.Println("Modules that failed to load:")
		for _, result := range hd.BrokenModules {
			if result.ConfigurationLoadError != nil {
				fmt.Printf("\t%s: configuration load failure - %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.ConfigurationLoadError)
			} else if result.GlobalAdapterContainerStatus != management.ContainerStatus_Running {
				fmt.Printf("\t%s: global adapter failed to start\n", result.Descriptor.GetFullyQualifiedModuleName())
			} else if result.GlobalAdapterRuntimeFileError != nil {
				fmt.Printf("\t%s: global adapter file failure - %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.GlobalAdapterRuntimeFileError)
			} else if result.DescriptorLoadError != nil {
				fmt.Printf("S\t%s: descriptor load failure - %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.DescriptorLoadError)
			} else {
				fmt.Printf("\t%s: unknown reason\n", result.Descriptor.GetFullyQualifiedModuleName())
			}
		}
		fmt.Println()
	}

	return nil
}

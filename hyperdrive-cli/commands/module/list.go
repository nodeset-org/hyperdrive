package module

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/urfave/cli/v2"
)

func listModules(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Load the current config
	var projectName string
	settings, exists, err := hd.LoadMainSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading main config file: %w", err)
	}
	if !exists {
		// Use the default project name
		projectName = hd.GetHyperdriveConfiguration().ProjectName.Default
		fmt.Printf("NOTE: Hyperdrive has not been configured yet, using the default project name (%s)\n", projectName)
	} else {
		projectName = settings.ProjectName
	}

	// Get the list of modules
	results, err := hd.LoadModules(projectName)
	if err != nil {
		return fmt.Errorf("error loading modules: %w", err)
	}
	if len(results) == 0 {
		fmt.Println("No modules are currently installed.")
		return nil
	}

	// Check each one's status
	failedModules := []*shared.ModuleInfoLoadResult{}
	succeededModules := []*shared.ModuleInfoLoadResult{}
	for _, result := range results {
		if result.LoadError != nil {
			failedModules = append(failedModules, result)
			continue
		}
		succeededModules = append(succeededModules, result)
	}

	// Print the successfully loaded modules
	if len(succeededModules) > 0 {
		fmt.Println("Successfully loaded modules:")
		for _, result := range succeededModules {
			descriptor := result.Info.Descriptor
			fmt.Printf("\t%s - %s (%s)\n", descriptor.Author, descriptor.Name, descriptor.Version)
		}
		fmt.Println()
	}

	// Print the failed modules
	if len(failedModules) > 0 {
		fmt.Println("Modules that failed to load:")
		for _, result := range failedModules {
			fmt.Printf("\t%s: %s\n", result.Info.Descriptor.GetFullyQualifiedModuleName(), result.LoadError)
		}
		fmt.Println()
	}

	return nil
}

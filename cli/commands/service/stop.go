package service

import (
	"context"
	"fmt"

	"github.com/nodeset-org/hyperdrive/cli/client"
	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/urfave/cli/v2"
)

// Stop the Hyperdrive service, stopping the Docker containers for all modules
func stopService(c *cli.Context) error {
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
		fmt.Println("The above modules cannot be stopped (if they are running) until their problems are resolved.")
		if !(c.Bool(cliutils.YesFlag.Name) || cliutils.Confirm("Are you sure you want to continue?")) {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	return nil
}

// For each project, have the project adapter stop the services that are marked for a restart in the provided settings.
func stopServicesMarkedForRestart(
	modMgr *utils.ModuleManager,
	hdSettings *hdconfig.HyperdriveSettings,
	infos map[string]*modconfig.ModuleInfo,
) error {
	for _, info := range infos {
		// Get the list of services to stop
		fqmn := info.Descriptor.GetFullyQualifiedModuleName()
		services := hdSettings.Modules[fqmn].Restart
		if len(services) == 0 {
			continue
		}

		// Stop the services
		pac, err := modMgr.GetProjectAdapterClient(hdSettings.ProjectName, fqmn)
		if err != nil {
			return fmt.Errorf("error getting project adapter client for module [%s]: %w", info.Descriptor.Name, err)
		}
		err = pac.Stop(context.Background(), hdSettings.ProjectName+"-"+string(info.Descriptor.Shortcut), services)
		if err != nil {
			return fmt.Errorf("error stopping module [%s]: %w", info.Descriptor.Name, err)
		}
	}
	return nil
}

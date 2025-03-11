package module

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nodeset-org/hyperdrive/cli/utils"
	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/nodeset-org/hyperdrive/management"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
	"github.com/urfave/cli/v2"
)

// Install a module
func installModule(c *cli.Context, moduleFile string) error {
	hd, err := utils.NewHyperdriveManagerFromCtx(c)
	if err != nil {
		return err
	}

	// Check if we have permissions to install the module
	mgr := hd.GetModuleManager()
	modDir := mgr.GetModuleSystemDir()
	testFile := filepath.Join(modDir, "test")
	_, err = os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY, management.ModuleFileMode)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			// We need privileges to install modules so try to escalate and re-run
			escalationCmd, err := command.GetEscalationCommand()
			if err != nil {
				return fmt.Errorf("escalated privileges are required to install this module but the escalation command could not be found: %w", err)
			}
			appPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("escalated privileges are required to install this module but error getting executable path: %w", err)
			}
			fmt.Println("Privilege escalation is required to install modules to the system directory.")
			args := []string{
				escalationCmd,
				appPath,
				"--" + cliutils.AllowRootFlag.Name,
				"--" + cliutils.UserDirPathFlag.Name,
				hd.Context.UserDirPath,
				"--" + cliutils.SystemDirPathFlag.Name,
				hd.Context.SystemDirPath,
				"module",
				"install",
				moduleFile,
			}
			cmd := command.NewCommand(strings.Join(args, " "))
			cmd.SetStdin(os.Stdin)
			cmd.SetStdout(os.Stdout)
			cmd.SetStderr(os.Stderr)
			return cmd.Run()
		} else {
			return fmt.Errorf("error testing installation permissions: %w", err)
		}
	}

	// Install the module
	err = mgr.InstallModule(moduleFile)
	if err != nil {
		return fmt.Errorf("error installing module: %w", err)
	}
	fmt.Println("Module installed successfully.")

	// Start the global adapters
	err = hd.LoadModules()
	if err != nil {
		return fmt.Errorf("error loading module info: %w", err)
	}
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
		fmt.Println("The above modules will be disabled until their load errors are resolved.")
		if !c.Bool(utils.YesFlag.Name) {
			fmt.Println("Press any key to continue.")
			_, _ = fmt.Scanln()
		}
	}
	return nil
}

package module

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	cliutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	hdutils "github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
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
	mgr := hd.GetModuleManager()
	modDir := mgr.GetModuleSystemDir()
	testFile := filepath.Join(modDir, "test")
	_, err = os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY, hdutils.ModuleFileMode)
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

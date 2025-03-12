package service

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/nodeset-org/hyperdrive/management"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
	"github.com/urfave/cli/v2"
)

func terminateService(c *cli.Context) error {
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

	// Print the warning message
	fmt.Println("WARNING: This will stop and delete all Docker containers (including module servces), remove all volumes, and delete all data (including your Hyperdrive configuration and all module data). This action cannot be undone!")
	if !(c.Bool(cliutils.YesFlag.Name) || cliutils.Confirm("Are you SURE you want to terminate Hyperdrive?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Delete the services
	err = hd.DownService(settings, true)
	if err != nil {
		return fmt.Errorf("error deleting services: %w", err)
	}

	// Delete the data directory
	err = purgeDataImpl(hd, settings)
	if err != nil {
		return err
	}

	// Delete the user folder
	err = deleteUserFolderImpl(hd)
	if err != nil {
		return err
	}

	return nil
}

// Delete the user folder. Reserved for the `delete-user-folder` command.
func deleteUserFolder(c *cli.Context) error {
	hd, err := cliutils.NewHyperdriveManagerFromCtx(c)
	if err != nil {
		return err
	}

	// Load the current settings
	_, isNew, err := hd.LoadMainSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}
	if isNew {
		return nil
	}

	// Delete the user folder
	return deleteUserFolderImpl(hd)
}

// Internal implementation of deleting the user folder that can be escalated
// TODO: escalation shouldn't be needed here if the data directory was already purged.
func deleteUserFolderImpl(hd *management.HyperdriveManager) error {
	err := hd.DeleteUserFolder(hd.Context)
	if err == nil {
		return nil
	}

	if errors.Is(err, fs.ErrPermission) {
		// We don't have permission so escalate and try again
		escalationCmd, err := command.GetEscalationCommand()
		if err != nil {
			return fmt.Errorf("escalated privileges are required to delete the user folder but the escalation command could not be found: %w", err)
		}
		appPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("escalated privileges are required to delete the user folder but error getting executable path: %w", err)
		}
		fmt.Println("Privilege escalation is required to delete the user folder.")
		args := []string{
			escalationCmd,
			appPath,
			"--" + cliutils.AllowRootFlag.Name,
			"--" + cliutils.UserDirPathFlag.Name,
			hd.Context.UserDirPath,
			"--" + cliutils.SystemDirPathFlag.Name,
			hd.Context.SystemDirPath,
			"service",
			"delete-user-folder",
		}
		cmd := command.NewCommand(strings.Join(args, " "))
		cmd.SetStdin(os.Stdin)
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)
		return cmd.Run()
	}

	return fmt.Errorf("error deleting user folder: %w", err)
}

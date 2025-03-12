package service

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/nodeset-org/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive/management"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
	"github.com/urfave/cli/v2"
)

var (
	noRestartFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "no-restart",
		Aliases: []string{"nr"},
		Usage:   "Don't restart the containers after purging the data",
		Value:   false,
	}
)

// Purge the data from the Hyperdrive service, stopping and deleting the Docker containers for all modules
func purgeData(c *cli.Context) error {
	hd, err := cliutils.NewHyperdriveManagerFromCtx(c)
	if err != nil {
		return err
	}

	// Load the current settings
	settings, isNew, err := hd.LoadMainSettingsFile()
	if err != nil {
		return err
	}
	if isNew {
		return nil
	}

	// Stop the service
	err = hd.StopService(settings)
	if err != nil {
		return fmt.Errorf("error stopping service: %w", err)
	}

	// Try to purge the data
	err = purgeDataImpl(hd, settings)
	if err != nil {
		return err
	}

	// Restart the service if not disabled
	if !c.Bool(noRestartFlag.Name) {
		err = hd.StartService(settings, nil) // Ignore pending settings for this command
		if err != nil {
			return fmt.Errorf("error starting service: %w", err)
		}
	}
	return nil
}

// Purge the data folder without interacting with Docker. Reserved for the `purge-data` command.
func purgeDataExclusive(c *cli.Context) error {
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
		return nil
	}

	// Purge the data
	return purgeDataImpl(hd, settings)
}

// Internal implementation of purging the data that can be escalated
func purgeDataImpl(hd *management.HyperdriveManager, settings *config.HyperdriveSettings) error {
	err := hd.PurgeData(settings)
	if err == nil {
		return nil
	}

	if errors.Is(err, fs.ErrPermission) {
		// We don't have permission so escalate and try again
		escalationCmd, err := command.GetEscalationCommand()
		if err != nil {
			return fmt.Errorf("escalated privileges are required to purge data but the escalation command could not be found: %w", err)
		}
		appPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("escalated privileges are required to purge data but error getting executable path: %w", err)
		}
		fmt.Println("Privilege escalation is required to purge data.")
		args := []string{
			escalationCmd,
			appPath,
			"--" + cliutils.AllowRootFlag.Name,
			"--" + cliutils.UserDirPathFlag.Name,
			hd.Context.UserDirPath,
			"--" + cliutils.SystemDirPathFlag.Name,
			hd.Context.SystemDirPath,
			"service",
			"purge-data-exclusive",
		}
		cmd := command.NewCommand(strings.Join(args, " "))
		cmd.SetStdin(os.Stdin)
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)
		return cmd.Run()
	}

	return fmt.Errorf("error purging data: %w", err)
}

package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

// Destroy and resync the Execution client from scratch
func resyncExecutionClient(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Get the config
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return err
	}
	if isNew {
		return fmt.Errorf("Settings file not found. Please run `hyperdrive service config` to set up Hyperdrive.")
	}

	fmt.Println("This will delete the chain data of your primary Execution client and resync it from scratch.")
	fmt.Printf("%sYou should only do this if your Execution client has failed and can no longer start or sync properly.\nThis is meant to be a last resort.%s\n", terminal.ColorYellow, terminal.ColorReset)

	// Prompt for confirmation
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm(fmt.Sprintf("%sAre you SURE you want to delete and resync your main Execution client from scratch? This cannot be undone!%s", terminal.ColorRed, terminal.ColorReset))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Stop Execution
	executionContainerName := cfg.Hyperdrive.GetDockerArtifactName(string(config.ContainerID_ExecutionClient))
	fmt.Printf("Stopping %s...\n", executionContainerName)
	err = hd.StopContainer(executionContainerName)
	if err != nil {
		fmt.Printf("%sWARNING: Stopping main Execution container failed: %s%s\n", terminal.ColorYellow, err.Error(), terminal.ColorReset)
	}

	// Get Execution volume name
	volume, err := hd.GetClientVolumeName(executionContainerName, clientDataVolumeName)
	if err != nil {
		return fmt.Errorf("Error getting Execution client volume name: %w", err)
	}

	// Remove the EC
	fmt.Printf("Deleting %s...\n", executionContainerName)
	err = hd.RemoveContainer(executionContainerName)
	if err != nil {
		return fmt.Errorf("Error deleting main Execution client container: %w", err)
	}

	// Delete the EC volume
	fmt.Printf("Deleting volume %s...\n", volume)
	err = hd.DeleteVolume(volume)
	if err != nil {
		return fmt.Errorf("Error deleting volume: %w", err)
	}

	// Restart Hyperdrive
	fmt.Printf("Rebuilding %s and restarting Hyperdrive...\n", executionContainerName)
	err = startService(c, true)
	if err != nil {
		return fmt.Errorf("Error starting Hyperdrive: %s", err)
	}

	fmt.Printf("\nDone! Your main Execution client is now resyncing. You can follow its progress with `hyperdrive service logs ec`.\n")
	return nil
}

package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/urfave/cli/v2"
)

// Destroy and resync the Beacon Node from scratch
func resyncBeaconNode(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Get the merged config
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return err
	}
	if isNew {
		return fmt.Errorf("Settings file not found. Please run `hyperdrive service config` to set up Hyperdrive.")
	}

	fmt.Println("This will delete the chain data of your Beacon Node and resync it from scratch.")
	fmt.Printf("%sYou should only do this if your Beacon Node has failed and can no longer start or sync properly.\nThis is meant to be a last resort.%s\n\n", terminal.ColorYellow, terminal.ColorReset)

	// Check the client mode
	if cfg.Hyperdrive.ClientMode.Value == config.ClientMode_External {
		fmt.Println("You use an externally-managed Beacon Node. Hyperdrive cannot resync it for you.")
		return nil
	}

	// Get the current checkpoint sync URL
	checkpointSyncUrl := cfg.Hyperdrive.LocalBeaconClient.CheckpointSyncProvider.Value
	if checkpointSyncUrl == "" {
		fmt.Printf("%sYou do not have a checkpoint sync provider configured.\nIf you have active validators, they %swill be considered offline and will lose ETH%s%s until your Beacon Node finishes syncing.\nWe strongly recommend you configure a checkpoint sync provider with `hyperdrive service config` so it syncs instantly before running this.%s\n\n", terminal.ColorRed, terminal.ColorBold, terminal.ColorReset, terminal.ColorRed, terminal.ColorReset)
	} else {
		fmt.Printf("You have a checkpoint sync provider configured (%s).\nYour Beacon Node will use it to sync to the head of the Beacon Chain instantly after being rebuilt.\n\n", checkpointSyncUrl)
	}

	// Prompt for confirmation
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm(fmt.Sprintf("%sAre you SURE you want to delete and resync your main Beacon Node from scratch? This cannot be undone!%s", terminal.ColorRed, terminal.ColorReset))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Stop the BN
	beaconContainerName := cfg.Hyperdrive.GetDockerArtifactName(string(config.ContainerID_BeaconNode))
	fmt.Printf("Stopping %s...\n", beaconContainerName)
	err = hd.StopContainer(beaconContainerName)
	if err != nil {
		fmt.Printf("%sWARNING: Stopping Beacon Node container failed: %s%s\n", terminal.ColorYellow, err.Error(), terminal.ColorReset)
	}

	// Get the BN volume name
	volume, err := hd.GetClientVolumeName(beaconContainerName, clientDataVolumeName)
	if err != nil {
		return fmt.Errorf("Error getting Beacon Node volume name: %w", err)
	}

	// Remove the BN
	fmt.Printf("Deleting %s...\n", beaconContainerName)
	err = hd.RemoveContainer(beaconContainerName)
	if err != nil {
		return fmt.Errorf("Error deleting Beacon Node container: %w", err)
	}

	// Delete the BN volume
	fmt.Printf("Deleting volume %s...\n", volume)
	err = hd.DeleteVolume(volume)
	if err != nil {
		return fmt.Errorf("Error deleting volume: %w", err)
	}

	// Restart Hyperdrive
	fmt.Printf("Rebuilding %s and restarting Hyperdrive...\n", beaconContainerName)
	err = startService(c, StartMode_NoUpdate)
	if err != nil {
		return fmt.Errorf("Error starting Hyperdrive: %s", err)
	}

	fmt.Printf("\nDone! Your Beacon Node is now resyncing. You can follow its progress with `hyperdrive service logs bn`.\n")
	return nil
}

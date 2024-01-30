package service

import (
	"fmt"
	"strings"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/urfave/cli/v2"
)

// Settings
const (
	clientDataVolumeName string = "/ethclient"
	dataFolderVolumeName string = "/.hyperdrive/data"

	PruneFreeSpaceRequired uint64 = 50 * 1024 * 1024 * 1024
	dockerImageRegex       string = ".*/(?P<image>.*):.*"
)

// Get the compose file paths for a CLI context
func getComposeFiles(c *cli.Context) []string {
	return c.StringSlice(utils.ComposeFileFlag.Name)
}

// Handle a network change by terminating the service, deleting everything, and starting over
// TODO
func changeNetworks(c *cli.Context, hd *client.HyperdriveClient) error {
	return fmt.Errorf("NYI")

	// Stop all of the containers
	fmt.Print("Stopping containers... ")
	err := hd.PauseService(getComposeFiles(c))
	if err != nil {
		return fmt.Errorf("error stopping service: %w", err)
	}
	fmt.Println("done")

	// Delete the data folder
	fmt.Print("Removing data folder... ")
	_, err = hd.Api.Service.TerminateDataFolder()
	if err != nil {
		return err
	}
	fmt.Println("done")

	// Terminate the current setup
	fmt.Print("Removing old installation... ")
	err = hd.StopService(getComposeFiles(c))
	if err != nil {
		return fmt.Errorf("error terminating old installation: %w", err)
	}
	fmt.Println("done")

	// Start the service
	fmt.Print("Starting Hyperdrive... ")
	err = hd.StartService(getComposeFiles(c))
	if err != nil {
		return fmt.Errorf("error starting service: %w", err)
	}
	fmt.Println("done")

	return nil
}

// Gets the prefix specified for Hyperdrive's Docker containers
func getContainerPrefix(hd *client.HyperdriveClient) (string, error) {
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return "", err
	}
	if isNew {
		return "", fmt.Errorf("Settings file not found. Please run `hyperdrive service config` to set up Hyperdrive.")
	}

	return cfg.ProjectName.Value, nil
}

// Get the amount of free space available in the target dir
func getPartitionFreeSpace(hd *client.HyperdriveClient, targetDir string) (uint64, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return 0, fmt.Errorf("error getting partition list: %w", err)
	}
	longestPath := 0
	bestPartition := disk.PartitionStat{}
	for _, partition := range partitions {
		if strings.HasPrefix(targetDir, partition.Mountpoint) && len(partition.Mountpoint) > longestPath {
			bestPartition = partition
			longestPath = len(partition.Mountpoint)
		}
	}
	diskUsage, err := disk.Usage(bestPartition.Mountpoint)
	if err != nil {
		return 0, fmt.Errorf("error getting free disk space available: %w", err)
	}
	return diskUsage.Free, nil
}

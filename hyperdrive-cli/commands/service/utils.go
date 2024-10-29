package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
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
func changeNetworks(c *cli.Context) error {
	// Create a new Hyperdrive client - important to ensure the config is loaded from disk and isn't the stale old one
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	composeFiles := getComposeFiles(c)

	// Purge the data folder
	fmt.Print("Purging data folder... ")
	err = hd.PurgeData(composeFiles, false)
	if err != nil {
		return fmt.Errorf("error purging data folder: %w", err)
	}

	// Terminate the current setup
	fmt.Print("Removing old installation... ")
	err = hd.DownService(composeFiles, true)
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

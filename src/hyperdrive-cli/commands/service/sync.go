package service

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/urfave/cli/v2"
)

// When printing sync percents, we should avoid printing 100%.
// This function is only called if we're still syncing,
// and the `%0.2f` token will round up if we're above 99.99%.
func SyncRatioToPercent(in float64) float64 {
	return math.Min(99.99, in*100)
	// TODO: INCORPORATE THIS
}

// Settings
const (
	ethClientRecentBlockThreshold time.Duration = 5 * time.Minute
)

func printClientStatus(status *api.ClientStatus, name string) {

	if status.Error != "" {
		fmt.Printf("Your %s is unavailable (%s).\n", name, status.Error)
		return
	}

	if status.IsSynced {
		fmt.Printf("Your %s is fully synced.\n", name)
		return
	}

	fmt.Printf("Your %s is still syncing (%0.2f%%).\n", name, client.SyncRatioToPercent(status.SyncProgress))
	if strings.Contains(name, "execution") && status.SyncProgress == 0 {
		fmt.Printf("\tNOTE: your %s may not report sync progress.\n\tYou should check its logs to review it.\n", name)
	}
}

func printSyncProgress(status *api.ClientManagerStatus, name string) {

	// Print primary client status
	printClientStatus(&status.PrimaryClientStatus, fmt.Sprintf("primary %s client", name))

	if !status.FallbackEnabled {
		fmt.Printf("You do not have a fallback %s client enabled.\n", name)
		return
	}

	// A fallback is enabled, so print fallback client status
	printClientStatus(&status.FallbackClientStatus, fmt.Sprintf("fallback %s client", name))
}

func getSyncProgress(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Get the config
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("Error loading configuration: %w", err)
	}

	// Print what network we're on
	err = utils.PrintNetwork(cfg.Hyperdrive.Network.Value, isNew)
	if err != nil {
		return err
	}

	// Get node status
	status, err := hd.Api.Service.ClientStatus()
	if err != nil {
		return err
	}

	// Print EC status
	printSyncProgress(&status.Data.EcManagerStatus, "execution")

	// Print CC status
	printSyncProgress(&status.Data.BcManagerStatus, "beacon")

	// Return
	return nil
}

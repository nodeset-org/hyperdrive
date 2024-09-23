package wallet

import (
	"fmt"
	"time"

	csconfig "github.com/nodeset-org/hyperdrive-constellation/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

var (
	rebuildAllFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "rebuild-all",
		Aliases: []string{"a"},
		Usage:   "Rebuild all validator private keys for all minipools, including ones for validators that have already exited the Beacon chain",
	}
	rebuildStartIndexFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "start-index",
		Aliases: []string{"s"},
		Usage:   "The index of the validator key path to start with when rebuilding keys",
		Value:   0,
	}
	rebuildSearchLimitFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "search-limit",
		Aliases: []string{"l"},
		Usage:   "The maximum number of validator key paths to search for each key before failing",
		Value:   500,
	}
)

func rebuildValidatorKeys(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Get the all flag
	rebuildAll := c.Bool(rebuildAllFlag.Name)

	// Get the node's minipool info
	minipoolInfo, err := cs.Api.Minipool.GetPubkeys(rebuildAll)
	if err != nil {
		return fmt.Errorf("error getting node's minipool info: %w", err)
	}

	// Print some details to the user
	if len(minipoolInfo.Data.Infos) == 0 {
		fmt.Println("No minipools are eligible for rebuilding.")
		return nil
	}
	fmt.Println("The following keys will be rebuilt:")
	for _, info := range minipoolInfo.Data.Infos {
		fmt.Printf("\tMinipool: %s\n", info.Address.Hex())
		fmt.Printf("\tValidator: %s\n", info.Pubkey.HexWithPrefix())
		if info.Index == "" {
			fmt.Println("\tIndex: N/A (not seen yet)")
		} else {
			fmt.Printf("\tIndex: %s\n", info.Index)
		}
		fmt.Println()
	}

	// Confirm
	if !(c.Bool("yes") || utils.Confirm("Are you ready to rebuild these keys?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Rebuild the keys
	failedCount := 0
	startIndex := c.Uint64(rebuildStartIndexFlag.Name)
	searchLimit := c.Uint64(rebuildSearchLimitFlag.Name)
	for _, info := range minipoolInfo.Data.Infos {
		start := time.Now()
		fmt.Printf("Rebuilding %s... ", info.Pubkey.HexWithPrefix())
		response, err := cs.Api.Wallet.CreateValidatorKey(info.Pubkey, startIndex, searchLimit)
		if err != nil {
			fmt.Printf("%serror (%s)%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
		} else {
			fmt.Printf("done! (%s)\n", time.Since(start))
			startIndex = response.Data.Index + 1
		}
	}

	// Done
	successCount := len(minipoolInfo.Data.Infos) - failedCount
	fmt.Printf("%d keys have been rebuilt.\n", successCount)
	if successCount == 0 {
		return nil
	}

	// Restart the VC if some keys were rebuilt
	fmt.Println()
	fmt.Println("Your Constellation Validator Client must be restarted in order to load the validator keys and resume attesting wth them.")
	if c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to restart the Constellation Validator Client now?") {
		_, err := hd.Api.Service.RestartContainer(string(csconfig.ContainerID_ConstellationValidator))
		if err != nil {
			fmt.Printf("%sWARNING: Error restarting Constellation Validator Client: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			fmt.Println("Please restart the Constellation Validator Client manually before your validator becomes active in order to load the rebuilt validator keys.")
			fmt.Printf("%sIf you don't restart it, you will miss attestations for validators without loaded keys and lose ETH!%s\n", terminal.ColorYellow, terminal.ColorReset)
		} else {
			fmt.Println("Successfully restarted the Constellation Validator Client. Your rebuilt validator keys are now loaded.")
			return nil
		}
	} else {
		fmt.Println("Please restart the Constellation Validator Client manually to load the rebuilt validator keys.")
		fmt.Printf("%sIf you don't restart it, you will miss attestations for validators without loaded keys and lose ETH, and may be ejected from NodeSet!%s\n", terminal.ColorYellow, terminal.ColorReset)
	}
	return nil
}

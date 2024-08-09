package minipool

import (
	"fmt"

	csapi "github.com/nodeset-org/hyperdrive-constellation/shared/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/rocketpool-go/v2/types"
	"github.com/urfave/cli/v2"

	ncli "github.com/rocket-pool/node-manager-core/cli/utils"
)

var (
	exitMinipoolsFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "minipools",
		Aliases: []string{"m"},
		Usage:   "A comma-separated list of addresses for minipools to exit (or 'all' to exit all eligible minipools)",
	}

	exitVerboseFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"v"},
		Usage:   "Display all minipool details, not just ones that are eligible for exiting",
	}
)

func exitMinipools(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Get eligible minipools
	verbose := c.Bool(exitVerboseFlag.Name)
	detailsResponse, err := cs.Api.Minipool.GetExitDetails(verbose)
	if err != nil {
		return err
	}

	// Check for active minipools
	details := detailsResponse.Data.Details
	if len(details) == 0 {
		fmt.Println("No minipools can be exited.")
		return nil
	}

	// Get selected minipools
	options := make([]ncli.SelectionOption[csapi.MinipoolExitDetails], len(details))
	for i, mp := range details {
		option := &options[i]
		option.Element = &details[i]
		option.ID = fmt.Sprint(mp.Address)

		if mp.MinipoolStatus == types.MinipoolStatus_Staking {
			option.Display = fmt.Sprintf("%s (staking since %s)", mp.Address.Hex(), mp.MinipoolStatusTime.Format(TimeFormat))
		} else {
			option.Display = fmt.Sprintf("%s (dissolved since %s)", mp.Address.Hex(), mp.MinipoolStatusTime.Format(TimeFormat))
		}
	}
	selectedMinipools, err := utils.GetMultiselectIndices(c, exitMinipoolsFlag.Name, options, "Please select a minipool to exit:")
	if err != nil {
		return fmt.Errorf("error determining minipool selection: %w", err)
	}

	// Show a warning message
	fmt.Printf("%sNOTE:\n", terminal.ColorYellow)
	fmt.Println("You are about to exit your minipool. This will tell each one's validator to stop all activities on the Beacon Chain.")
	fmt.Println("Please continue to run your validators until each one you've exited has been processed by the exit queue.\nYou can watch their progress on the https://beaconcha.in explorer.")
	fmt.Println("Your funds will be locked on the Beacon Chain until they've been withdrawn, which will happen automatically (this may take a few days).")
	fmt.Printf("Once your funds have been withdrawn, you can run `rocketpool minipool close` to distribute them to your withdrawal address and close the minipool.\n\n%s", terminal.ColorReset)

	// Prompt for confirmation
	if !(c.Bool("yes") || utils.ConfirmWithIAgree(fmt.Sprintf("Are you sure you want to exit %d minipool(s)? This action cannot be undone!", len(selectedMinipools)))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Exit minipools
	exitInfos := make([]csapi.MinipoolExitInfo, len(selectedMinipools))
	for i, mp := range selectedMinipools {
		exitInfos[i] = csapi.MinipoolExitInfo{
			Address: mp.Address,
			Pubkey:  mp.Pubkey,
			Index:   mp.Index,
		}
	}
	if _, err := cs.Api.Minipool.Exit(exitInfos); err != nil {
		return fmt.Errorf("error while exiting minipools: %w\n", err)
	} else {
		fmt.Println("Successfully exited all selected minipools.")
		fmt.Println("It may take several hours for your minipools' status to be reflected.")
	}

	// Return
	return nil
}
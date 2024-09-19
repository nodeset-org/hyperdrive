package minipool

import (
	"fmt"
	"strconv"

	csapi "github.com/nodeset-org/hyperdrive-constellation/shared/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/rocketpool-go/v2/types"
	"github.com/urfave/cli/v2"

	"github.com/rocket-pool/node-manager-core/beacon"
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

	exitManualModeFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "enable-manual-mode",
		Aliases: []string{"emm"},
		Usage:   "Bypass the exit viability check and enable manual exit mode, using the 'mp' and 'mi' flags. Only use this if you know what you're doing.",
	}

	exitManualPubkeyFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "manual-pubkey",
		Aliases: []string{"mp"},
		Usage:   "Manually specify the pubkey of the locally-stored validator to exit. Requires the 'emm' and 'mi' flags.",
	}

	exitManualIndexFlag *cli.IntFlag = &cli.IntFlag{
		Name:    "manual-index",
		Aliases: []string{"mi"},
		Usage:   "Manually specify the index of the locally-stored validator to exit. Requires the 'emm' and 'mp' flags.",
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

	// Get the local force flag
	enableManualMode := c.Bool(exitManualModeFlag.Name)

	// Normal exit
	var selectedMinipools []*csapi.MinipoolExitDetails
	if !enableManualMode {
		// Get eligible minipools
		verbose := c.Bool(exitVerboseFlag.Name)
		detailsResponse, err := cs.Api.Minipool.GetExitDetails(verbose)
		if err != nil {
			return err
		}

		// Check for active minipools
		var eligibleMpDetails []csapi.MinipoolExitDetails
		if !verbose {
			eligibleMpDetails = detailsResponse.Data.Details
		} else {
			eligibleMpDetails = []csapi.MinipoolExitDetails{}
			alreadyFinalizedDetails := []csapi.MinipoolExitDetails{}
			invalidStatusDetails := []csapi.MinipoolExitDetails{}
			invalidValidatorDetails := []csapi.MinipoolExitDetails{}
			validatorTooYoungDetails := []csapi.MinipoolExitDetails{}
			validatorNotSeenYetDetails := []csapi.MinipoolExitDetails{}
			for _, mp := range detailsResponse.Data.Details {
				if mp.CanExit {
					eligibleMpDetails = append(eligibleMpDetails, mp)
					continue
				}
				if mp.AlreadyFinalized {
					alreadyFinalizedDetails = append(alreadyFinalizedDetails, mp)
					continue
				}
				if mp.InvalidMinipoolStatus {
					invalidStatusDetails = append(invalidStatusDetails, mp)
					continue
				}
				if mp.InvalidValidatorStatus {
					invalidValidatorDetails = append(invalidValidatorDetails, mp)
					continue
				}
				if mp.ValidatorTooYoung {
					validatorTooYoungDetails = append(validatorTooYoungDetails, mp)
					continue
				}
				if mp.ValidatorNotSeenYet {
					validatorNotSeenYetDetails = append(validatorNotSeenYetDetails, mp)
					continue
				}
			}

			// Print details
			if len(alreadyFinalizedDetails) > 0 {
				fmt.Println("The following minipools are already finalized:")
				for _, mp := range alreadyFinalizedDetails {
					fmt.Printf("  %s\n", mp.Address.Hex())
				}
				fmt.Println()
			}
			if len(invalidStatusDetails) > 0 {
				fmt.Println("The following minipools are not in an exitable state:")
				for _, mp := range invalidStatusDetails {
					fmt.Printf("  %s\n", mp.Address.Hex())
				}
				fmt.Println()
			}
			if len(invalidValidatorDetails) > 0 {
				fmt.Println("The following minipools have validators that are not in an exitable state:")
				for _, mp := range invalidValidatorDetails {
					fmt.Printf("  %s\n", mp.Address.Hex())
				}
				fmt.Println()
			}
			if len(validatorTooYoungDetails) > 0 {
				fmt.Println("The following minipools have validators that are too young to exit:")
				for _, mp := range validatorTooYoungDetails {
					fmt.Printf("  %s\n", mp.Address.Hex())
				}
				fmt.Println()
			}
			if len(validatorNotSeenYetDetails) > 0 {
				fmt.Println("The following minipools have validators that have not been seen on the Beacon Chain yet:")
				for _, mp := range validatorNotSeenYetDetails {
					fmt.Printf("  %s\n", mp.Address.Hex())
				}
				fmt.Println()
			}
		}
		if len(eligibleMpDetails) == 0 {
			fmt.Println("No minipools can be exited.")
			return nil
		}

		// Get selected minipools
		options := make([]ncli.SelectionOption[csapi.MinipoolExitDetails], len(eligibleMpDetails))
		for i, mp := range eligibleMpDetails {
			option := &options[i]
			option.Element = &eligibleMpDetails[i]
			option.ID = fmt.Sprint(mp.Address)

			if mp.MinipoolStatus == types.MinipoolStatus_Staking {
				option.Display = fmt.Sprintf("%s (staking since %s)", mp.Address.Hex(), mp.MinipoolStatusTime.Format(TimeFormat))
			} else {
				option.Display = fmt.Sprintf("%s (dissolved since %s)", mp.Address.Hex(), mp.MinipoolStatusTime.Format(TimeFormat))
			}
		}
		selectedMinipools, err = utils.GetMultiselectIndices(c, exitMinipoolsFlag.Name, options, "Please select a minipool to exit:")
		if err != nil {
			return fmt.Errorf("error determining minipool selection: %w", err)
		}
	} else {
		// Manual exit
		pubkeyString := c.String(exitManualPubkeyFlag.Name)
		index := c.Int(exitManualIndexFlag.Name)
		if pubkeyString == "" {
			return fmt.Errorf("pubkey is required when using the manual exit flag")
		}
		if index == 0 {
			return fmt.Errorf("index is required when using the manual exit flag")
		}

		// Make a fake exit details with the manual info
		pubkey, err := beacon.HexToValidatorPubkey(pubkeyString)
		if err != nil {
			return fmt.Errorf("error parsing pubkey: %w", err)
		}
		selectedMinipools = []*csapi.MinipoolExitDetails{
			{
				Pubkey: pubkey,
				Index:  strconv.FormatInt(int64(index), 10),
			},
		}
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

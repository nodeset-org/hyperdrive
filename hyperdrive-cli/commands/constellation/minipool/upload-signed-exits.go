package minipool

import (
	"fmt"

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
	uploadSignedExitsMinipoolsFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "minipools",
		Aliases: []string{"m"},
		Usage:   "A comma-separated list of addresses for minipools to upload signed exit messages for (or 'all' to upload exits for all eligible minipools)",
	}
)

func uploadSignedExits(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Normal exit
	var selectedMinipools []*csapi.MinipoolExitDetails
	detailsResponse, err := cs.Api.Minipool.GetExitDetails(true)
	if err != nil {
		return err
	}

	// Do some custom filtering
	filteredMinipools := []csapi.MinipoolExitDetails{}
	for _, mp := range detailsResponse.Data.Details {
		// Skip ineligible minipools
		if mp.InvalidMinipoolStatus || mp.AlreadyFinalized || mp.ValidatorNotSeenYet {
			continue
		}

		// Allow ones that are too young or still pending
		switch mp.ValidatorStatus {
		case beacon.ValidatorState_PendingInitialized,
			beacon.ValidatorState_PendingQueued,
			beacon.ValidatorState_ActiveOngoing:
			filteredMinipools = append(filteredMinipools, mp)
		default:
			continue
		}
	}

	// Get the list of validators with exits already uploaded
	outstandingMinipools := []csapi.MinipoolExitDetails{}
	validatorsResponse, err := hd.Api.NodeSet_Constellation.GetValidators()
	if err != nil {
		return err
	}
	for _, mp := range filteredMinipools {
		hasExit := false
		for _, validatorInfo := range validatorsResponse.Data.Validators {
			if mp.Pubkey != validatorInfo.Pubkey {
				continue
			}
			hasExit = validatorInfo.ExitMessageUploaded
			break
		}
		if !hasExit {
			outstandingMinipools = append(outstandingMinipools, mp)
		}
	}

	// Check for active minipools
	details := outstandingMinipools
	if len(details) == 0 {
		fmt.Println("No minipools can have signed exits uploaded to NodeSet.")
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
	selectedMinipools, err = utils.GetMultiselectIndices(c, exitMinipoolsFlag.Name, options, "Please select a minipool to upload a signed exit message to NodeSet for:")
	if err != nil {
		return fmt.Errorf("error determining minipool selection: %w", err)
	}

	// Show a warning message
	fmt.Printf("%sNOTE:\n", terminal.ColorYellow)
	fmt.Println("Your node will sign a voluntary exit message for each of the selected minipools and upload them to the NodeSet service.")
	fmt.Printf("NodeSet will not submit these message under normal circumstances; if you want to voluntarily exit the validator yourself while it's still in good standing, you must do so manually via the Hyperdrive exit command.\n\n%s", terminal.ColorReset)

	// Prompt for confirmation
	if !(c.Bool("yes") || utils.ConfirmWithIAgree(fmt.Sprintf("Are you ready to send signed exits for %d minipool(s)?", len(selectedMinipools)))) {
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
	if _, err := cs.Api.Minipool.UploadSignedExits(exitInfos); err != nil {
		return fmt.Errorf("error while uploading signed exits: %w\n", err)
	} else {
		fmt.Println("Successfully uploaded signed exit messages for all selected minipools to NodeSet.")
	}

	// Return
	return nil
}

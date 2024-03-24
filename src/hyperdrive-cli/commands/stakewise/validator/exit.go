package validator

import (
	"fmt"
	"time"

	"github.com/nodeset-org/hyperdrive/shared/types"

	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

var (
	pubkeysFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "pubkeys",
		Aliases: []string{"p"},
		Usage:   "Comma-separated list of pubkeys (including 0x prefix) to get the exit message for",
	}
	epochFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "epoch",
		Aliases: []string{"e"},
		Usage:   "(Optional) the epoch to use when creating the signed exit messages. If not specified, the current chain head will be used.",
	}
	noBroadcastFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "no-broadcast",
		Aliases: []string{"n"},
		Usage:   "(Optional) pass this flag to skip broadcasting the exit message(s) and print them instead",
	}
)

func exit(c *cli.Context) error {
	// Get the client
	sw := client.NewStakewiseClientFromCtx(c)

	// Get all active validators
	activeValidatorResponse, err := sw.Api.Status.GetValidatorStatuses()
	if err != nil {
		return fmt.Errorf("error while getting active validators: %w", err)
	}
	var activeValidators []string // []beacon.ValidatorPubkey.HexWithPrefix()
	for pubKey, status := range activeValidatorResponse.Data.BeaconStatus {
		if status == types.ValidatorState_ActiveOngoing {
			activeValidators = append(activeValidators, pubKey)
		}
	}
	if len(activeValidators) == 0 {
		fmt.Println("None of your validators are active, so they cannot be exited.")
		return nil
	}

	// Get selected validators
	options := make([]utils.SelectionOption[beacon.ValidatorPubkey], len(activeValidators))
	for i, pubkey := range activeValidators {
		pubKey, err := beacon.HexToValidatorPubkey(activeValidators[i])
		if err != nil {
			return fmt.Errorf("error while converting validator pubkey: %w", err)
		}
		option := &options[i]
		option.Element = &pubKey
		option.ID = activeValidators[i]
		option.Display = fmt.Sprintf("%s (active since %s)", pubkey, time.Unix(0, 0)) // Placeholder, fill in with status details
	}
	selectedValidators, err := utils.GetMultiselectIndices(c, pubkeysFlag.Name, options, "Please select a validator to get the signed exit for:")
	if err != nil {
		return fmt.Errorf("error determining validator selection: %w", err)
	}

	// Get the epoch if set
	var epochPtr *uint64
	if c.IsSet(epochFlag.Name) {
		epoch := c.Uint64(epochFlag.Name)
		epochPtr = &epoch
	}

	// Get the no broadcast flag
	noBroadcastBool := false
	if c.IsSet(noBroadcastFlag.Name) {
		noBroadcastBool = c.Bool(noBroadcastFlag.Name)
	}

	// Get the pubkeys
	pubkeys := make([]beacon.ValidatorPubkey, len(selectedValidators))
	for i, validator := range selectedValidators {
		pubkeys[i] = *validator
	}

	if !noBroadcastBool {
		// Show a warning message
		fmt.Printf("%sNOTE:\n", terminal.ColorYellow)
		fmt.Println("You are about to exit your validator(s). This will tell each one to stop all activities on the Beacon Chain.")
		fmt.Println("Please continue to run them until each one you've exited has been processed by the exit queue. It will no longer earn staking rewards after this point.")
		fmt.Printf("Your funds will be locked on the Beacon Chain until they've been withdrawn, which will happen automatically (typically after a few days).%s\n", terminal.ColorReset)

		// Prompt for confirmation
		if !(c.Bool(utils.YesFlag.Name) || utils.ConfirmWithIAgree(fmt.Sprintf("Are you sure you want to exit %d validator(s)? This action cannot be undone!", len(selectedValidators)))) {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Get signed exit messages
	response, err := sw.Api.Validator.Exit(pubkeys, epochPtr, noBroadcastBool)
	if err != nil {
		return fmt.Errorf("error while getting validator exit messages: %w", err)
	}

	// Log successand return if not broadcasting
	if !noBroadcastBool {
		fmt.Println("Successfully exited the selected validator(s). It will take some time before their status is reflected on the Beacon Chain.")
		return nil
	}

	// Print them all
	fmt.Printf("Exit epoch: %d\n", response.Data.Epoch)
	fmt.Println()
	for _, info := range response.Data.ExitInfos {
		fmt.Printf("Validator %d (%s):\n", info.Index, info.Pubkey.HexWithPrefix())
		fmt.Printf("\tSignature: %s\n", info.Signature.HexWithPrefix())
		fmt.Println()
	}

	// Return
	return nil
}

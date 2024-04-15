package status

import (
	"fmt"

	swtypes "github.com/nodeset-org/hyperdrive-stakewise/shared/types"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/urfave/cli/v2"
)

func getNodeStatus(c *cli.Context) error {
	sw := client.NewStakewiseClientFromCtx(c)
	response, err := sw.Api.Status.GetValidatorStatuses()
	if err != nil {
		fmt.Printf("error fetching validator statuses: %v\n", err)
		return err
	}

	if len(response.Data.States) == 0 {
		fmt.Println("You do not have any validators.")
		return nil
	}

	for _, state := range response.Data.States {
		fmt.Printf("%s:\n", state.Pubkey.HexWithPrefix())

		// Print Beacon status
		if state.Index == "" {
			fmt.Println("\tBeacon State: Not seen by Beacon Chain yet")
		} else {
			fmt.Printf("\tBeacon Index: %s\n", state.Index)
			fmt.Printf("\tBeacon State: %s\n", getBeaconStatusLabel(state.BeaconStatus))
		}

		// Print NodeSet status
		fmt.Printf("\tNodeSet State: %s\n", getNodeSetStateLabel(state.NodesetStatus))
		fmt.Println()
	}

	return nil
}

func getBeaconStatusLabel(state beacon.ValidatorState) string {
	switch state {
	case beacon.ValidatorState_ActiveExiting:
		return "Active (Exiting in Progress)"
	case beacon.ValidatorState_ActiveOngoing:
		return "Active"
	case beacon.ValidatorState_ActiveSlashed:
		return "Slashed (Exit in Progress)"
	case beacon.ValidatorState_ExitedSlashed:
		return "Slashed (Exited)"
	case beacon.ValidatorState_ExitedUnslashed:
		return "Exited (Withdrawal Pending)"
	case beacon.ValidatorState_PendingInitialized:
		return "Seen on Beacon, Waiting for More Deposits"
	case beacon.ValidatorState_PendingQueued:
		return "In Beacon Activation Queue"
	case beacon.ValidatorState_WithdrawalDone:
		return "Exited and Withdrawn"
	case beacon.ValidatorState_WithdrawalPossible:
		return "Exited (Waiting for Wihdrawal)"
	default:
		return fmt.Sprintf("<Unknown Beacon Status: %s>", state)
	}
}

func getNodeSetStateLabel(state swtypes.NodesetStatus) string {
	switch state {
	case swtypes.NodesetStatus_Generated:
		return "Generated (Not Yet Uploaded)"
	case swtypes.NodesetStatus_RegisteredToStakewise:
		return "Registered with Stakewise"
	case swtypes.NodesetStatus_UploadedStakewise:
		return "Uploaded to Stakewise"
	case swtypes.NodesetStatus_UploadedToNodeset:
		return "Uploaded to NodeSet"
	default:
		return fmt.Sprintf("<Unknown NodeSet Status: %s>", state)
	}
}

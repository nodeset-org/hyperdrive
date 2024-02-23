package status

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func getNodeStatus(c *cli.Context) error {
	sw := client.NewStakewiseClientFromCtx(c)
	response, err := sw.Api.Status.GetValidatorStatuses()
	if err != nil {
		fmt.Printf("error fetching active validators: %v\n", err)
		return err
	}

	if len(response.Data.WithdrawalDone) != 0 {
		fmt.Printf("Withdrawal Done: \n")

		for _, validator := range response.Data.WithdrawalDone {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.WithdrawalPossible) != 0 {
		fmt.Printf("Withdrawal Possible: \n")

		for _, validator := range response.Data.WithdrawalPossible {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.ExitedSlashed) != 0 {
		fmt.Printf("Exited Slashed: \n")

		for _, validator := range response.Data.ExitedSlashed {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.ExitedUnslashed) != 0 {
		fmt.Printf("Exited Unslashed: \n")

		for _, validator := range response.Data.ExitedUnslashed {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.ActiveSlashed) != 0 {
		fmt.Printf("Active Slashed: \n")

		for _, validator := range response.Data.ActiveSlashed {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.ActiveExited) != 0 {
		fmt.Printf("Active Exited: \n")

		for _, validator := range response.Data.ActiveExited {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.ActiveOngoing) != 0 {
		fmt.Printf("Active Ongoing: \n")

		for _, validator := range response.Data.ActiveOngoing {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.PendingQueued) != 0 {
		fmt.Printf("Pending Queued: \n")

		for _, validator := range response.Data.PendingQueued {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.PendingInitialized) != 0 {
		fmt.Printf("Pending Initialized: \n")

		for _, validator := range response.Data.PendingInitialized {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.RegisteredToStakewise) != 0 {
		fmt.Printf("Registered To Stakewise: \n")

		for _, validator := range response.Data.RegisteredToStakewise {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.UploadToStakewise) != 0 {
		fmt.Printf("Uploaded To Stakewise: \n")

		for _, validator := range response.Data.UploadToStakewise {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.UploadedToNodeset) != 0 {
		fmt.Printf("Uploaded To Nodeset: \n")

		for _, validator := range response.Data.UploadedToNodeset {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.Generated) != 0 {
		fmt.Printf("Generated Validator Pubkeys: \n")

		for _, validator := range response.Data.Generated {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	return nil
}

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

	if len(response.Data.Active) != 0 {
		fmt.Printf("Active Validator Pubkeys: \n")

		for _, validator := range response.Data.Active {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.Exited) != 0 {
		fmt.Printf("Exited Validator Pubkeys: \n")

		for _, validator := range response.Data.Exited {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.Exiting) != 0 {
		fmt.Printf("Exiting Validator Pubkeys: \n")

		for _, validator := range response.Data.Exiting {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}
	fmt.Println("")

	if len(response.Data.Deposited) != 0 {
		fmt.Printf("Deposited Validator Pubkeys: \n")

		for _, validator := range response.Data.Deposited {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.Depositing) != 0 {
		fmt.Printf("Depositing Validator Pubkeys: \n")

		for _, validator := range response.Data.Depositing {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.RegisteredToStakewise) != 0 {
		fmt.Printf("Registered to Stakewise Validator Pubkeys: \n")

		for _, validator := range response.Data.RegisteredToStakewise {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.UploadedToNodeset) != 0 {
		fmt.Printf("Uploaded to Nodeset Validator Pubkeys: \n")

		for _, validator := range response.Data.UploadedToNodeset {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.UploadToStakewise) != 0 {
		fmt.Printf("Upload to Stakewise Validator Pubkeys: \n")

		for _, validator := range response.Data.UploadToStakewise {
			fmt.Printf("%v\n", validator.HexWithPrefix())
		}
		fmt.Println("")
	}

	if len(response.Data.WaitingDepositConfirmation) != 0 {
		fmt.Printf("Waiting Deposit Confirmation Validator Pubkeys: \n")

		for _, validator := range response.Data.WaitingDepositConfirmation {
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

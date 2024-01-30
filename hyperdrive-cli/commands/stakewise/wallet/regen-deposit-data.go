package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func regenerateDepositData(c *cli.Context) error {
	// Get Stakewise client
	sw := client.NewStakewiseClientFromCtx(c)

	fmt.Println("Regenerating complete deposit data, please wait...")
	regenResponse, err := sw.Api.Wallet.RegenerateDepositData()
	if err != nil {
		fmt.Println("%sThere was an error regenerating your deposit data. Please run it manually with `hyperdrive stakewise wallet regen-deposit-data` to try again.%s", terminal.ColorYellow, terminal.ColorReset)
		return fmt.Errorf("error regenerating deposit data: %w", err)
	}

	// Print the total
	fmt.Printf("Total keys loaded: %s%d%s\n", terminal.ColorGreen, len(regenResponse.Data.Pubkeys), terminal.ColorReset)
	fmt.Println()

	if c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to restart the Stakewise Operator service so it loads the new keys and deposit data?") {
		fmt.Println("NYI")
	} else {
		fmt.Println("Please restart the container at your convenience.")
	}
	fmt.Println()

	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to upload the deposit data with the new keys to the NodeSet server, so they can be used for new validator assignments?")) {
		fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise service upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
		return nil
	}

	// TODO

	fmt.Println("<NYI>")

	return nil
}

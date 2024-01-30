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

	fmt.Printf("Regenerating aggregated deposit data... ")

	err := regenDepositData(c, sw)
	if err != nil {
		fmt.Println("error")
		return err
	}

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

	fmt.Printf("Uploading deposit data to the NodeSet server... ")

	_, err = sw.Api.Nodeset.UploadDepositData()
	if err != nil {
		fmt.Println("error")
		fmt.Printf("%sWARNING: error uploading deposit data to nodeset: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
		fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
		fmt.Println()
	} else {
		fmt.Println("done!")
	}

	return nil
}

package nodeset

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func uploadDepositData(c *cli.Context) error {
	// Get the client
	sw := client.NewStakewiseClientFromCtx(c)

	fmt.Printf("Uploading deposit data to the NodeSet server... ")
	_, err := sw.Api.Nodeset.UploadDepositData()
	if err != nil {
		fmt.Println("error")
		fmt.Printf("%sWARNING: error uploading deposit data to nodeset: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
		fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
		fmt.Println()
	}

	fmt.Println("done!")
	fmt.Println("Your Stakewise Operator container should reflect the new deposit data shortly.")
	return nil
}

package swcmdutils

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	swconfig "github.com/nodeset-org/hyperdrive/shared/config/modules/stakewise"
)

func RegenDepositData(hd *client.HyperdriveClient, sw *client.StakewiseClient, noRestart bool) error {
	// Regen the aggregated deposit data
	fmt.Println("Regenerating complete deposit data, please wait...")
	regenResponse, err := sw.Api.Wallet.RegenerateDepositData()
	if err != nil {
		fmt.Printf("%sThere was an error regenerating your deposit data. Please run it manually with `hyperdrive stakewise wallet regen-deposit-data` to try again.%s", terminal.ColorYellow, terminal.ColorReset)
		return fmt.Errorf("error regenerating deposit data: %w", err)
	}

	// Print the total
	fmt.Printf("Total keys loaded: %s%d%s\n", terminal.ColorGreen, len(regenResponse.Data.Pubkeys), terminal.ColorReset)
	fmt.Println()

	// Restart the Stakewise Operator
	if noRestart {
		fmt.Printf("%sYou have automatic restarting turned off.\nPlease restart your Stakewise Operator service at your earliest convenience so it can deposit your validators when it's your turn in the queue.%s\n", terminal.ColorYellow, terminal.ColorReset)
	} else {
		fmt.Print("Restarting Stakewise Operator... ")
		_, err = hd.Api.Service.RestartContainer(swconfig.OperatorContainerSuffix)
		if err != nil {
			fmt.Println("error")
			fmt.Printf("%sWARNING: error restarting stakewise operator: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			fmt.Println("Please restart it in order to assign deposits to your new keys.")
			fmt.Println()
		} else {
			fmt.Println("done!")
		}
	}
	fmt.Println()
	return nil
}

// Upload deposit data to the server
func UploadDepositData(sw *client.StakewiseClient) error {
	fmt.Printf("Uploading deposit data to the NodeSet server... ")
	response, err := sw.Api.Nodeset.UploadDepositData()
	if err != nil {
		fmt.Println("error")
		fmt.Printf("%sWARNING: error uploading deposit data to nodeset: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
		fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
		fmt.Println()
	} else {
		fmt.Println("done!")
		fmt.Println("Server returned: %s\n", string(response.Data.ServerResponse))
	}
	return nil
}

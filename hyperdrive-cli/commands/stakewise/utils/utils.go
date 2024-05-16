package swcmdutils

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
)

func printUploadError(err error) {
	fmt.Println("Error")
	fmt.Printf("%sWARNING: Error uploading deposit data to NodeSet: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
	fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
	fmt.Println()
}

// Upload deposit data to the server
func UploadDepositData(sw *client.StakewiseClient) (bool, error) {
	// Initial attempt to upload all deposit data
	fmt.Println("Uploading deposit data to the NodeSet server...")
	response, err := sw.Api.Nodeset.UploadDepositData()
	if err != nil {
		printUploadError(err)
		return false, nil
	}

	data := response.Data
	sw.Logger.Debug("Server response", "data", data.ServerResponse)
	fmt.Println()
	newKeyCount := len(data.UnregisteredPubkeys)

	if data.SufficientBalance {
		if newKeyCount == 0 {
			fmt.Println("All of your validator keys are registered.")
			return false, nil
		}
		fmt.Printf("Registered %s%d%s new validator keys:\n", terminal.ColorGreen, newKeyCount, terminal.ColorReset)
		for _, key := range data.UnregisteredPubkeys {
			fmt.Println(key.HexWithPrefix())
		}
		fmt.Println()
		fmt.Printf("Total keys registered: %s%d%s\n", terminal.ColorGreen, data.TotalCount, terminal.ColorReset)
	} else {
		fmt.Println("Not all keys were uploaded due to insufficient balance.")
		fmt.Printf("ETH required per key: %s%f%s\n", terminal.ColorGreen, data.EthPerKey, terminal.ColorReset)
		fmt.Printf("Current Balance: %s%f%s\n", terminal.ColorGreen, data.Balance, terminal.ColorReset)
		fmt.Printf("Additional ETH required for remaining keys: %s%f%s\n", terminal.ColorGreen, data.RemainingEthRequired, terminal.ColorReset)

		fmt.Printf("\nUploaded %d out of %d keys:\n", newKeyCount, data.TotalCount)
		for _, key := range data.UnregisteredPubkeys {
			fmt.Println(key.HexWithPrefix())
		}
	}

	return true, nil
}

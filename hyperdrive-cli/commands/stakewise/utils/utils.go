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
	newKeyCount := len(data.UnregisteredPubkeys)

	if newKeyCount == 0 && data.SufficientBalance {
		fmt.Println("All of your validator keys were already registered.")
		return false, nil
	}

	if !data.SufficientBalance {
		fmt.Println("Not all keys were uploaded due to insufficient balance.")
		fmt.Printf("Uploaded %d out of %d keys.", newKeyCount, data.TotalCount)
	}

	data = response.Data
	sw.Logger.Debug("Server response", "data", data.ServerResponse)
	fmt.Println()
	fmt.Printf("Registered %s%d%s new validator keys:\n", terminal.ColorGreen, len(data.UnregisteredPubkeys), terminal.ColorReset)
	for _, key := range data.UnregisteredPubkeys {
		fmt.Println(key.HexWithPrefix())
	}
	fmt.Println()

	fmt.Printf("Total keys registered: %s%d%s\n", terminal.ColorGreen, data.TotalCount, terminal.ColorReset)
	return true, nil
}

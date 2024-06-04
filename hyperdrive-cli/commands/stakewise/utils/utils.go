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
	if response.Data.UnregisteredNode {
		fmt.Println("Your node is not registered with NodeSet yet. Please register your node first.")
		return false, nil
	}

	data := response.Data

	newKeyCount := len(data.NewPubkeys)
	remainingKeyCount := len(data.RemainingPubkeys)

	if data.SufficientBalance {
		if newKeyCount == 0 {
			fmt.Printf("All of your validator keys are already registered (%s%d%s in total).\n", terminal.ColorGreen, data.TotalCount, terminal.ColorReset)
			fmt.Printf("%s%d%s are pending activation.\n", terminal.ColorGreen, data.PendingCount, terminal.ColorReset)
			fmt.Printf("%s%d%s have been activated already.\n", terminal.ColorGreen, data.ActiveCount, terminal.ColorReset)
			return false, nil
		}
		fmt.Printf("Registered %s%d%s new validator keys:\n", terminal.ColorGreen, newKeyCount, terminal.ColorReset)
		for _, key := range data.NewPubkeys {
			fmt.Println(key.HexWithPrefix())
		}
	} else {
		fmt.Println("Not all keys were uploaded due to insufficient balance.")
		fmt.Printf("Current wallet balance: %s%f%s\n", terminal.ColorGreen, data.Balance, terminal.ColorReset)
		fmt.Printf("Remaining unregistered keys: %s%d%s\n", terminal.ColorGreen, remainingKeyCount, terminal.ColorReset)
		fmt.Printf("You need %s%f%s more ETH to register your remaining keys.\n", terminal.ColorGreen, data.RemainingEthRequired, terminal.ColorReset)

		totalUnregisteredKeyCount := len(data.NewPubkeys) + len(data.RemainingPubkeys)
		if newKeyCount == 0 {
			fmt.Printf("\nUploaded 0 out of %d new keys.\n", totalUnregisteredKeyCount)
		} else {
			fmt.Printf("\nUploaded %d out of %d new keys:\n", newKeyCount, totalUnregisteredKeyCount)
			for _, key := range data.NewPubkeys {
				fmt.Println(key.HexWithPrefix())
			}
		}
	}
	data.PendingCount += uint64(newKeyCount)
	fmt.Println()
	fmt.Printf("Total keys: %s%d%s\n", terminal.ColorGreen, data.TotalCount, terminal.ColorReset)
	fmt.Printf("%s%d%s are registered and pending activation.\n", terminal.ColorGreen, data.PendingCount, terminal.ColorReset)
	fmt.Printf("%s%d%s are registered and have been activated already.\n", terminal.ColorGreen, data.ActiveCount, terminal.ColorReset)

	return true, nil
}

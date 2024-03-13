package swcmdutils

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/eth"
)

func printUploadError(err error) {
	fmt.Println("Error")
	fmt.Printf("%sWARNING: Error uploading deposit data to NodeSet: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
	fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
	fmt.Println()
}

// Upload deposit data to the server
func UploadDepositData(sw *client.StakewiseClient) error {
	fmt.Println("Uploading deposit data to the NodeSet server...")
	response, err := sw.Api.Nodeset.UploadDepositData(false)
	if err != nil {
		printUploadError(err)
		return nil
	}

	data := response.Data
	if !data.SufficientBalance {
		// Prompt the user to upload anyway
		fmt.Printf("You're attempting to upload %d keys, but you only have %.6f ETH in your account. We recommend you have at least %.6f ETH", data.UnregisteredKeyCount, eth.WeiToEth(data.Balance), eth.WeiToEth(data.RequiredBalance))
		if !utils.Confirm("Do you want to upload these keys anyway? You may not be able to register them if your wallet doesn't have sufficient ETH in it!") {
			fmt.Println("Cancelled.")
			return nil
		}

		response, err = sw.Api.Nodeset.UploadDepositData(true)
		if err != nil {
			printUploadError(err)
			return nil
		}
	}

	data = response.Data
	fmt.Println("done!")
	if len(data.NewPubkeys) == 0 {
		fmt.Println("All of your validator keys were already registered.")
	} else {
		fmt.Printf("Server returned: %s\n", string(data.ServerResponse))
		fmt.Println()
		fmt.Printf("Registered %s%d%s new validator keys:\n", terminal.ColorGreen, len(data.NewPubkeys), terminal.ColorReset)
		for _, key := range data.NewPubkeys {
			fmt.Println(key.HexWithPrefix())
		}
		fmt.Println()
	}

	fmt.Printf("Total keys registered: %s%d%s\n", terminal.ColorGreen, data.TotalCount, terminal.ColorReset)
	return nil
}

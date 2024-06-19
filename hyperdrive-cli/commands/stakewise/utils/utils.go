package swcmdutils

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func printUploadError(err error) {
	fmt.Println("Error")
	fmt.Printf("%sWARNING: Error uploading deposit data to NodeSet: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
	fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
	fmt.Println()
}

// Upload deposit data to the server
func UploadDepositData(c *cli.Context, sw *client.StakewiseClient) (bool, error) {
	// Warn user prior to uploading deposit data
	fmt.Println("NOTE: There is currently no way to remove a validator's deposit data from the NodeSet service once you've uploaded it. The key will be eligible for activation at any time, so this node must remain online at all times to handle activation and validation duties.")
	fmt.Printf("%sIf you turn the node off, you may be removed from NodeSet for negligence of duty!%s\n", terminal.ColorYellow, terminal.ColorReset)
	fmt.Println()

	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Do you want to continue uploading your deposit data?")) {
		fmt.Println("Cancelled.")
		return false, nil
	}

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

	if data.InvalidWithdrawalCredentials {
		fmt.Printf("%sWARNING: Your deposit data contained withdrawal credentials that do not correspond to a valid StakeWise vault. Please contact the Hyperdrive developers and report this issue.%s\n", terminal.ColorYellow, terminal.ColorReset)
		return false, nil
	}
	if data.NotAuthorizedForMainnet {
		fmt.Printf("%sWARNING: Your deposit data was rejected because you are not currently authorized to access the Mainnet vault. You will need to run on the Holesky testnet first before being given access to Mainnet.%s\n", terminal.ColorYellow, terminal.ColorReset)
		return false, nil
	}

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
		fmt.Printf("%sWarning: not all keys were uploaded due to insufficient balance.%s\n", terminal.ColorYellow, terminal.ColorReset)
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

	return newKeyCount != 0, nil
}

package swcmdutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
)

// Upload deposit data to the server
func UploadDepositData(sw *client.StakewiseClient) error {
	fmt.Printf("Uploading deposit data to the NodeSet server... ")
	response, err := sw.Api.Nodeset.UploadDepositData(false)
	if err != nil {
		if strings.Contains(err.Error(), "balance_check_failed") {
			// Prompt the user for decision on balance check error
			fmt.Println("balance check failed blah blah")
			fmt.Println("Do you want to continue regardless? (yes/no)")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if strings.ToLower(input) != "yes" {
				fmt.Println("Operation aborted by the user.")
				return err
			} else {
				response, _ = sw.Api.Nodeset.UploadDepositData(true)
			}

		} else {
			fmt.Println("error")
			fmt.Printf("%sWARNING: error uploading deposit data to nodeset: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
			fmt.Println()
			return err
		}
	}

	data := response.Data
	fmt.Println("done!")
	if len(data.NewPubkeys) == 0 {
		fmt.Println("All of your validator keys were already registered.")
	} else {
		fmt.Printf("Server returned: %s\n", string(data.ServerResponse))
		fmt.Println()
		fmt.Printf("Registered %s%d%s new validator keys:\n", terminal.ColorGreen, len(data.NewPubkeys), terminal.ColorReset)
		for _, key := range response.Data.NewPubkeys {
			fmt.Println(key.HexWithPrefix())
		}
		fmt.Println()
	}

	fmt.Printf("Total keys registered: %s%d%s\n", terminal.ColorGreen, data.TotalCount, terminal.ColorReset)

	return nil
}

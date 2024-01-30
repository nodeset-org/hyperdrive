package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func regenDepositData(c *cli.Context, sw *client.StakewiseClient) error {
	fmt.Println("Regenerating complete deposit data, please wait...")
	regenResponse, err := sw.Api.Wallet.RegenerateDepositData()
	if err != nil {
		fmt.Println("%sThere was an error regenerating your deposit data. Please run it manually with `hyperdrive stakewise wallet regen-deposit-data` to try again.%s", terminal.ColorYellow, terminal.ColorReset)
		return fmt.Errorf("error regenerating deposit data: %w", err)
	}

	// Print the total
	fmt.Printf("Total keys loaded: %s%d%s\n", terminal.ColorGreen, len(regenResponse.Data.Pubkeys), terminal.ColorReset)
	fmt.Println()
	return nil
}

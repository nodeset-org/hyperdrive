package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	"github.com/urfave/cli/v2"
)

var (
	generateKeysCountFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "count",
		Aliases: []string{"c"},
		Usage:   "The number of keys to generate",
	}
)

func generateKeys(c *cli.Context) error {
	// Get Stakewise client
	hd := client.NewStakewiseClientFromCtx(c)

	// Get the count
	var err error
	count := c.Uint64(generateKeysCountFlag.Name)
	if count == 0 {
		countString := utils.Prompt("How many keys would you like to generate?", "^\\d+$", "Invalid count, try again")
		count, err = input.ValidateUint("count", countString)
		if err != nil {
			return fmt.Errorf("invalid count [%s]: %w", countString, err)
		}
	}

	// Generate the new keys
	response, err := hd.Api.Wallet.GenerateKeys(count)
	if err != nil {
		return fmt.Errorf("error generating keys: %w", err)
	}

	// Print them
	fmt.Println("New keys:")
	for _, pubkey := range response.Data.Pubkeys {
		fmt.Println(pubkey.Hex())
	}
	fmt.Println()

	fmt.Printf("You now have %d keys ready for validation.\n", response.Data.TotalCount)
	fmt.Println()

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

	// TODO

	fmt.Println("<NYI>")

	return nil
}

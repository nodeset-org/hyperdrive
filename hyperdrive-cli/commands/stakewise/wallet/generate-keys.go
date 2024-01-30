package wallet

import (
	"fmt"
	"time"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	swconfig "github.com/nodeset-org/hyperdrive/shared/config/modules/stakewise"
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
	// Get the client
	hd := client.NewClientFromCtx(c)
	sw := client.NewStakewiseClientFromCtx(c)

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

	fmt.Println("Note: key generation is an expensive process, this may take a long time! Progress will be printed as each key is generated.")
	fmt.Println()

	// Generate the new keys
	startTime := time.Now()
	latestTime := startTime
	for i := uint64(0); i < count; i++ {
		response, err := sw.Api.Wallet.GenerateKeys(1, false)
		if err != nil {
			return fmt.Errorf("error generating keys: %w", err)
		}
		if len(response.Data.Pubkeys) == 0 {
			return fmt.Errorf("server did not return any pubkeys")
		}

		elapsed := time.Since(latestTime)
		latestTime = time.Now()
		pubkey := response.Data.Pubkeys[0]
		fmt.Printf("Generated %s (%d/%d) in %s\n", pubkey.Hex(), (i + 1), count, elapsed)
	}
	fmt.Printf("Completed in %s.\n", time.Since(startTime))
	fmt.Println()

	if c.Bool(utils.YesFlag.Name) || utils.Confirm(fmt.Sprintf("Would you like to restart your Stakewise Validator Client so it loads the new keys?\n%sYou will not be able to attest with these new keys until your VC is restarted and the keys are loaded into it.%s", terminal.ColorYellow, terminal.ColorReset)) {
		fmt.Print("Restarting Validator Client... ")
		_, err = hd.Api.Service.RestartContainer(swconfig.VcContainerSuffix)
		if err != nil {
			fmt.Println("error")
			fmt.Printf("%sWARNING: error restarting validator client: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			fmt.Println("Please restart your Validator Client in order to attest with your new keys!")
		} else {
			fmt.Println("done!")
		}
	} else {
		fmt.Println("Please restart your Validator Client at your earliest convenience in order to attest with your new keys.")
	}
	fmt.Println()

	// Regenerate the deposit data
	err = regenDepositData(c, sw)
	if err != nil {
		return err
	}

	if c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to restart the Stakewise Operator service so it loads the new keys and deposit data?") {
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
	} else {
		fmt.Println("Please restart the container at your convenience.")
	}
	fmt.Println()

	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to upload the deposit data with the new keys to the NodeSet server, so they can be used for new validator assignments?")) {
		fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
		return nil
	}

	fmt.Printf("Uploading deposit data to the NodeSet server... ")

	_, err = sw.Api.Nodeset.UploadDepositData()
	if err != nil {
		fmt.Println("error")
		fmt.Printf("%sWARNING: error uploading deposit data to nodeset: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
		fmt.Println("Please upload the deposit data for all of your keys with `hyperdrive stakewise nodeset upload-deposit-data` when you're ready. Without it, NodeSet won't be able to assign new deposits to your validators.")
		fmt.Println()
	} else {
		fmt.Println("done!")
	}

	return nil
}

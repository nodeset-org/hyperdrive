package wallet

import (
	"fmt"
	"time"

	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	swcmdutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	cliutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/utils/input"
	"github.com/urfave/cli/v2"
)

var (
	generateKeysCountFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "count",
		Aliases: []string{"c"},
		Usage:   "The number of keys to generate",
	}
	generateKeysNoRestartFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:  "no-restart",
		Usage: fmt.Sprintf("Don't automatically restart the Stakewise Operator or Validator Client containers after generating keys. %sOnly use this if you know what you're doing and can restart them manually.%s", terminal.ColorRed, terminal.ColorReset),
	}
)

func generateKeys(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}
	noRestart := c.Bool(generateKeysNoRestartFlag.Name)

	// Check wallet status
	_, ready, err := utils.CheckIfWalletReady(hd)
	if err != nil {
		return err
	}
	if !ready {
		return nil
	}

	// Get the count
	count := c.Uint64(generateKeysCountFlag.Name)
	if count == 0 {
		countString := cliutils.Prompt("How many keys would you like to generate?", "^\\d+$", "Invalid count, try again")
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
		fmt.Printf("Generated %s (%d/%d) in %s\n", pubkey.HexWithPrefix(), (i + 1), count, elapsed)
	}
	fmt.Printf("Completed in %s.\n", time.Since(startTime))
	fmt.Println()

	// Restart the Stakewise Operator
	if noRestart {
		fmt.Printf("%sYou have automatic restarting turned off.\nPlease restart your Stakewise Operator container at your earliest convenience in order to deposit your new keys once it's your turn. Failure to do so will prevent your validators from ever being activated.%s\n", terminal.ColorYellow, terminal.ColorReset)
	} else {
		fmt.Print("Restarting Stakewise Operator... ")
		_, err = hd.Api.Service.RestartContainer(string(swconfig.ContainerID_StakewiseOperator))
		if err != nil {
			fmt.Println("error")
			fmt.Printf("%sWARNING: error restarting stakewise operator: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			fmt.Println("Please restart your Stakewise Operator container in order to be able to deposit for your new keys,")
		} else {
			fmt.Println("done!")
		}
	}
	fmt.Println()

	// Restart the VC
	if noRestart {
		fmt.Printf("%sYou have automatic restarting turned off.\nPlease restart your Validator Client at your earliest convenience in order to attest with your new keys. Failure to do so will result in any new validators being offline and *losing ETH* until you restart it.%s\n", terminal.ColorYellow, terminal.ColorReset)
	} else {
		fmt.Print("Restarting Validator Client... ")
		_, err = hd.Api.Service.RestartContainer(string(swconfig.ContainerID_StakewiseValidator))
		if err != nil {
			fmt.Println("error")
			fmt.Printf("%sWARNING: error restarting validator client: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			fmt.Println("Please restart your Validator Client in order to attest with your new keys!")
		} else {
			fmt.Println("done!")
		}
	}
	fmt.Println()

	// Upload to the server
	newKeysUploaded, err := swcmdutils.UploadDepositData(c, hd, sw)
	if err != nil {
		return err
	}

	if newKeysUploaded {
		if !noRestart {
			fmt.Println()
			fmt.Println("Your new keys are now ready for use. When one of them is selected for activation, your system will deposit it and begin attesting automatically.")
		} else {
			fmt.Println("Your new keys are uploaded, but you *must* restart your Validator Client at your earliest convenience to begin attesting once they are selected for depositing.")
		}
	}

	return nil
}

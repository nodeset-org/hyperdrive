package wallet

import (
	"fmt"
	"time"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	swcmdutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
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
	hd := client.NewHyperdriveClientFromCtx(c)
	sw := client.NewStakewiseClientFromCtx(c)
	noRestart := c.Bool(generateKeysNoRestartFlag.Name)

	// Make sure there's a wallet loaded
	response, err := hd.Api.Wallet.Status()
	if err != nil {
		return fmt.Errorf("error checking wallet status: %w", err)
	}
	status := response.Data.WalletStatus
	if !status.Wallet.IsLoaded {
		if !status.Wallet.IsOnDisk {
			fmt.Println("Your node wallet has not been initialized yet. Please run `hyperdrive wallet init` first to create it, then run this again.")
			return nil
		}
		if !status.Password.IsPasswordSaved {
			fmt.Println("Your node wallet has been initialized, but Hyperdrive doesn't have a password loaded for it so it cannot be used. Please run `hyperdrive wallet set-password` to enter it, then run this command again.")
			return nil
		}
		return nil
	}

	// Get the count
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
	err = swcmdutils.UploadDepositData(sw)
	if err != nil {
		return err
	}

	if !noRestart {
		fmt.Println()
		fmt.Println("Your new keys are now ready for use. When one of them is selected for activation, your system will deposit it and begin attesting automatically.")
	} else {
		fmt.Println("Your new keys are uploaded, but you *must* restart your Validator Client at your earliest convenience to begin attesting once they are selected for depositing.")
	}

	return nil
}

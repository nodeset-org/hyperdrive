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
		Usage: fmt.Sprintf("Don't automatically restart the Validator Client after generating keys. %sOnly use this if you know what you're doing and can restart it manually.%s", terminal.ColorYellow, terminal.ColorReset),
	}
)

func generateKeys(c *cli.Context) error {
	hd := client.NewClientFromCtx(c)
	sw := client.NewStakewiseClientFromCtx(c)
	noRestart := c.Bool(generateKeysNoRestartFlag.Name)

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
		fmt.Printf("Generated %s (%d/%d) in %s\n", pubkey.HexWithPrefix(), (i + 1), count, elapsed)
	}
	fmt.Printf("Completed in %s.\n", time.Since(startTime))
	fmt.Println()

	// Restart the VC
	if noRestart {
		fmt.Printf("%sYou have automatic restarting turned off.\nPlease restart your Validator Client at your earliest convenience in order to attest with your new keys. Failure to do so will result in any new validators being offline and *losing ETH* until you restart it.%s\n", terminal.ColorYellow, terminal.ColorReset)
	} else {
		fmt.Print("Restarting Validator Client... ")
		_, err = hd.Api.Service.RestartContainer(swconfig.VcContainerSuffix)
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
		fmt.Println("Your new keys are now ready for use. When NodeSet selects one of them for a new deposit, your system will deposit it and begin attesting automatically.")
	} else {
		fmt.Println("Your new keys are uploaded, but you *must* restart your Validator Client at your earliest convenience to begin attesting once they are selected for depositing.")
	}

	return nil
}

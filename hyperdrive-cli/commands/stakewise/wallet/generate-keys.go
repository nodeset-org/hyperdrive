package wallet

import (
	"fmt"
	"time"

	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	cliutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/utils/input"
	"github.com/urfave/cli/v2"
)

var (
	generateKeysCountFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "count",
		Aliases: []string{"c"},
		Usage:   "The number of keys to generate",
	}
	noRestartFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:  "no-restart",
		Usage: fmt.Sprintf("Don't automatically restart the Validator Client after the operation. %sOnly use this if you know what you're doing and can restart it manually.%s", terminal.ColorRed, terminal.ColorReset),
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
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading Hyperdrive config: %w", err)
	}
	if !cfg.StakeWise.Enabled.Value {
		fmt.Println("The StakeWise module is not enabled in your Hyperdrive configuration.")
		return nil
	}

	// Check wallet status
	_, ready, err := cliutils.CheckIfWalletReady(hd)
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
	newPubkeys := make([]beacon.ValidatorPubkey, count)
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
		newPubkeys[i] = pubkey
		fmt.Printf("Generated %s (%d/%d) in %s\n", pubkey.HexWithPrefix(), (i + 1), count, elapsed)
	}
	fmt.Printf("Completed in %s.\n", time.Since(startTime))
	fmt.Println()

	// Get the list of available keys
	response, err := sw.Api.Wallet.GetAvailableKeys(false)
	if err != nil {
		return fmt.Errorf("error getting available keys: %w", err)
	}
	data := response.Data
	availableKeyMap := map[beacon.ValidatorPubkey]struct{}{}
	for _, key := range data.AvailablePubkeys {
		availableKeyMap[key] = struct{}{}
	}

	// Sort into new and already used keys
	newKeys := []beacon.ValidatorPubkey{}
	usedKeys := []beacon.ValidatorPubkey{}
	for _, key := range newPubkeys {
		if _, exists := availableKeyMap[key]; exists {
			newKeys = append(newKeys, key)
		} else {
			usedKeys = append(usedKeys, key)
		}
	}

	// Print warnings about used keys
	if len(usedKeys) > 0 {
		fmt.Printf("%sNOTE: %d of the new keys belong to existing validators and cannot be used:%s\n", terminal.ColorYellow, len(usedKeys), terminal.ColorReset)
		for _, key := range usedKeys {
			fmt.Println("\t" + key.HexWithPrefix())
		}
		fmt.Println()
	}
	if len(newKeys) == 0 {
		fmt.Println("None of the keys can be used for new deposits. Please run this command again to generate more keys.")
		return nil
	}

	fmt.Printf("You now have %s%d%s validator keys ready for deposits:\n", terminal.ColorGreen, len(data.AvailablePubkeys), terminal.ColorReset)
	for _, key := range data.AvailablePubkeys {
		fmt.Println("\t" + key.HexWithPrefix())
	}
	fmt.Println()
	if !data.SufficientBalance {
		fmt.Println()
		fmt.Printf("%sWARNING: your wallet has less ETH than StakeWise recommends (%.2f ETH per key).%s\n", terminal.ColorYellow, data.EthPerKey, terminal.ColorReset)
		fmt.Printf("Current wallet balance: %s%f%s\n", terminal.ColorGreen, data.Balance, terminal.ColorReset)
		fmt.Printf("You need %s%f%s more ETH to use all of these keys.\n", terminal.ColorGreen, data.RemainingEthRequired, terminal.ColorReset)
		fmt.Println()
	}

	// Restart the Stakewise Operator
	/* TODO: possibly not needed with SW v2 now
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
	*/

	// Restart the VC
	if c.Bool(noRestartFlag.Name) {
		fmt.Printf("%sYou have automatic restarting turned off.\nPlease restart your Validator Client at your earliest convenience in order to attest with your new keys. Failure to do so will result in any new validators being offline and *losing ETH* until you restart it.%s\n", terminal.ColorYellow, terminal.ColorReset)
	} else {
		fmt.Print("Restarting Validator Client to load the new keys... ")
		_, err = hd.Api.Service.RestartContainer(string(swconfig.ContainerID_StakewiseValidator))
		if err != nil {
			fmt.Println("error")
			fmt.Printf("%sWARNING: error restarting validator client: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			fmt.Println("Please restart your Validator Client in order to attest with your new keys!")
		} else {
			fmt.Println("done!")
			fmt.Println("Your new keys are now loaded.")
			fmt.Println("Your node will deposit with them automatically once the vault has been funded.")
			fmt.Println("It will start attesting for those validators automatically once they have been activated.")
		}
	}
	return nil
}

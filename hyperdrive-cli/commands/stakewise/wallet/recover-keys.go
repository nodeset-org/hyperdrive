package wallet

import (
	"fmt"

	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/urfave/cli/v2"
)

const (
	// The limit for a single instance of key recovery
	singleRecoverSearchLimit uint64 = 5
)

var (
	startIndexFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "start-index",
		Aliases: []string{"i"},
		Usage:   "The index to start recovering keys from. Default is 0.",
		Value:   0,
	}
	searchLimitFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "search-limit",
		Aliases: []string{"l"},
		Usage:   "The maximum number of continuous keys to search unsuccessfully before stopping. Once a key is found, this limit will reset and key recovery will continue.",
		Value:   100,
	}
)

func recoverKeys(c *cli.Context) error {
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
	_, ready, err := utils.CheckIfWalletReady(hd)
	if err != nil {
		return err
	}
	if !ready {
		return nil
	}

	// Print status info
	response, err := sw.Api.Wallet.GetRegisteredKeys()
	if err != nil {
		return fmt.Errorf("error getting registered keys: %w", err)
	}
	data := response.Data
	if data.NotRegisteredWithNodeSet {
		fmt.Println("Your wallet is not registered with NodeSet. Please register with `hyperdrive nodeset register-node`.")
		return nil
	}
	if data.InvalidPermissions {
		fmt.Println("Your node currently doesn't have permissions to access the vaults on this deployment.")
		return nil
	}
	if len(data.Vaults) == 0 {
		fmt.Println("This deployment doesn't have any StakeWise vaults yet.")
		return nil
	}

	// Check if there are any keys to recover
	hasKeys := false
	for _, vault := range data.Vaults {
		if vault.HasPermission && len(vault.Validators) > 0 {
			hasKeys = true
			break
		}
	}
	if !hasKeys {
		fmt.Println("You do not have any validator keys registered yet, so there's nothing to recover.")
	}

	// Print the vault info
	keysToRecover := []beacon.ValidatorPubkey{}
	fmt.Println("The following validator keys have been registered and can be recovered:")
	for _, vault := range data.Vaults {
		if !vault.HasPermission {
			continue
		}
		if len(vault.Validators) == 0 {
			continue
		}
		fmt.Printf("%s (%s%s%s):\n", vault.Name, terminal.ColorGreen, vault.Address.Hex(), terminal.ColorReset)
		for _, key := range vault.Validators {
			fmt.Printf("\t%s\n", key.HexWithPrefix())
			keysToRecover = append(keysToRecover, key)
		}
		fmt.Println()
	}

	// Prompt for confirmation
	fmt.Println("NOTE: Key recovering may take a long time. Progress will be printed after checking every 5 keys.")
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Are you ready to begin key recovery?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Recover the keys
	startIndex := c.Uint64(startIndexFlag.Name)
	searchLimit := c.Uint64(searchLimitFlag.Name)
	keyMap := make(map[beacon.ValidatorPubkey]struct{})
	nextEndIndex := startIndex + searchLimit
	for _, key := range keysToRecover {
		keyMap[key] = struct{}{}
	}
	keysRecovered := false
	for len(keyMap) > 0 {
		fmt.Printf("Searching index %d to %d...\n", startIndex, startIndex+singleRecoverSearchLimit-1)
		response, err := sw.Api.Wallet.RecoverKeys(keysToRecover, startIndex, 1, singleRecoverSearchLimit, false)
		if err != nil {
			return fmt.Errorf("error recovering keys: %w", err)
		}
		data := response.Data
		if data.NotRegisteredWithNodeSet {
			fmt.Println("Your wallet is not registered with NodeSet. Please register with `hyperdrive nodeset register-node`.")
			return nil
		}
		for _, key := range data.Keys {
			delete(keyMap, key.Pubkey)
			fmt.Printf("Recovered %s (index %d)\n", key.Pubkey.HexWithPrefix(), key.Index)
			keysRecovered = true
			nextEndIndex = data.SearchEnd + 1 + searchLimit
		}

		startIndex = data.SearchEnd + 1
		if startIndex > nextEndIndex {
			fmt.Println("Reached the search limit. Stopping key recovery.")
			break
		}
	}
	fmt.Println("Key recovery complete.")
	fmt.Println()

	if keysRecovered {
		// Restart the VC
		if c.Bool(noRestartFlag.Name) {
			fmt.Printf("%sYou have automatic restarting turned off.\nPlease restart your Validator Client at your earliest convenience in order to attest with your recovered keys. Failure to do so will result in the validators being offline and *losing ETH* until you restart it.%s\n", terminal.ColorYellow, terminal.ColorReset)
		} else {
			fmt.Print("Restarting Validator Client to load the recovered keys... ")
			_, err = hd.Api.Service.RestartContainer(string(swconfig.ContainerID_StakewiseValidator))
			if err != nil {
				fmt.Println("error")
				fmt.Printf("%sWARNING: error restarting validator client: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
				fmt.Println("Please restart your Validator Client in order to attest with your recovered keys!")
			} else {
				fmt.Println("done!")
				fmt.Println("Your recovered keys are now loaded.")
				fmt.Println("Your node can now attest for these validators.")
			}
		}
	}

	return nil
}

package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func getAvailableKeys(c *cli.Context) error {
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

	// Get the list of available keys
	response, err := sw.Api.Wallet.GetAvailableKeys()
	if err != nil {
		return fmt.Errorf("error getting available keys: %w", err)
	}

	// Print the available keys and balance info
	data := response.Data
	newKeyCount := len(response.Data.AvailablePubkeys)
	if newKeyCount == 0 {
		fmt.Println("You do not have any validator keys ready for deposits.")
		fmt.Println("Please generate new keys with `hyperdrive stakewise wallet generate-keys`.")
		return nil
	}
	fmt.Printf("You have %s%d%s validator keys ready for deposits:\n", terminal.ColorGreen, newKeyCount, terminal.ColorReset)
	for _, key := range data.AvailablePubkeys {
		fmt.Println("\t" + key.HexWithPrefix())
	}
	if !data.SufficientBalance {
		fmt.Println()
		fmt.Printf("%sWarning: your wallet has less ETH than StakeWise recommends (%.2f ETH per key).%s\n", terminal.ColorYellow, data.EthPerKey, terminal.ColorReset)
		fmt.Printf("Current wallet balance: %s%f%s\n", terminal.ColorGreen, data.Balance, terminal.ColorReset)
		fmt.Printf("You need %s%f%s more ETH to use all of these keys.\n", terminal.ColorGreen, data.RemainingEthRequired, terminal.ColorReset)
	}
	return nil
}

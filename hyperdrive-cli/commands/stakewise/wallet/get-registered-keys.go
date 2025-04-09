package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func getRegisteredKeys(c *cli.Context) error {
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

	// Print the vault info
	for _, vault := range data.Vaults {
		fmt.Printf("%s (%s%s%s):\n", vault.Name, terminal.ColorGreen, vault.Address.Hex(), terminal.ColorReset)
		if !vault.HasPermission {
			fmt.Println("\tYou do not have permission to access this vault.")
			fmt.Println()
			continue
		}
		if len(vault.Validators) == 0 {
			fmt.Println("\tYou do not have any validators registered with this vault.")
			fmt.Println()
			continue
		}
		for _, key := range vault.Validators {
			fmt.Printf("\t%s\n", key.HexWithPrefix())
		}
		fmt.Println()
	}
	return nil
}

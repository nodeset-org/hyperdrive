package network

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/nodeset"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/urfave/cli/v2"
)

func getStatus(c *cli.Context) error {
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

	// Check the registration status first
	shouldContinue, err := nodeset.CheckRegistrationStatus(c, hd)
	if err != nil {
		return fmt.Errorf("error checking nodeset registration status: %w", err)
	}
	if !shouldContinue {
		return nil
	}

	// Get the network status
	response, err := sw.Api.Network.Status()
	if err != nil {
		return fmt.Errorf("error fetching network status: %w", err)
	}
	if response.Data.NotRegisteredWithNodeSet {
		fmt.Println("You are not registered with NodeSet yet.")
		fmt.Println("Please register with `hyperdrive nodeset register-node`.")
		return nil
	}
	if response.Data.InvalidPermissions {
		fmt.Println("Your node doesn't have permission to use the StakeWise module yet.")
		return nil
	}
	if len(response.Data.Vaults) == 0 {
		fmt.Println("There are no StakeWise vaults for this deployment yet.")
		return nil
	}

	// Print the network info
	fmt.Println()
	for _, vault := range response.Data.Vaults {
		fmt.Printf("%s (%s%s%s):\n", vault.Name, terminal.ColorGreen, vault.Address, terminal.ColorReset)
		fmt.Printf("\tMax Validators per User:    %d\n", vault.MaxValidators)
		fmt.Printf("\tYour Registered Validators: %d\n", vault.RegisteredValidators)
		fmt.Printf("\tYour Available Validators:  %d\n", vault.AvailableValidators)
		fmt.Printf("\tETH Available for staking:  %.6f\n", eth.WeiToEth(vault.Balance))
		fmt.Println()
	}

	return nil
}

package wallet

import (
	"fmt"

	"github.com/nodeset-org/eth-utils/eth"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func getStatus(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Get the config
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("Error loading configuration: %w", err)
	}

	// Print what network we're on
	err = utils.PrintNetwork(cfg.Hyperdrive.Network.Value, isNew)
	if err != nil {
		return err
	}

	// Get wallet response
	response, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}

	// Print status & return
	status := response.Data.WalletStatus
	if !status.Address.HasAddress {
		fmt.Println("The node wallet has not been initialized with an address yet.")
		return nil
	}
	if !status.Wallet.IsLoaded {
		if !status.Wallet.IsOnDisk {
			fmt.Println("The node wallet has not been initialized yet.")
			fmt.Printf("Your node is currently masquerading as %s%s%s.\n", terminal.ColorBlue, status.Address.NodeAddress.Hex(), terminal.ColorReset)
			fmt.Printf("%sIt is running in 'read-only' mode and cannot transact, as does not have that node's private wallet key.%s\n", terminal.ColorYellow, terminal.ColorReset)
			return nil
		}
		if !status.Password.IsPasswordSaved {
			fmt.Println("The node wallet has been initialized, but Hyperdrive doesn't have a password loaded for your node wallet so it cannot be used.")
			fmt.Printf("Your node is currently running as %s%s%s in %s'read-only' mode%s.\n", terminal.ColorBlue, status.Address.NodeAddress.Hex(), terminal.ColorReset, terminal.ColorYellow, terminal.ColorReset)
			return nil
		}
	}

	if status.Address.NodeAddress != status.Wallet.WalletAddress {
		fmt.Printf("The node wallet is initialized, but you are currently masquerading as %s%s%s.\n", terminal.ColorBlue, status.Address.NodeAddress.Hex(), terminal.ColorReset)
		fmt.Printf("Your node wallet is for %s%s%s.\n", terminal.ColorBlue, status.Wallet.WalletAddress.Hex(), terminal.ColorReset)
		fmt.Printf("%sDue to this mismatch, your node is running in 'read-only' mode and cannot submit transactions.%s\n", terminal.ColorYellow, terminal.ColorReset)
	} else {
		fmt.Println("The node wallet is initialized and ready.")
		fmt.Printf("Node account: %s%s%s\n", terminal.ColorGreen, status.Wallet.WalletAddress.Hex(), terminal.ColorReset)
		fmt.Printf("%sThe node's wallet keystore matches this address; it will be able to submit transactions.%s", terminal.ColorGreen, terminal.ColorReset)
	}

	fmt.Println()
	if status.Password.IsPasswordSaved {
		fmt.Printf("The node wallet's password %sis saved to disk%s.\n", terminal.ColorGreen, terminal.ColorReset)
		fmt.Println("The node will be able to submit transactions automatically after a restart.")
	} else {
		fmt.Printf("The node wallet's password %sis not saved to disk%s.\n", terminal.ColorYellow, terminal.ColorReset)
		fmt.Println("You will have to manually re-enter it with `hyperdrive wallet set-password` after a restart to be able to submit transactions.")
	}

	// Try to get the wallet balance
	balanceResponse, err := hd.Api.Wallet.Balance()
	if err != nil {
		if hd.Context.DebugEnabled {
			fmt.Printf("The node address's ETH balance is currently unavailable (%s).\n", err.Error())
		} else {
			fmt.Println("The node address's ETH balance is currently unavailable.")
		}
		return nil
	}

	fmt.Printf("Address %s%s%s's balance is %s%.6f%s ETH.\n", terminal.ColorBlue, status.Address.NodeAddress.Hex(), terminal.ColorReset, terminal.ColorGreen, eth.WeiToEth(balanceResponse.Data.Balance), terminal.ColorReset)

	return nil
}

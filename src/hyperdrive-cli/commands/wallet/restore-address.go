package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func restoreAddress(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c).WithReady()
	if err != nil {
		return err
	}

	// Get the wallet status
	response, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}
	status := response.Data.WalletStatus

	if !status.Wallet.IsLoaded {
		fmt.Println("You do not currently have a node wallet loaded, so there is no address to restore. Please see `hyperdrive wallet status` for more details.")
		return nil
	}
	if status.Wallet.WalletAddress == status.Address.NodeAddress {
		fmt.Println("Your node address is set to your wallet address; you are not currently masquerading.")
		return nil
	}

	fmt.Printf("Your node wallet is %s%s%s. If you restore it, you will no longer be masquerading as %s%s%s.\n\n", terminal.ColorBlue, status.Wallet.WalletAddress.Hex(), terminal.ColorReset, terminal.ColorBlue, status.Address.NodeAddress, terminal.ColorReset)

	// Confirm
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Are you sure you want to end your masquerade and restore your node address to your wallet address?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Run it
	_, err = hd.Api.Wallet.RestoreAddress()
	if err != nil {
		return fmt.Errorf("error restoring address: %w", err)
	}

	fmt.Printf("Your node address has been reset to your wallet address.")
	return nil
}

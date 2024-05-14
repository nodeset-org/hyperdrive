package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func initialize(c *cli.Context) error {
	// Get client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}

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

	// Initialize the Stakewise wallet
	swResponse, err := sw.Api.Wallet.Initialize()
	if err != nil {
		return err
	}

	fmt.Printf("Your node wallet has been successfully copied to the Stakewise module with address %s%s%s.", terminal.ColorBlue, swResponse.Data.AccountAddress.Hex(), terminal.ColorReset)
	return nil
}

package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
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

	// Check wallet status
	_, ready, err := utils.CheckIfWalletReady(hd)
	if err != nil {
		return err
	}
	if !ready {
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

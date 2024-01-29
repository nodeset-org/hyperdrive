package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func initialize(c *cli.Context) error {
	// Get Stakewise client
	hd := client.NewStakewiseClientFromCtx(c)

	//
	response, err := hd.Api.Wallet.Initialize()
	if err != nil {
		return err
	}

	fmt.Printf("Your node wallet has been successfully copied to the Stakewise module with address %s%s%s.", terminal.ColorBlue, response.Data.AccountAddress.Hex(), terminal.ColorReset)
	return nil
}

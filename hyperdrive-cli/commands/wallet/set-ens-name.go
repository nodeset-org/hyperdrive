package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	"github.com/urfave/cli/v2"
)

func setEnsName(c *cli.Context, name string) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	fmt.Printf("This will confirm the node's ENS name as '%s'.\n\n%sNOTE: to confirm your name, you must first register it with the ENS application at https://app.ens.domains.\nWe recommend using a hardware wallet as the base domain, and registering your node as a subdomain of it.%s\n\n", name, terminal.ColorYellow, terminal.ColorReset)

	// Build the TX
	response, err := hd.Api.Wallet.SetEnsName(name)
	if err != nil {
		return err
	}

	// Run the TX
	validated, err := tx.HandleTx(c, hd, response.Data.TxInfo,
		"Are you sure you want to confirm your node's ENS name?",
		"setting ENS name",
		"Setting ENS name...",
	)
	if err != nil {
		return err
	}
	if !validated {
		return nil
	}

	fmt.Printf("The ENS name associated with your node account is now '%s'.\n\n", name)
	return nil
}

package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

func deletePassword(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Get & check wallet status
	statusResponse, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}
	status := statusResponse.Data.WalletStatus

	// Check if it's already set
	if !status.Password.IsPasswordSaved {
		fmt.Println("The node wallet password is not saved to disk.")
		return nil
	}

	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Are you sure you want to delete your password from disk? Your node will not be able to submit transactions after a restart until you manually enter the password")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Run it
	_, err = hd.Api.Wallet.DeletePassword()
	if err != nil {
		return fmt.Errorf("error deleting password: %w", err)
	}

	// Log & return
	fmt.Println("The password has been successfully removed from disk storage.")
	return nil
}

package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	"github.com/urfave/cli/v2"
)

func setPassword(c *cli.Context) error {
	// Get RP client
	hd := client.NewClientFromCtx(c)

	// Get & check wallet status
	statusResponse, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}
	status := statusResponse.Data.WalletStatus

	// Check if it's already set
	if status.IsPasswordSaved {
		fmt.Println("The node wallet password is already loaded and saved to disk.")
		return nil
	}

	if status.HasKeystore {

	}

	// Get the password
	passwordString := c.String(passwordFlag.Name)
	if passwordString == "" {
		passwordString = promptPassword()
	}
	password, err := input.ValidateNodePassword("password", passwordString)
	if err != nil {
		return fmt.Errorf("error validating password: %w", err)
	}

	// Get the save flag
	savePassword := c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to save the password to disk? If you do, your node will be able to handle transactions automatically after a client restart; otherwise, you will have to repeat this command to manually enter the password after each restart.")

	// Run it
	_, err = hd.Api.Wallet.SetPassword(password, savePassword)
	if err != nil {
		return fmt.Errorf("error setting password: %w", err)
	}

	// Log & return
	fmt.Println("The password has been successfully uploaded to the daemon.")
	return nil
}

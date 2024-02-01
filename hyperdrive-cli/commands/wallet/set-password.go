package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	"github.com/urfave/cli/v2"
)

func setPassword(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Get & check wallet status
	statusResponse, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}
	status := statusResponse.Data.WalletStatus

	// Check if it's already set properly and the wallet has been loaded
	if status.Wallet.IsLoaded {
		if status.Password.IsPasswordSaved {
			fmt.Println("The node wallet password is already loaded and saved to disk.")
			return nil
		}
		fmt.Println("The node wallet is loaded, but the password is not saved to disk.")
	}

	// Get the password
	passwordString := c.String(PasswordFlag.Name)
	if passwordString == "" {
		if status.Wallet.IsOnDisk {
			passwordString = PromptExistingPassword()
		} else {
			passwordString = PromptNewPassword()
		}
	}
	password, err := input.ValidateNodePassword("password", passwordString)
	if err != nil {
		return fmt.Errorf("error validating password: %w", err)
	}

	// Get the save flag
	savePassword := c.Bool(SavePasswordFlag.Name) || utils.Confirm("Would you like to save the password to disk? If you do, your node will be able to handle transactions automatically after a client restart; otherwise, you will have to repeat this command to manually enter the password after each restart.")

	if status.Wallet.IsLoaded && !status.Password.IsPasswordSaved && !savePassword {
		fmt.Println("You've elected not to save the password but the node wallet is already loaded, so there's nothing to do.")
		return nil
	}

	// Run it
	_, err = hd.Api.Wallet.SetPassword(password, savePassword)
	if err != nil {
		return fmt.Errorf("error setting password: %w", err)
	}

	// Log & return
	if status.Wallet.IsLoaded {
		fmt.Println("The password has been successfully saved.")
	} else {
		fmt.Println("The password has been successfully uploaded to the daemon and the node wallet has been loaded.")
	}
	return nil
}

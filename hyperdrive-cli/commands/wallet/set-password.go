package wallet

import (
	"fmt"
	"log/slog"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/utils/input"
	"github.com/urfave/cli/v2"
)

func setPassword(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Get the config
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error getting Hyperdrive configuration: %w", err)
	}

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
	if !status.Wallet.IsOnDisk {
		fmt.Println("The node wallet has not been initialized yet. Please run `hyperdrive wallet init` or `hyperdrive wallet recover` first, then run this again.")
		return nil
	}

	// Print a debug log warning
	if cfg.Hyperdrive.Logging.Level.Value == slog.LevelDebug {
		fmt.Printf("%sWARNING: You have debug logging enabled. Your node's wallet password will be saved to the log file if you run this command.%s\n\n", terminal.ColorRed, terminal.ColorReset)
		if !utils.Confirm("Are you sure you want to continue?") {
			fmt.Println("Cancelled.")
			return nil
		}
		fmt.Println()
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

package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/term"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

var (
	initConfirmMnemonicFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "confirm-mnemonic",
		Aliases: []string{"c"},
		Usage:   "Automatically confirm the mnemonic phrase",
	}
)

// If hd is provided, this is assumed to be called from another function so the wallet check will be skipped.
func InitWallet(c *cli.Context, hd *client.HyperdriveClient) error {
	if hd == nil {
		// Get Hyperdrive client
		hd = client.NewHyperdriveClientFromCtx(c)

		// Get & check wallet status
		statusResponse, err := hd.Api.Wallet.Status()
		if err != nil {
			return err
		}
		status := statusResponse.Data.WalletStatus
		if status.Wallet.IsOnDisk {
			fmt.Println("The node wallet is already initialized.")
			return nil
		}
	}

	// Get the config
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error getting Hyperdrive configuration: %w", err)
	}

	// Prompt for user confirmation before printing sensitive information
	if !(hd.Context.SecureSession ||
		utils.ConfirmSecureSession("Creating a wallet will print sensitive information to your screen.")) {
		return nil
	}

	// Set password if not set
	var password string
	if c.String(PasswordFlag.Name) != "" {
		password = c.String(PasswordFlag.Name)
	} else {
		password = PromptNewPassword()
	}

	// Ask about saving
	savePassword := utils.Confirm("Would you like to save the password to disk? If you do, your node will be able to handle transactions automatically after a client restart; otherwise, you will have to manually enter the password after each restart with `hyperdrive wallet set-password`.")

	// Get the derivation path
	derivationPathString := c.String(derivationPathFlag.Name)
	var derivationPath *string
	if derivationPathString != "" {
		fmt.Printf("Using a custom derivation path (%s).\n\n", derivationPathString)
		derivationPath = &derivationPathString
	}

	// Get the wallet index
	walletIndexVal := c.Uint64(walletIndexFlag.Name)
	var walletIndex *uint64
	if walletIndexVal != 0 {
		fmt.Printf("Using a custom wallet index (%d).\n", walletIndex)
		walletIndex = &walletIndexVal
	}

	// Initialize wallet
	response, err := hd.Api.Wallet.Initialize(derivationPath, walletIndex, false, password, false)
	if err != nil {
		return fmt.Errorf("error initializing wallet: %w", err)
	}

	// Print mnemonic
	fmt.Println("Your mnemonic phrase to recover your wallet is printed below. It can be used to recover your node account and validator keys if they are lost.")
	fmt.Println("Record this phrase somewhere secure and private. Do not share it with anyone as it will give them control of your node account and validators.")
	fmt.Println("==============================================================================================================================================")
	fmt.Println("")
	fmt.Println(response.Data.Mnemonic)
	fmt.Println("")
	fmt.Println("==============================================================================================================================================")
	fmt.Println("")

	// Confirm mnemonic
	if !c.Bool(initConfirmMnemonicFlag.Name) {
		confirmMnemonic(response.Data.Mnemonic)
	}

	// Do a recover to verify and save the wallet
	recoverResponse, err := hd.Api.Wallet.Recover(derivationPath, &response.Data.Mnemonic, walletIndex, password, savePassword)
	if err != nil {
		return fmt.Errorf("error saving wallet: %w", err)
	}

	// Sanity check the addresses
	if recoverResponse.Data.AccountAddress != response.Data.AccountAddress {
		return fmt.Errorf("expected %s, but generated %s upon testing recovery", response.Data.AccountAddress, recoverResponse.Data.AccountAddress)
	}

	// Clear terminal output
	_ = term.Clear()

	fmt.Println("The node wallet was successfully initialized.")
	fmt.Printf("Node account: %s%s%s\n", terminal.ColorBlue, response.Data.AccountAddress.Hex(), terminal.ColorReset)

	// Initialize the Stakewise wallet if it's enabled
	if cfg.Stakewise.Enabled.Value {
		fmt.Println()
		fmt.Println("You have the Stakewise module enabled. Initializing it with your new wallet...")
		sw := client.NewStakewiseClientFromCtx(c)
		_, err = sw.Api.Wallet.Initialize()
		if err != nil {
			return fmt.Errorf("error initializing Stakewise wallet: %w", err)
		}
		fmt.Println("Stakewise wallet initialized.")
	}
	return nil
}

package wallet

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	swapi "github.com/nodeset-org/hyperdrive-stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func recoverWallet(c *cli.Context) error {
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
	if status.Wallet.IsOnDisk {
		fmt.Println("The node wallet is already initialized.")
		return nil
	}

	// Print a debug log warning
	if cfg.Hyperdrive.Logging.Level.Value == slog.LevelDebug {
		fmt.Printf("%sWARNING: You have debug logging enabled. Your mnemonic and node wallet password will be saved to the log file if you run this command.%s\n\n", terminal.ColorRed, terminal.ColorReset)
		if !utils.Confirm("Are you sure you want to continue?") {
			fmt.Println("Cancelled.")
			return nil
		}
		fmt.Println()
	}

	// Prompt a notice about test recovery
	fmt.Printf("%sNOTE:\nThis command will restore your node wallet's private key.\nIf you just want to test recovery to ensure it works without actually regenerating the files, please use `hyperdrive wallet test-recovery` instead.%s\n\n", terminal.ColorYellow, terminal.ColorReset)

	// Set password if not set
	var password string
	var savePassword bool
	if c.String(PasswordFlag.Name) != "" {
		password = c.String(PasswordFlag.Name)
	} else {
		password = PromptNewPassword()
	}

	// Ask about saving
	savePassword = utils.Confirm("Would you like to save the password to disk? If you do, your node will be able to handle transactions automatically after a client restart; otherwise, you will have to manually enter the password after each restart with `hyperdrive wallet set-password`.")

	// Prompt for mnemonic
	var mnemonic string
	if c.String(mnemonicFlag.Name) != "" {
		mnemonic = c.String(mnemonicFlag.Name)
	} else {
		mnemonic = PromptMnemonic()
	}
	mnemonic = strings.TrimSpace(mnemonic)

	// Check for a search-by-address operation
	addressString := c.String(addressFlag.Name)
	if addressString != "" {
		// Get the address to search for
		address := common.HexToAddress(addressString)
		fmt.Printf("Searching for the derivation path and index for wallet %s...\nNOTE: this may take several minutes depending on how large your wallet's index is.\n", address.Hex())

		// Recover wallet
		response, err := hd.Api.Wallet.SearchAndRecover(mnemonic, address, password, savePassword)
		if err != nil {
			return err
		}

		// Log & return
		fmt.Println("The node wallet was successfully recovered.")
		fmt.Printf("Derivation path: %s\n", response.Data.DerivationPath)
		fmt.Printf("Wallet index:    %d\n", response.Data.Index)
		fmt.Printf("Node account:    %s\n", response.Data.AccountAddress.Hex())
	} else {
		// Get the derivation path
		derivationPathString := c.String(derivationPathFlag.Name)
		var derivationPath *string
		if derivationPathString != "" {
			fmt.Printf("Using a custom derivation path (%s).\n", derivationPathString)
			derivationPath = &derivationPathString
		}

		// Get the wallet index
		walletIndexVal := c.Uint64(walletIndexFlag.Name)
		var walletIndex *uint64
		if walletIndexVal != 0 {
			fmt.Printf("Using a custom wallet index (%d).\n", walletIndex)
			walletIndex = &walletIndexVal
		}

		fmt.Println()
		fmt.Println("Recovering node wallet...")

		// Recover wallet
		response, err := hd.Api.Wallet.Recover(derivationPath, mnemonic, walletIndex, password, savePassword)
		if err != nil {
			return err
		}

		// Log & return
		fmt.Println("The node wallet was successfully recovered.")
		fmt.Printf("Node account: %s\n", response.Data.AccountAddress.Hex())
	}

	// Initialize the StakeWise wallet if it's enabled
	if cfg.Stakewise.Enabled.Value {
		fmt.Println()
		fmt.Println("You have the Stakewise module enabled. Initializing it with your new wallet...")
		sw, err := client.NewStakewiseClientFromCtx(c, hd)
		if err != nil {
			return err
		}
		_, err = sw.Api.Wallet.Initialize()
		if err != nil {
			return fmt.Errorf("error initializing Stakewise wallet: %w", err)
		}
		fmt.Println("Stakewise wallet initialized.")
		fmt.Println()

		// Check if the wallet is registered with NodeSet
		regResponse, err := sw.Api.Nodeset.RegistrationStatus()
		if err != nil {
			fmt.Println("Hyperdrive couldn't check your node's registration status:")
			fmt.Println(err.Error())
			fmt.Println("If your node isn't registered yet, you'll have to register it later.")
			return nil
		}
		switch regResponse.Data.Status {
		case swapi.NodesetRegistrationStatus_Registered:
			fmt.Println("Your node is already registered with NodeSet.")

		case swapi.NodesetRegistrationStatus_Unregistered:
			fmt.Println("Please whitelist your node on your `nodeset.io` dashboard, then register it with `hyperdrive sw ns register`.")

		case swapi.NodesetRegistrationStatus_Unknown:
			fmt.Println("Hyperdrive couldn't check your node's registration status:")
			fmt.Println(regResponse.Data.ErrorMessage)
			fmt.Println("If your node isn't registered yet, you'll have to register it later.")
		}

	}
	return nil
}

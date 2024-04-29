package wallet

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func recoverWallet(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

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

	return nil
}

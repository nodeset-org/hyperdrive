package wallet

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

func testRecovery(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Prompt a notice about test recovery
	fmt.Printf("%sNOTE:\nThis command will test the recovery of your node wallet's private key, but will not actually write any files; it's simply a \"dry run\" of recovery.\nUse `hyperdrive wallet recover` to actually recover the wallet.%s\n\n", terminal.ColorYellow, terminal.ColorReset)

	// Get the config
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error getting Hyperdrive configuration: %w", err)
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

		// Test recover wallet
		response, err := hd.Api.Wallet.TestSearchAndRecover(mnemonic, address)
		if err != nil {
			return err
		}

		// Log & return
		fmt.Println("The node wallet was successfully found - recovery is possible.")
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
		fmt.Println("Testing recovery of node wallet...")

		// Test recover wallet
		response, err := hd.Api.Wallet.TestRecover(derivationPath, mnemonic, walletIndex)
		if err != nil {
			return err
		}

		// Log & return
		fmt.Println("The node wallet was successfully found - recovery is possible.")
		fmt.Printf("Node account: %s\n", response.Data.AccountAddress.Hex())
	}

	return nil
}

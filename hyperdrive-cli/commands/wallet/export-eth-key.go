package wallet

import (
	"fmt"
	"os"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/rocket-pool/node-manager-core/wallet"
	"github.com/urfave/cli/v2"
)

func exportEthKey(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Get & check wallet status
	status, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}
	if !status.Data.WalletStatus.Wallet.IsLoaded {
		fmt.Println("The node wallet is not loaded and ready for usage. Please run `hyperdrive wallet status` for more details.")
		return nil
	}
	if status.Data.WalletStatus.Wallet.Type != wallet.WalletType_Local {
		fmt.Println("This command can only be run on local wallets; hardware wallets cannot have their keys exported.")
		return nil
	}

	if !hd.Context.SecureSession {
		// Check if stdout is interactive
		stat, err := os.Stdout.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while determining whether or not the output is a tty: %s\n"+
				"Use 'hyperdrive --secure-session wallet export-eth-key' to bypass.\n", err.Error())
			os.Exit(1)
		}

		if (stat.Mode()&os.ModeCharDevice) == os.ModeCharDevice &&
			!utils.ConfirmSecureSession("Exporting a wallet will print sensitive information to your screen.") {
			return nil
		}
	}

	// Get the wallet in ETH key format
	ethKey, err := hd.Api.Wallet.ExportEthKey()
	if err != nil {
		return err
	}

	// Print wallet & return
	fmt.Println("Wallet in ETH Key Format:")
	fmt.Println("============")
	fmt.Println()
	fmt.Println(string(ethKey.Data.EthKeyJson))
	fmt.Println()
	fmt.Println("============")
	fmt.Println()
	fmt.Println("Wallet password:")
	fmt.Println()
	fmt.Println(ethKey.Data.Password)
	fmt.Println()
	return nil
}

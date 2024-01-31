package wallet

import (
	"fmt"
	"os"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

func exportWallet(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Get & check wallet status
	status, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}
	if !status.Data.WalletStatus.Wallet.IsLoaded {
		fmt.Println("The node wallet is not loaded and ready for usage. Please run `hyperdrive wallet status` for more details.")
		return nil
	}

	if !hd.Context.SecureSession {
		// Check if stdout is interactive
		stat, err := os.Stdout.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while determining whether or not the output is a tty: %w\n"+
				"Use \"hyperdrive --secure-session wallet export\" to bypass.\n", err)
			os.Exit(1)
		}

		if (stat.Mode()&os.ModeCharDevice) == os.ModeCharDevice &&
			!utils.ConfirmSecureSession("Exporting a wallet will print sensitive information to your screen.") {
			return nil
		}
	}

	// Export wallet
	export, err := hd.Api.Wallet.Export()
	if err != nil {
		return err
	}

	// Print wallet & return
	fmt.Println("Node account private key:")
	fmt.Println("")
	fmt.Println(export.Data.AccountPrivateKey)
	fmt.Println("")
	fmt.Println("Wallet password:")
	fmt.Println("")
	fmt.Println(export.Data.Password)
	fmt.Println("")
	fmt.Println("Wallet file:")
	fmt.Println("============")
	fmt.Println("")
	fmt.Println(export.Data.Wallet)
	fmt.Println("")
	fmt.Println("============")
	return nil
}

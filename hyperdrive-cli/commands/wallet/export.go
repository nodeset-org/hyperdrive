package wallet

import (
	"fmt"
	"os"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	nutils "github.com/rocket-pool/node-manager-core/utils"
	"github.com/urfave/cli/v2"
)

func exportWallet(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Check wallet status
	_, ready, err := utils.CheckIfWalletReady(hd)
	if err != nil {
		return err
	}
	if !ready {
		return nil
	}

	if !hd.Context.SecureSession {
		// Check if stdout is interactive
		stat, err := os.Stdout.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while determining whether or not the output is a tty: %s\n"+
				"Use 'hyperdrive --secure-session wallet export' to bypass.\n", err.Error())
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
	fmt.Println()
	fmt.Println(nutils.EncodeHexWithPrefix(export.Data.AccountPrivateKey))
	fmt.Println()
	fmt.Println("Wallet file:")
	fmt.Println("============")
	fmt.Println()
	fmt.Println(export.Data.Wallet)
	fmt.Println()
	fmt.Println("============")
	fmt.Println()
	fmt.Println("Wallet password:")
	fmt.Println()
	if export.Data.Password == "" {
		fmt.Println("<Unknown - password not saved to disk>")
	} else {
		fmt.Println(export.Data.Password)
	}
	fmt.Println()
	return nil
}

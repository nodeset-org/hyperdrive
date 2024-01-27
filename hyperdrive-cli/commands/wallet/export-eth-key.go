package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func exportEthKey(c *cli.Context) error {
	// Get RP client
	hd := client.NewClientFromCtx(c)

	// Get & check wallet status
	status, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}
	if !status.Data.WalletStatus.HasKeystore {
		fmt.Println("The node wallet is not initialized.")
		return nil
	}

	// Get the wallet in ETH key format
	ethKey, err := hd.Api.Wallet.ExportEthKey()
	if err != nil {
		return err
	}

	// Print wallet & return
	fmt.Println("Wallet in ETH Key Format:")
	fmt.Println(string(ethKey.Data.EthKeyJson))
	return nil
}

package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	"github.com/urfave/cli/v2"
)

func claimRewards(c *cli.Context) error {
	fmt.Printf("Claiming rewards...\n")
	hd := client.NewHyperdriveClientFromCtx(c)
	sw := client.NewStakewiseClientFromCtx(c)
	resp, err := sw.Api.Wallet.ClaimRewards()
	if err != nil {
		return err
	}
	err = tx.HandleTx(c, hd, resp.Data.TxInfo,
		"Are you sure you want to set the validators root?",
		"setting validators root",
		"Setting the validators root...",
	)
	if err != nil {
		return err
	}
	fmt.Println("Rewards successfully claimed.")
	return nil
}

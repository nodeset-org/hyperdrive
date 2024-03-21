package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	"github.com/urfave/cli/v2"
)

func claimRewards(c *cli.Context) error {
	hd := client.NewHyperdriveClientFromCtx(c)
	sw := client.NewStakewiseClientFromCtx(c)
	resp, err := sw.Api.Wallet.ClaimRewards()
	if err != nil {
		return err
	}
	err = tx.HandleTx(c, hd, resp.Data.TxInfo,
		"Are you sure you want to claim rewards?",
		"claiming rewards",
		"Claiming rewards...",
	)
	if err != nil {
		return err
	}
	fmt.Println("Rewards successfully claimed.")
	return nil
}

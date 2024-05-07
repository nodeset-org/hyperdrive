package wallet

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/urfave/cli/v2"
)

func claimRewards(c *cli.Context) error {
	hd := client.NewHyperdriveClientFromCtx(c)
	sw := client.NewStakewiseClientFromCtx(c)
	resp, err := sw.Api.Wallet.ClaimRewards()
	if err != nil {
		return err
	}

	// Get the list of rewards available
	fmt.Println("Your withdrawable rewards:")
	fmt.Printf("%.4f %s (%s)\n", eth.WeiToEth(resp.Data.WithdrawableToken), resp.Data.TokenSymbol, resp.Data.TokenName)
	fmt.Printf("%.4f ETH\n", eth.WeiToEth(resp.Data.WithdrawableEth))
	fmt.Println()
	fmt.Println("NOTE: this list only shows rewards that Stakewise has already returned to NodeSet. Your share may include more rewards, but Stakewise hasn't returned yet.")
	fmt.Println()

	// Check if both balances are zero
	sum := big.NewInt(0)
	sum.Add(sum, resp.Data.WithdrawableEth)
	sum.Add(sum, resp.Data.WithdrawableToken)
	if sum.Cmp(common.Big0) == 0 {
		fmt.Println("You don't have any rewards to claim.")
		return nil
	}

	// Run the TX
	validated, err := tx.HandleTx(c, hd, resp.Data.TxInfo,
		"Are you sure you want to claim rewards?",
		"claiming rewards",
		"Claiming rewards...",
	)
	if err != nil {
		return err
	}
	if !validated {
		return nil
	}

	fmt.Println("Rewards successfully claimed.")
	return nil
}
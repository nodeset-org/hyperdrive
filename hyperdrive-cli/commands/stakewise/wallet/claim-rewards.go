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
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Get the list of rewards available
	resp, err := sw.Api.Wallet.ClaimRewards()
	if err != nil {
		return err
	}
	fmt.Println("Your withdrawable rewards:")
	fmt.Printf("%.4f %s (%s)\n", eth.WeiToEth(resp.Data.WithdrawableToken), resp.Data.TokenSymbol, resp.Data.TokenName)
	fmt.Printf("%.4f ETH\n", eth.WeiToEth(resp.Data.WithdrawableNativeToken))
	fmt.Println()
	fmt.Println("NOTE: this list only shows rewards that StakeWise has already returned to NodeSet. Your share may include more rewards, but StakeWise hasn't returned yet.")
	fmt.Println()

	// Check if both balances are zero
	sum := big.NewInt(0)
	sum.Add(sum, resp.Data.WithdrawableNativeToken)
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

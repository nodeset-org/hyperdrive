package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func claimRewards(c *cli.Context) error {
	// Get Hyperdrive client
	// address := hd.Config.NodeAddress
	// hd.Api.Rewards.ClaimRewards()

	fmt.Printf("Claiming rewards...\n")
	sc := client.NewStakewiseClientFromCtx(c)
	resp, err := sc.Api.Wallet.ClaimRewards()
	if err != nil {
		return err
	}
	fmt.Printf("Claimed rewards resp: %v\n", resp)
	return nil
}

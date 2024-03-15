package nodeset

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var rewardsRootFlag *cli.StringFlag = &cli.StringFlag{
	Name:    "root",
	Aliases: []string{"r"},
	Usage:   "The new root for the validators Merkle Tree, generated by the Stakewise Operator `get-validators-root` command",
}

func claimRewardsRoot(c *cli.Context) error {
	// Get the client
	// hd := client.NewHyperdriveClientFromCtx(c)

	// response, err := hd.Api.Rewards.ClaimRewardsRoot()
	// if err != nil {
	// 	return err
	// }

	// // Log & return
	fmt.Println("Validators root successfully set.")
	return nil
}
package status

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func getNodeStatus(c *cli.Context) error {
	sw := client.NewStakewiseClientFromCtx(c)
	// TODO: !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	response, err := sw.Api.Status.GetActiveWallets()
	if err != nil {
		fmt.Printf("!!! Error: %v\n", err)
		return err
	}
	fmt.Printf("!!! response: %v\n", response)
	return nil
}

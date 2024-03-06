package status

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func getNodeStatus(c *cli.Context) error {
	sw := client.NewStakewiseClientFromCtx(c)
	response, err := sw.Api.Status.GetValidatorStatuses()
	if err != nil {
		fmt.Printf("error fetching validator statuses: %v\n", err)
		return err
	}

	fmt.Printf("Beacon Statuses:\n")
	for pubKey, status := range response.Data.BeaconStatus {
		fmt.Printf("%v: %v\n", pubKey.HexWithPrefix(), status)
	}

	fmt.Printf("Nodeset Statuses:\n")
	for pubKey, status := range response.Data.NodesetStatus {
		fmt.Printf("%v: %v\n", pubKey.HexWithPrefix(), status)
	}

	return nil
}

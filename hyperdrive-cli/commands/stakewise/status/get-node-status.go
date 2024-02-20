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
		fmt.Printf("error fetching active validators: %v\n", err)
		return err
	}

	fmt.Printf("Active Validator Pubkeys: \n")

	for _, validator := range response.Data.ActiveValidators {
		fmt.Printf("%v\n", validator.HexWithPrefix())
	}

	return nil
}

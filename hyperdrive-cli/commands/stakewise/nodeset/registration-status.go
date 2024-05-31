package nodeset

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func registrationStatus(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	resp, err := sw.Api.Nodeset.RegistrationStatus()
	if err != nil {
		return err
	}
	if resp.Data.Registered {
		fmt.Println("Your node is registered.")
	} else {
		fmt.Println("Your node is not registered.")
	}

	return nil
}

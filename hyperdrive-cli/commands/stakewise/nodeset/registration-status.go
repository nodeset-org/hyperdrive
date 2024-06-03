package nodeset

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	cliutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
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
		if cliutils.Confirm("Would you like to upload your validator keys to NodeSet?") {
			if c.String(RegisterEmailFlag.Name) == "" {
				fmt.Printf("Please provide an email address with the %s flag.\n", RegisterEmailFlag.Name)
				return nil
			}
			_, err := sw.Api.Nodeset.RegisterNode(c.String(RegisterEmailFlag.Name))
			if err != nil {
				return err
			}
		}

	}

	return nil
}

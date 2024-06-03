package nodeset

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func registerNode(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	_, err = sw.Api.Nodeset.RegisterNode(RegisterEmailFlag.Name)

	return err
}

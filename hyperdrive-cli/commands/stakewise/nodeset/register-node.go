package nodeset

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

var (
	registerEmailFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "email",
		Aliases: []string{"e"},
		Usage:   "Email address to register with NodeSet.",
	}
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

	_, err = sw.Api.Nodeset.RegisterNode(registerEmailFlag.Name)

	return err
}

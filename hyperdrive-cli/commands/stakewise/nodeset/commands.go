package nodeset

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

var (
	registerEmailFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "email",
		Aliases: []string{"e"},
		Usage:   "Email address to register with NodeSet.",
	}
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your account with the Stakewise vault in NodeSet.",
		Subcommands: []*cli.Command{
			{
				Name:    "upload-deposit-data",
				Aliases: []string{"u"},
				Flags:   []cli.Flag{},
				Usage:   "Uploads the combined deposit data for all of your validator keys to NodeSet's Stakewise vault, so they can be assigned new deposits.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return uploadDepositData(c)
				},
			},
			{
				Name:    "register-node",
				Aliases: []string{"r"},
				Flags: []cli.Flag{
					registerEmailFlag,
				},
				Usage: "Register node with NodeSet",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return registerNode(c)
				},
			},
		},
	})
}

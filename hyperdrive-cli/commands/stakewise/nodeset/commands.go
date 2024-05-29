package nodeset

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
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
				Name:    "blah",
				Aliases: []string{"b"},
				Flags:   []cli.Flag{},
				Usage:   "Testing endpoint for new features",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return blah(c)
				},
			},
		},
	})
}

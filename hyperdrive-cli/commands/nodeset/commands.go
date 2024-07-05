package nodeset

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your node within the nodeset.io service.",
		Subcommands: []*cli.Command{
			{
				Name:    "registration-status",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					utils.YesFlag,
					RegisterEmailFlag,
				},
				Usage: "Check the registration status of your validator with NodeSet.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					return registrationStatus(c)
				},
			},
			{
				Name:    "register-node",
				Aliases: []string{"r"},
				Flags: []cli.Flag{
					RegisterEmailFlag,
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

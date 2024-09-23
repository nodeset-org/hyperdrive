package node

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your Constellation node",
		Subcommands: []*cli.Command{
			{
				Name:    "status",
				Aliases: []string{"s"},
				Flags:   []cli.Flag{},
				Usage:   "Get the node's status.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getStatus(c)
				},
			},
			{
				Name:    "register",
				Aliases: []string{"r"},
				Flags: []cli.Flag{
					utils.YesFlag,
				},
				Usage: "Registers your node with Constellation so you can create and run minipools.",
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

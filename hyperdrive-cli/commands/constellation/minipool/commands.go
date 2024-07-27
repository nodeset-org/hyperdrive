package minipool

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
			/*
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
			*/
			{
				Name:    "create",
				Aliases: []string{"c"},
				Flags: []cli.Flag{
					utils.YesFlag,
					saltFlag,
				},
				Usage: "Create a new minipool.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return createMinipool(c)
				},
			},
			/*
				{
					Name:    "stake",
					Aliases: []string{"k"},
					Flags: []cli.Flag{
						utils.YesFlag,
					},
					Usage: "Stake one or minipools that are still in prelaunch but have passed the Rocket Pool scrub check.",
					Action: func(c *cli.Context) error {
						// Validate args
						utils.ValidateArgCount(c, 0)

						// Run
						return stake(c)
					},
				},
			*/
		},
	})
}

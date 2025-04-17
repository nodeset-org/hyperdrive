package network

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage NodeSet's StakeWise network",
		Subcommands: []*cli.Command{
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "Get the status of the network",
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getStatus(c)
				},
			},
		},
	})
}

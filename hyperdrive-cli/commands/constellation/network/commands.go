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
		Usage:   "Interact with the Constellation network as a whole",
		Subcommands: []*cli.Command{
			{
				Name:    "stats",
				Aliases: []string{"s"},
				Flags:   []cli.Flag{},
				Usage:   "Get information about the Constellation network's current settings and stats",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getStats(c)
				},
			},
		},
	})
}

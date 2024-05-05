package validator

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your Stakewise validator keys",
		Subcommands: []*cli.Command{
			{
				Name:    "exit",
				Aliases: []string{"e"},
				Usage:   "Exit a validator",
				Flags: []cli.Flag{
					pubkeysFlag,
					epochFlag,
					noBroadcastFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return exit(c)
				},
			},
		},
	})
}

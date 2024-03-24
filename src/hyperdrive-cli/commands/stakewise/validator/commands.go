package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/utils"
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
					if err := utils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return exit(c)
				},
			},
		},
	})
}

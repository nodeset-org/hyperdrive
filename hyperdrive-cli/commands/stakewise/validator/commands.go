package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
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
				Name:    "get-signed-exit-messages",
				Aliases: []string{"s"},
				Usage:   "Get signed exit messages for one or more validators",
				Flags: []cli.Flag{
					pubkeysFlag,
					epochFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getSignedExitMessages(c)
				},
			},
		},
	})
}

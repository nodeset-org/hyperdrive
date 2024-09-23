package wallet

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your Constellation wallet",
		Subcommands: []*cli.Command{
			{
				Name:    "rebuild",
				Aliases: []string{"b"},
				Flags: []cli.Flag{
					rebuildAllFlag,
					rebuildStartIndexFlag,
					rebuildSearchLimitFlag,
					utils.YesFlag,
				},
				Usage: "Regenerate the validator private keys for your node's minipools and save them to disk. Useful for disaster recovery.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return rebuildValidatorKeys(c)
				},
			},
		},
	})
}

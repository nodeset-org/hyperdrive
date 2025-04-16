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
		Usage:   "Manage the Stakewise wallet",
		Subcommands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Clone the node wallet file into a wallet that the Stakewise operator service can use.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return initialize(c)
				},
			},
			{
				Name:    "get-available-keys",
				Aliases: []string{"a"},
				Usage:   "Retrieve the list of validator pubkeys that have been created locally and available for new deposits on this network.",
				Flags: []cli.Flag{
					lookbackFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getAvailableKeys(c)
				},
			},
			{
				Name:    "generate-keys",
				Aliases: []string{"g"},
				Usage:   "Generate new validator keys derived from your node wallet.",
				Flags: []cli.Flag{
					utils.YesFlag,
					generateKeysCountFlag,
					noRestartFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return generateKeys(c)
				},
			},
			{
				Name:    "recover-keys",
				Aliases: []string{"r"},
				Usage:   "Recover all of the registered validator keys derived from your node wallet.",
				Flags: []cli.Flag{
					utils.YesFlag,
					startIndexFlag,
					searchLimitFlag,
					noRestartFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return recoverKeys(c)
				},
			},
		},
	})
}

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
				Name:    "generate-keys",
				Aliases: []string{"g"},
				Usage:   "Generate new validator keys derived from your node wallet.",
				Flags: []cli.Flag{
					generateKeysCountFlag,
					generateKeysNoRestartFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return generateKeys(c)
				},
			},
			{
				Name:    "claim-rewards",
				Aliases: []string{"cr"},
				Usage:   "Claim rewards",
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					// Run
					return claimRewards(c)
				},
			},
		},
	})
}

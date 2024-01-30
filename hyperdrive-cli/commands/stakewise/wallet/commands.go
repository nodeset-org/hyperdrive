package wallet

import (
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
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
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

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
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return generateKeys(c)
				},
			},
			{
				Name:    "regen-deposit-data",
				Aliases: []string{"r"},
				Usage:   "Regenerate the combined deposit data for all of your validator keys and upload them to NodeSet's Stakewise vault, so they can be assigned new deposits.",
				Flags: []cli.Flag{
					regenDepositDataNoRestartFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return regenerateDepositData(c)
				},
			},
		},
	})
}

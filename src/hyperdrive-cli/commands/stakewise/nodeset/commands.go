package nodeset

import (
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your account with the Stakewise vault in NodeSet.",
		Subcommands: []*cli.Command{
			// DEBUG ONLY
			/*
				{
					Name:    "set-validators-root",
					Aliases: []string{"v"},
					Usage:   "(Temp cmd) Sets the validators root for the vault once a new aggregated deposit data file has been generated by the server. Use the Stakewise Operator's `get-validators-root` to generate it.",
					Flags: []cli.Flag{
						validatorsRootFlag,
					},
					Action: func(c *cli.Context) error {
						// Validate args
						if err := input.ValidateArgCount(c, 0); err != nil {
							return err
						}

						// Run
						return setValidatorsRoot(c)
					},
				},
			*/
			{
				Name:    "upload-deposit-data",
				Aliases: []string{"u"},
				Usage:   "Uploads the combined deposit data for all of your validator keys to NodeSet's Stakewise vault, so they can be assigned new deposits.",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return uploadDepositData(c)
				},
			},
		},
	})
}

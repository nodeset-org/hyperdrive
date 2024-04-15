package validator

import (
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your Constellation validators in NodeSet.",
		Subcommands: []*cli.Command{
			{
				Name:    "claim-rewards",
				Aliases: []string{"cr"},
				Usage:   "Claim rewards for your validators",
				Action: func(c *cli.Context) error {
					// Validate args
					// if err := utils.ValidateArgCount(c, 0); err != nil {
					// 	return err
					// }

					return nil
					// return uploadDepositData(c)
				},
			},
		},
	})
}

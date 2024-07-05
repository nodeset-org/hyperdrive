package nodeset

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your account with the Stakewise vault in NodeSet.",
		Subcommands: []*cli.Command{
			{
				Name:    "upload-deposit-data",
				Aliases: []string{"u"},
				Flags: []cli.Flag{
					utils.YesFlag,
				},
				Usage: "Uploads the combined deposit data for all of your validator keys to NodeSet's Stakewise vault, so they can be assigned new deposits.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return uploadDepositData(c)
				},
			},
			{
				Name:    "generate-deposit-data",
				Aliases: []string{"g"},
				Flags: []cli.Flag{
					generatePubkeyFlag,
					generateIndentFlag,
				},
				Usage: "Generates and prints the deposit data for your validators without uploading it to NodeSet. Useful for debugging.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return generateDepositData(c)
				},
			},
		},
	})
}

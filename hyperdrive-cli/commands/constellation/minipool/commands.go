package minipool

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage your Constellation minipools",
		Subcommands: []*cli.Command{
			{
				Name:    "status",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					statusIncludeFinalizedFlag,
				},
				Usage: "Get the status of the node's minipools.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getMinipoolStatus(c)
				},
			},
			{
				Name:    "create",
				Aliases: []string{"c"},
				Flags: []cli.Flag{
					utils.YesFlag,
					saltFlag,
				},
				Usage: "Create a new minipool.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return createMinipool(c)
				},
			},
			{
				Name:    "stake",
				Aliases: []string{"k"},
				Flags: []cli.Flag{
					utils.YesFlag,
					stakeMinipoolsFlag,
				},
				Usage: "Stake one or minipools that are still in prelaunch but have passed the Rocket Pool scrub check.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return stakeMinipools(c)
				},
			},
			{
				Name:    "upload-signed-exits",
				Aliases: []string{"u"},
				Flags: []cli.Flag{
					utils.YesFlag,
					uploadSignedExitsMinipoolsFlag,
				},
				Usage: "Upload signed exit messages for one or more validators to the NodeSet service if it doesn't already have them.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return uploadSignedExits(c)
				},
			},
			{
				Name:    "exit",
				Aliases: []string{"e"},
				Flags: []cli.Flag{
					utils.YesFlag,
					exitMinipoolsFlag,
					exitVerboseFlag,
					exitManualModeFlag,
					exitManualPubkeyFlag,
					exitManualIndexFlag,
				},
				Usage: "Voluntarily exit one or more minipools from the Beacon Chain, ending validation duties and withdrawing their full balances back to the Execution layer.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return exitMinipools(c)
				},
			},
			{
				Name:    "find-vanity-address",
				Aliases: []string{"v"},
				Usage:   "Search for a custom vanity minipool address",
				Flags: []cli.Flag{
					vanityPrefixFlag,
					vanitySaltFlag,
					vanityThreadsFlag,
					vanityAddressFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return findVanitySalt(c)
				},
			},
		},
	})
}

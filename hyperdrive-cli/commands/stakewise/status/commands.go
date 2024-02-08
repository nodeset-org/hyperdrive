package status

import (
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Get active validators",
		Action: func(c *cli.Context) error {
			return getNodeStatus(c)
		},
	})
}

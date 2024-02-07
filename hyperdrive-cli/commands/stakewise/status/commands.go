package status

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Get SW node status",
		Action: func(c *cli.Context) error {
			fmt.Printf("!!! Get SW node status\n")
			return getNodeStatus(c)
		},
	})
}

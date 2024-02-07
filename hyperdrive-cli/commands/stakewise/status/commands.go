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
		Usage:   "Blah",
		Action: func(c *cli.Context) error {
			fmt.Printf("!!! Get wallet statuses\n")
			return nil
		},
	})
}

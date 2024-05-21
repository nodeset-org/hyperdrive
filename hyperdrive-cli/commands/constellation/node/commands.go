package node

import (
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(cmd *cli.Command, name string, aliases []string) {
	cmd.Subcommands = append(cmd.Subcommands, &cli.Command{
		Name:        name,
		Aliases:     aliases,
		Usage:       "Manage your Constellation node in NodeSet.",
		Subcommands: []*cli.Command{},
	})
}

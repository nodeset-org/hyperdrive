package swcmd

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/constellation/node"

	"github.com/urfave/cli/v2"
)

// Register commands

func RegisterCommands(app *cli.App, name string, aliases []string) {
	cmd := &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage the Constellation module",
	}
	node.RegisterCommands(cmd, "node", []string{"n"})

	app.Commands = append(app.Commands, cmd)
}

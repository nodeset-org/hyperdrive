package constellation

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/constellation/minipool"
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
	minipool.RegisterCommands(cmd, "minipool", []string{"m"})
	node.RegisterCommands(cmd, "node", []string{"n"})

	app.Commands = append(app.Commands, cmd)
}

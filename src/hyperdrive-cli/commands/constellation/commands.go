package constcmd

import (
	"github.com/urfave/cli/v2"
)

// Register commands

func RegisterCommands(app *cli.App, name string, aliases []string) {
	cmd := &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage the Constellation module",
	}
	// TODO: HUY: Add commands here

	app.Commands = append(app.Commands, cmd)
}

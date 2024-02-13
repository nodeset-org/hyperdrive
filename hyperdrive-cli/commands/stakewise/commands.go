package swcmd

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/nodeset"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/status"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/validator"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/wallet"

	"github.com/urfave/cli/v2"
)

// Register commands

func RegisterCommands(app *cli.App, name string, aliases []string) {
	cmd := &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage the Stakewise module",
	}
	nodeset.RegisterCommands(cmd, "nodeset", []string{"ns"})
	wallet.RegisterCommands(cmd, "wallet", []string{"w"})
	status.RegisterCommands(cmd, "status", []string{"s"})
	validator.RegisterCommands(cmd, "validator", []string{"v"})

	app.Commands = append(app.Commands, cmd)
}

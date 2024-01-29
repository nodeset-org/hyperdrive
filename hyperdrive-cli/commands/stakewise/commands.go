package stakewise

import (
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
	wallet.RegisterCommands(cmd, "wallet", []string{"w"})
	app.Commands = append(app.Commands, cmd)
}

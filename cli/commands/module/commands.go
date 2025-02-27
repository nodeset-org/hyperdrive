package module

import (
	"github.com/nodeset-org/hyperdrive/cli/utils"
	"github.com/urfave/cli/v2"
)

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage Hyperdrive modules",
		Flags: []cli.Flag{
			utils.ComposeFileFlag,
		},
		Subcommands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List the modules currently installed on the system",
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return listModules(c)
				},
			}, {
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Install a module from a package file",
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 1)

					// Run command
					return installModule(c, c.Args().Get(0))
				},
			},
		},
	})
}

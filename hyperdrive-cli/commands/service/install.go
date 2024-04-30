package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive-daemon/shared"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

var (
	installVerboseFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"r"},
		Usage:   "Print installation script command output",
	}
	installNoDepsFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "no-deps",
		Aliases: []string{"d"},
		Usage:   "Do not install Operating System dependencies",
	}
	installPathFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "path",
		Aliases: []string{"p"},
		Usage:   "A custom path to install Hyperdrive to",
	}
	installVersionFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "The Hyperdrive package version to install",
		Value:   fmt.Sprintf("v%s", shared.HyperdriveVersion),
	}
	installLocalFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "local-script",
		Aliases: []string{"l"},
		Usage:   fmt.Sprintf("Use a local installer script instead of pulling it down from the source repository. The script and the installer package must be in your current working directory.%sMake sure you absolutely trust the script before using this flag.%s", terminal.ColorRed, terminal.ColorReset),
	}
)

// Install the Hyperdrive service
func installService(c *cli.Context) error {
	// Prompt for confirmation
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm(fmt.Sprintf(
		"Hyperdrive will be installed --Version: %s\n\n%sIf you're upgrading, your existing configuration will be backed up and preserved.\nAll of your previous settings will be migrated automatically.%s\nAre you sure you want to continue?",
		c.String(installVersionFlag.Name), terminal.ColorGreen, terminal.ColorReset,
	))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Install service
	err := hd.InstallService(
		c.Bool(installVersionFlag.Name),
		c.Bool(installNoDepsFlag.Name),
		c.String(installVersionFlag.Name),
		c.String(installPathFlag.Name),
		c.Bool(installLocalFlag.Name),
	)
	if err != nil {
		return err
	}

	// Print success message & return
	fmt.Println("")
	fmt.Println("The Hyperdrive service was successfully installed!")

	printPatchNotes()

	// Reload the config after installation
	_, isNew, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading new configuration: %w", err)
	}

	// Report next steps
	fmt.Printf("%s\n=== Next Steps ===\n", terminal.ColorBlue)
	fmt.Printf("Run 'hyperdrive service config' to review the settings changes for this update, or to continue setting up your node.%s\n", terminal.ColorReset)

	// Print the docker permissions notice
	if isNew {
		fmt.Printf("\n%sNOTE:\nSince this is your first time installing Hyperdrive, please start a new shell session by logging out and back in or restarting the machine.\n", terminal.ColorYellow)
		fmt.Printf("This is necessary for your user account to have permissions to use Docker.%s", terminal.ColorReset)
	}

	return nil

}

// Print the latest patch notes for this release
// TODO: get this from an external source and don't hardcode it into the CLI
func printPatchNotes() {
	fmt.Println()
	fmt.Println(shared.Logo)
	fmt.Printf("%s=== Hyperdrive v%s ===%s\n\n", terminal.ColorGreen, shared.HyperdriveVersion, terminal.ColorReset)
	fmt.Printf("Changes you should be aware of before starting:\n\n")

	fmt.Printf("%s=== Stakewise Module ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("Stakewise functions now initialize the Stakewise wallet if it's missing instead of erroring out.")
}

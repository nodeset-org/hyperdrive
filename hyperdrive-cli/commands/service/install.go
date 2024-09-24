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
	installVersionFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "The Hyperdrive package version to install",
		Value:   fmt.Sprintf("v%s", shared.HyperdriveVersion),
	}
	installLocalScriptFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "local-script",
		Aliases: []string{"ls"},
		Usage:   fmt.Sprintf("Path to a local installer script. If this is specified, Hyperdrive will use it instead of pulling the script down from the source repository. %sMake sure you absolutely trust the script before using this flag.%s", terminal.ColorRed, terminal.ColorReset),
	}
	installLocalPackageFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "local-package",
		Aliases: []string{"lp"},
		Usage:   fmt.Sprintf("Path to a local installer package. If this is specified, Hyperdrive will use it instead of pulling the package down from the source repository. Requires -ls. %sMake sure you absolutely trust the script before using this flag.%s", terminal.ColorRed, terminal.ColorReset),
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

	// Install service
	err := client.InstallService(client.InstallOptions{
		RequireEscalation:       true,
		Verbose:                 c.Bool(installVerboseFlag.Name),
		NoDeps:                  c.Bool(installNoDepsFlag.Name),
		Version:                 c.String(installVersionFlag.Name),
		InstallPath:             "",
		RuntimePath:             "",
		LocalInstallScriptPath:  c.String(installLocalScriptFlag.Name),
		LocalInstallPackagePath: c.String(installLocalPackageFlag.Name),
	})
	if err != nil {
		return err
	}

	// Print success message & return
	fmt.Println("")
	fmt.Println("The Hyperdrive service was successfully installed!")

	printPatchNotes()

	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

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

	fmt.Printf("%s=== Mainnet Support! ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("This version of Hyperdrive supports the Ethereum Mainnet. You can now access the Gravita vault if you're a StakeWise module user and stake ETH on Gravita's behalf.")
	fmt.Println()

	fmt.Printf("%s=== IPv6 Support ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("There's a new toggle in the Hyperdrive section of the settings to enable IPv6 on the Hyperdrive services. Enabling it requires adding support to the Docker service itself first; please read the full patch notes for more information.")
	fmt.Println()

	fmt.Printf("%s=== StakeWise DB ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("The StakeWise Operator database has been moved out of the user data directory and onto a dedicated Docker volume. When you start the StakeWise service for the first time if you were previously running v1.0.0, it will regenerate the database which will cause elevated CPU load for a few hours.")
	fmt.Println()
}

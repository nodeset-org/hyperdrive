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
	installNoRestartFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "no-restart",
		Aliases: []string{"nr"},
		Usage:   "Do not restart Hyperdrive services after installation",
	}
)

// Install the Hyperdrive service
func installService(c *cli.Context) error {
	fmt.Printf("Hyperdrive will be installed --Version: %s\n\n%sIf you're upgrading, your existing configuration will be backed up and preserved.\nAll of your previous settings will be migrated automatically.", c.String(installVersionFlag.Name), terminal.ColorGreen)
	fmt.Println()

	if !c.Bool(installNoRestartFlag.Name) {
		fmt.Print("The service will restart to apply the update (including restarting your clients). If you have doppelganger detection enabled, any active validators will miss the next few attestations while it runs.")
	}
	fmt.Printf("%s\n", terminal.ColorReset)
	fmt.Println()

	// Prompt for confirmation
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Are you ready to continue?")) {
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
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading new configuration: %w", err)
	}

	// Generate the daemon API keys
	err = hd.GenerateDaemonAuthKeys(cfg)
	if err != nil {
		return fmt.Errorf("error generating daemon API keys: %w", err)
	}

	if isNew {
		fmt.Printf("%s\n=== Next Steps ===\n", terminal.ColorBlue)
		fmt.Printf("Run 'hyperdrive service config' to review the settings changes for this update, or to continue setting up your node.%s\n", terminal.ColorReset)

		// Print the docker permissions notice on first install
		fmt.Printf("\n%sNOTE:\nSince this is your first time installing Hyperdrive, please start a new shell session by logging out and back in or restarting the machine.\n", terminal.ColorYellow)
		fmt.Printf("This is necessary for your user account to have permissions to use Docker.%s", terminal.ColorReset)
	} else if !c.Bool(installNoRestartFlag.Name) {
		// Restart services
		fmt.Println("Restarting Hyperdrive services...")
		err = startService(c, true)
		if err != nil {
			return fmt.Errorf("error restarting services: %w", err)
		}
	} else {
		fmt.Println("Remember to run `hyperdrive service start` to update to the new services!")
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

	fmt.Printf("%s=== Constellation Support ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("This version of Hyperdrive supports the long-awaited Constellation module. You can now register as a Constellation operator and create Rocket Pool minipools without needing to bond your own ETH or RPL (aside from a small 12-hour lockup for security). For more info, visit https://docs.nodeset.io/constellation/introduction.")
	fmt.Println()

	fmt.Printf("%s=== File-Based Networks ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("Hyperdrive now supports creating your own custom network definitions in the networks directory (default for Linux: /usr/share/hyperdrive/networks). Add your own files, and they'll become selectable choices in the TUI!")
	fmt.Println()
}

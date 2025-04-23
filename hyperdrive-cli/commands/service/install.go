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
		err = startService(c, StartMode_ForceUpdate)
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

	fmt.Printf("%s=== New StakeWise Module ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("The StakeWise module has been upgraded to support StakeWise's new v3 vaults, which dramatically improve the node operator experience. Deposits now happen automatically; all you need to do is generate keys in advance, and let it do the rest! Take a look at our documentation to get started: https://docs.nodeset.io/stakewise-integration/node-operator-guide-wip")
	fmt.Println()

	fmt.Printf("%s=== Reth Changes ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("Reth will now preserve event logs and transaction receipts by default, which are required for the new StakeWise module. If you previously used Reth without this configuration manually enabled, you will need to resync your node to regenerate the pruned events. If you are using an external Reth client, please ensure that it is configured to preserve event logs from all contracts, not just the deposit contract logs. Also, there is a new configuration parameter called 'State Prune Distance' that lets you fine-tune how many blocks Reth keeps state in its cache for.")
	fmt.Println()

	fmt.Printf("%s=== Pectra on Mainnet ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("The Pectra network upgrade is schedule for Mainnet on epoch 364032 (May 7 2025 - 10:05:11 AM UTC). This version includes clients that support it, except for Nethermind.")
	fmt.Println()

	fmt.Printf("%s=== Holesky Deprecation ===%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Println("The Holesky testnet is no longer supported deprecated. If you were using it, please run `hyperdrive wallet purge` to remove all of your wallet and validator keys, then run`hyperdrive service config` and select Hoodi as your network.")
	fmt.Println()
}

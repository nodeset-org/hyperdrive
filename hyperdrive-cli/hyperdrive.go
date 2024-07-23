package main

import (
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	csconfig "github.com/nodeset-org/hyperdrive-constellation/shared/config"
	"github.com/nodeset-org/hyperdrive-daemon/shared"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/constellation"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/nodeset"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/service"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/wallet"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/context"
	"github.com/urfave/cli/v2"
)

const (
	defaultConfigFolder string      = ".hyperdrive"
	traceMode           os.FileMode = 0644

	// System dir path for Linux
	linuxSystemDir string = "/usr/share/hyperdrive"

	// Subfolders under the system dir
	scriptsDir        string = "scripts"
	templatesDir      string = "templates"
	overrideSourceDir string = "override"
	networksDir       string = "networks"
)

// Flags
var (
	allowRootFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "allow-root",
		Aliases: []string{"r"},
		Usage:   "Allow hyperdrive to be run as the root user",
	}
	configPathFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "config-path",
		Aliases: []string{"c"},
		Usage:   "Directory to install and save all of Hyperdrive's configuration and data to",
	}
	maxFeeFlag *cli.Float64Flag = &cli.Float64Flag{
		Name:    "max-fee",
		Aliases: []string{"f"},
		Usage:   "The max fee (including the priority fee) you want a transaction to cost, in gwei. Use 0 to set it automatically based on network conditions.",
		Value:   0,
	}
	maxPriorityFeeFlag *cli.Float64Flag = &cli.Float64Flag{
		Name:    "max-priority-fee",
		Aliases: []string{"i"},
		Usage:   "The max priority fee you want a transaction to use, in gwei. Use 0 to set it automatically.",
		Value:   0,
	}
	nonceFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:  "nonce",
		Usage: "Use this flag to explicitly specify the nonce that the next transaction should use, so it can override an existing 'stuck' transaction. If running a command that sends multiple transactions, the first will be given this nonce and the rest will be incremented sequentially.",
		Value: 0,
	}
	debugFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:  "debug",
		Usage: "Enable debug printing of API commands",
	}
	secureSessionFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "secure-session",
		Aliases: []string{"s"},
		Usage:   "Some commands may print sensitive information to your terminal. Use this flag when nobody can see your screen to allow sensitive data to be printed without prompting",
	}
	apiAddressFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "api-address",
		Aliases: []string{"a"},
		Usage:   "The address of the Hyperdrive API server to connect to. If left blank it will default to 'localhost' at the port specified in the service configuration.",
	}
	httpTracePathFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "http-trace-path",
		Aliases: []string{"htp"},
		Usage:   "The path to save HTTP trace logs to. Leave blank to disable HTTP tracing",
	}
)

// Run
func main() {
	// Add logo and attribution to application help template
	attribution := "ATTRIBUTION:\n   Adapted from the Rocket Pool Smart Node (https://github.com/rocket-pool/smartnode) with love."
	cli.AppHelpTemplate = fmt.Sprintf("\n%s\n\n%s\n%s\n", shared.Logo, cli.AppHelpTemplate, attribution)
	cli.CommandHelpTemplate = fmt.Sprintf("%s\n%s\n", cli.CommandHelpTemplate, attribution)
	cli.SubcommandHelpTemplate = fmt.Sprintf("%s\n%s\n", cli.SubcommandHelpTemplate, attribution)

	// Initialize application
	app := cli.NewApp()

	// Set application info
	app.Name = "hyperdrive"
	app.Usage = "Hyperdrive CLI for NodeSet Node Operator Management"
	app.Version = shared.HyperdriveVersion
	app.Authors = []*cli.Author{
		{
			Name:  "Nodeset",
			Email: "info@nodeset.io",
		},
	}
	app.Copyright = "(c) 2024 NodeSet LLC"

	// Initialize app metadata
	app.Metadata = make(map[string]interface{})

	// Set application flags
	app.Flags = []cli.Flag{
		allowRootFlag,
		configPathFlag,
		apiAddressFlag,
		maxFeeFlag,
		maxPriorityFeeFlag,
		nonceFlag,
		utils.PrintTxDataFlag,
		utils.SignTxOnlyFlag,
		debugFlag,
		httpTracePathFlag,
		secureSessionFlag,
	}

	// Load the network settings
	installInfo := context.NewInstallationInfo()
	hdNetworkSettings, err := hdconfig.LoadSettingsFiles(installInfo.NetworksDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading network settings from path [%s]: %s", installInfo.NetworksDir, err.Error())
		os.Exit(1)
	}
	swNetSettingsDir := filepath.Join(installInfo.NetworksDir, hdconfig.ModulesName, swconfig.ModuleName)
	swNetworkSettings, err := swconfig.LoadSettingsFiles(swNetSettingsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading network settings from path [%s]: %s", swNetSettingsDir, err.Error())
		os.Exit(1)
	}
	csNetSettingsDir := filepath.Join(installInfo.NetworksDir, hdconfig.ModulesName, csconfig.ModuleName)
	csNetworkSettings, err := csconfig.LoadSettingsFiles(csNetSettingsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading network settings from path [%s]: %s", csNetSettingsDir, err.Error())
		os.Exit(1)
	}

	// Set default paths for flags before parsing the provided values
	setDefaultPaths()

	// Register commands
	constellation.RegisterCommands(app, "constellation", []string{"cs"})
	nodeset.RegisterCommands(app, "nodeset", []string{"ns"})
	service.RegisterCommands(app, "service", []string{"s"}, hdNetworkSettings, swNetworkSettings, csNetworkSettings)
	stakewise.RegisterCommands(app, "stakewise", []string{"sw"})
	wallet.RegisterCommands(app, "wallet", []string{"w"})

	var hdCtx *context.HyperdriveContext
	app.Before = func(c *cli.Context) error {
		// Check user ID
		if os.Getuid() == 0 && !c.Bool(allowRootFlag.Name) {
			fmt.Fprintln(os.Stderr, "hyperdrive should not be run as root. Please try again without 'sudo'.")
			fmt.Fprintf(os.Stderr, "If you want to run hyperdrive as root anyway, use the '--%s' option to override this warning.\n", allowRootFlag.Name)
			os.Exit(1)
		}

		var err error
		hdCtx, err = validateFlags(c)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		hdCtx.InstallationInfo = installInfo
		hdCtx.HyperdriveNetworkSettings = hdNetworkSettings
		hdCtx.StakeWiseNetworkSettings = swNetworkSettings
		hdCtx.ConstellationNetworkSettings = csNetworkSettings
		return nil
	}
	app.After = func(c *cli.Context) error {
		if hdCtx != nil && hdCtx.HttpTraceFile != nil {
			_ = hdCtx.HttpTraceFile.Close()
		}
		return nil
	}

	// Run application
	fmt.Println()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println()
}

// Set the default paths for various flags
func setDefaultPaths() {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Cannot get user's home directory: %s\n", err.Error())
		os.Exit(1)
	}

	// Default config folder path
	defaultConfigPath := filepath.Join(homeDir, defaultConfigFolder)
	configPathFlag.Value = defaultConfigPath
}

// Validate the global flags
func validateFlags(c *cli.Context) (*context.HyperdriveContext, error) {
	hdCtx := &context.HyperdriveContext{
		MaxFee:         c.Float64(maxFeeFlag.Name),
		MaxPriorityFee: c.Float64(maxPriorityFeeFlag.Name),
		DebugEnabled:   c.Bool(debugFlag.Name),
		SecureSession:  c.Bool(secureSessionFlag.Name),
	}

	// If set, validate custom nonce
	hdCtx.Nonce = big.NewInt(0)
	if c.IsSet(nonceFlag.Name) {
		customNonce := c.Uint64(nonceFlag.Name)
		hdCtx.Nonce.SetUint64(customNonce)
	}

	// Make sure the config directory exists
	configPath := c.String(configPathFlag.Name)
	path, err := homedir.Expand(strings.TrimSpace(configPath))
	if err != nil {
		return nil, fmt.Errorf("error expanding config path [%s]: %w", configPath, err)
	}
	hdCtx.ConfigPath = path

	// Get the API URL
	address := c.String(apiAddressFlag.Name)
	if address != "" {
		baseUrl, err := url.Parse(address)
		if err != nil {
			return nil, fmt.Errorf("error parsing API address [%s]: %w", hdCtx.ApiUrl, err)
		}
		hdCtx.ApiUrl = baseUrl.JoinPath(hdconfig.HyperdriveApiClientRoute)
	}

	// Get the HTTP trace flag
	httpTracePath := c.String(httpTracePathFlag.Name)
	if httpTracePath != "" {
		hdCtx.HttpTraceFile, err = os.OpenFile(httpTracePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, traceMode)
		if err != nil {
			return nil, fmt.Errorf("error opening HTTP trace file [%s]: %w", httpTracePath, err)
		}
	}

	// TODO: more here
	context.SetHyperdriveContext(c, hdCtx)
	return hdCtx, nil
}

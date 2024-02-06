package main

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/service"
	swcmd "github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/wallet"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/context"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/urfave/cli/v2"
)

const (
	defaultConfigFolder string = ".hyperdrive"
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
			Name:  "Joe Clapis",
			Email: "joe@nodeset.io",
		},
	}
	app.Copyright = "(c) 2024 NodeSet LLC"

	// Initialize app metadata
	app.Metadata = make(map[string]interface{})

	// Set application flags
	app.Flags = []cli.Flag{
		utils.PrintTxDataFlag,
		allowRootFlag,
		configPathFlag,
		maxFeeFlag,
		maxPriorityFeeFlag,
		nonceFlag,
		debugFlag,
		secureSessionFlag,
	}

	// Set default paths for flags before parsing the provided values
	setDefaultPaths()

	// Register commands
	service.RegisterCommands(app, "service", []string{"s"})
	swcmd.RegisterCommands(app, "stakewise", []string{"sw"})
	wallet.RegisterCommands(app, "wallet", []string{"w"})

	app.Before = func(c *cli.Context) error {
		// Check user ID
		if os.Getuid() == 0 && !c.Bool(allowRootFlag.Name) {
			fmt.Fprintln(os.Stderr, "hyperdrive should not be run as root. Please try again without 'sudo'.")
			fmt.Fprintf(os.Stderr, "If you want to run hyperdrive as root anyway, use the '--%s' option to override this warning.\n", allowRootFlag.Name)
			os.Exit(1)
		}

		err := validateFlags(c)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		return nil
	}

	// Run application
	fmt.Println()
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println()
}

// Set the default paths for various flags
func setDefaultPaths() {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Cannot get user's home directory: %w\n", err)
		os.Exit(1)
	}

	// Default config folder path
	defaultConfigPath := filepath.Join(homeDir, defaultConfigFolder)
	configPathFlag.Value = defaultConfigPath
}

// Validate the global flags
func validateFlags(c *cli.Context) error {
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
		return fmt.Errorf("error expanding config path [%s]: %w", configPath, err)
	}
	hdCtx.ConfigPath = path

	// TODO: more here
	context.SetHyperdriveContext(c, hdCtx)
	return nil
}

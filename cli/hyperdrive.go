package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/nodeset-org/hyperdrive/cli/commands/module"
	"github.com/nodeset-org/hyperdrive/cli/commands/service"
	cliutils "github.com/nodeset-org/hyperdrive/cli/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive/management"
	"github.com/nodeset-org/hyperdrive/shared"
	hdutils "github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/urfave/cli/v2"
)

const (
	// Default user config folder
	DefaultUserFolder string = ".hyperdrive"

	// System dir path for Linux
	LinuxSystemDir string = "/usr/share/hyperdrive"

	// Trace file mode for HTTP tracing
	traceMode os.FileMode = 0644
)

// Flags
var (
	debugFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:  "debug",
		Usage: "Enable debug printing of API commands",
	}
	secureSessionFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "secure-session",
		Aliases: []string{"ss"},
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
	app.Copyright = "(c) 2025 NodeSet LLC"

	// Initialize app metadata
	app.Metadata = make(map[string]interface{})

	// Enable Bash Completion
	app.EnableBashCompletion = true

	// Set application flags
	app.Flags = []cli.Flag{
		cliutils.AllowRootFlag,
		cliutils.UserDirPathFlag,
		cliutils.SystemDirPathFlag,
		apiAddressFlag,
		debugFlag,
		httpTracePathFlag,
		secureSessionFlag,
	}

	// Set default paths for flags before parsing the provided values
	setDefaultPaths()

	// Try to load the installed modules
	var systemDir string
	for i, flag := range os.Args {
		if flag == "--"+cliutils.SystemDirPathFlag.Name || flag == "-"+cliutils.SystemDirPathFlag.Aliases[0] {
			if i+1 >= len(os.Args) {
				fmt.Println("System directory flag was provided without a value")
				os.Exit(1)
			}
			systemDir = os.Args[i+1]
			break
		}
	}
	if systemDir == "" {
		systemDir = cliutils.SystemDirPathFlag.Value
	}
	descriptors, err := hdutils.GetInstalledDescriptors(filepath.Join(systemDir, shared.ModulesDir))
	if err != nil {
		fmt.Println("WARNING: Installed modules could not be loaded:", err)
		fmt.Println("Module commands will be disabled until this is resolved.")
	}
	if len(descriptors) > 0 {
		builder := strings.Builder{}
		for _, desc := range descriptors {
			builder.WriteString("   " + string(desc.Name) + ", " + string(desc.Shortcut) + " - " + "Interact with the " + string(desc.Name) + " module\n")
		}
		cli.AppHelpTemplate = fmt.Sprintf(AppHelpTemplate, builder.String())
	}

	// Register base commands
	module.RegisterCommands(app, "module", []string{"m"})
	service.RegisterCommands(app, "service", []string{"s"})

	// Register module commands
	app.CommandNotFound = HandleCommandNotFound
	/*
		for _, desc := range descriptors {
			app.Commands = append(app.Commands, &cli.Command{
				Name:    string(desc.Name),
				Aliases: []string{string(desc.Shortcut)},
				Usage:   "Interact with the " + string(desc.Name) + " module",
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					fmt.Println("Module command not implemented yet")
					return nil
				},
			})
		}
	*/

	var hdCtx *management.HyperdriveContext
	app.Before = func(c *cli.Context) error {
		// Check user ID
		if os.Getuid() == 0 && !c.Bool(cliutils.AllowRootFlag.Name) {
			fmt.Fprintln(os.Stderr, "hyperdrive should not be run as root. Please try again without 'sudo'.")
			fmt.Fprintf(os.Stderr, "If you want to run hyperdrive as root anyway, use the '--%s' option to override this warning.\n", cliutils.AllowRootFlag.Name)
			os.Exit(1)
		}

		var err error
		hdCtx, err = validateFlags(c)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		return nil
	}
	app.After = func(c *cli.Context) error {
		if hdCtx != nil && hdCtx.HttpTraceFile != nil {
			_ = hdCtx.HttpTraceFile.Close()
		}
		return nil
	}
	app.BashComplete = func(c *cli.Context) {
		// Load the context and flags prior to autocomplete
		err := app.Before(c)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		// Run the default autocomplete
		cli.DefaultAppComplete(c)
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

	// Default user dir path
	defaultUserDirPath := filepath.Join(homeDir, DefaultUserFolder)
	cliutils.UserDirPathFlag.Value = defaultUserDirPath

	// Default system directory path
	switch runtime.GOOS {
	// This is where to add different paths for different OS's like macOS
	default:
		// By default just use the Linux path for now
		cliutils.SystemDirPathFlag.Value = LinuxSystemDir
	}
}

// Validate the global flags
func validateFlags(c *cli.Context) (*management.HyperdriveContext, error) {
	// Expand the user and system paths
	configPath := c.String(cliutils.UserDirPathFlag.Name)
	fullConfigPath, err := homedir.Expand(strings.TrimSpace(configPath))
	if err != nil {
		return nil, fmt.Errorf("error expanding config path \"%s\": %w", configPath, err)
	}
	systemPath := c.String(cliutils.SystemDirPathFlag.Name)
	fullSystemPath, err := homedir.Expand(strings.TrimSpace(systemPath))
	if err != nil {
		return nil, fmt.Errorf("error expanding system path \"%s\": %w", systemPath, err)
	}

	hdCtx := management.NewHyperdriveContext(fullConfigPath, fullSystemPath)
	hdCtx.DebugEnabled = c.Bool(debugFlag.Name)
	hdCtx.SecureSession = c.Bool(secureSessionFlag.Name)

	// Get the API URL
	address := c.String(apiAddressFlag.Name)
	if address != "" {
		baseUrl, err := url.Parse(address)
		if err != nil {
			return nil, fmt.Errorf("error parsing API address \"%s\": %w", hdCtx.ApiUrl, err)
		}
		hdCtx.ApiUrl = baseUrl.JoinPath(hdconfig.HyperdriveApiClientRoute)
	}

	// Get the HTTP trace flag
	httpTracePath := c.String(httpTracePathFlag.Name)
	if httpTracePath != "" {
		hdCtx.HttpTraceFile, err = os.OpenFile(httpTracePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, traceMode)
		if err != nil {
			return nil, fmt.Errorf("error opening HTTP trace file \"%s\": %w", httpTracePath, err)
		}
	}

	// TODO: more here
	management.SetHyperdriveContext(c, hdCtx)
	return hdCtx, nil
}

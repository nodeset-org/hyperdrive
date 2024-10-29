package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/nodeset"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/wallet"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

var (
	ignoreSlashTimerFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:  "ignore-slash-timer",
		Usage: fmt.Sprintf("Bypass the safety timer that forces a delay when switching to a new Beacon Node.\n%sUsing this flag to bypass the slashing timer could result in a *major* loss of ETH! Only use this is if you absolutely understand the risks!%s", terminal.ColorRed, terminal.ColorReset),
	}
	tailFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "tail",
		Aliases: []string{"t"},
		Usage:   "The number of lines to show from the end of the logs (number or \"all\")",
		Value:   "100",
	}
)

// Creates CLI argument flags from the parameters of the configuration struct
// TODO: HEADLESS MODE
/*
func createFlagsFromConfigParams(sectionName string, params []config.IParameter, configFlags []cli.Flag, network config.Network) []cli.Flag {
	for _, param := range params {
		common := param.GetCommon()
		var paramName string
		if sectionName == "" {
			paramName = common.ID
		} else {
			paramName = fmt.Sprintf("%s-%s", sectionName, common.ID)
		}

		usage := common.Description
		options := param.GetOptions()
		if len(options) > 0 {
			optionStrings := []string{}
			for _, option := range options {
				optionStrings = append(optionStrings, fmt.Sprint(option.String()))
			}
			usage = fmt.Sprintf("%s\nOptions: %s\n", common.Description, strings.Join(optionStrings, ", "))
		}

		defaultVal := param.GetDefaultAsAny(network)
		configFlags = append(configFlags, &cli.StringFlag{
			Name:  paramName,
			Usage: usage,
			Value: fmt.Sprint(defaultVal),
		})
	}

	return configFlags
}
*/

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
	configFlags := []cli.Flag{
		configUpdateDefaultsFlag,
	}

	// TODO: HEADLESS MODE
	/*
		cfgTemplate := hdconfig.NewHyperdriveConfig("")
		network := cfgTemplate.Network.Value
		// Root params
		configFlags = createFlagsFromConfigParams("", cfgTemplate.GetParameters(), configFlags, network)

		// Subconfigs
		for sectionName, subconfig := range cfgTemplate.GetSubconfigs() {
			configFlags = createFlagsFromConfigParams(sectionName, subconfig.GetParameters(), configFlags, network)
		}
	*/

	app.Commands = append(app.Commands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage Hyperdrive service",
		Flags: []cli.Flag{
			utils.ComposeFileFlag,
		},
		Subcommands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Install the Hyperdrive service",
				Flags: []cli.Flag{
					utils.YesFlag,
					installVerboseFlag,
					installNoDepsFlag,
					installVersionFlag,
					installLocalScriptFlag,
					installLocalPackageFlag,
					installNoRestartFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return installService(c)
				},
			},

			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Configure the Hyperdrive service",
				Flags:   configFlags,
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return configureService(c)
				},
			},

			{
				Name:    "sync",
				Aliases: []string{"y"},
				Usage:   "Get the sync progress of the Execution and Beacon Nodes",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getSyncProgress(c)
				},
			},

			{
				Name:    "status",
				Aliases: []string{"u"},
				Usage:   "View the Hyperdrive service status",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return serviceStatus(c)
				},
			},

			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Start the Hyperdrive service",
				Flags: []cli.Flag{
					ignoreSlashTimerFlag,
					nodeset.RegisterEmailFlag,
					wallet.PasswordFlag,
					wallet.SavePasswordFlag,
					utils.YesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return startService(c, false)
				},
			},

			{
				Name:    "stop",
				Aliases: []string{"pause", "p"},
				Usage:   "Pause the Hyperdrive service",
				Flags: []cli.Flag{
					utils.YesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return stopService(c)
				},
			},

			{
				Name:    "down",
				Aliases: []string{"d"},
				Usage:   "Stop and delete the Hyperdrive service Docker containers and network. This will leave your private data (such as wallet, config, and validator keys) intact.",
				Flags: []cli.Flag{
					utils.YesFlag,
					includeVolumesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return downService(c)
				},
			},

			{
				Name:      "logs",
				Aliases:   []string{"l"},
				Usage:     "View the Hyperdrive service logs",
				ArgsUsage: "[service names]",
				Flags: []cli.Flag{
					tailFlag,
				},
				Action: func(c *cli.Context) error {
					// Run command
					return serviceLogs(c, c.Args().Slice()...)
				},
			},

			{
				Name:      "daemon-logs",
				Aliases:   []string{"dl"},
				Usage:     "View one or more of the logs from the Hyperdrive daemon, or module daemons",
				ArgsUsage: "[api | tasks | <module log names>]",
				Flags: []cli.Flag{
					tailFlag,
				},
				Action: func(c *cli.Context) error {
					// Run command
					return daemonLogs(c, c.Args().Slice()...)
				},
			},

			{
				Name:    "stats",
				Aliases: []string{"a"},
				Usage:   "DEPRECATED - No longer supported",
				Action: func(c *cli.Context) error {
					// Run command
					return serviceStats()
				},
			},

			{
				Name:  "compose",
				Usage: "View the Hyperdrive service docker compose config",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return serviceCompose(c)
				},
			},

			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "View the Hyperdrive service version information",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return serviceVersion(c)
				},
			},

			/*
				{
					Name:    "prune-ec",
					Aliases: []string{"prune-eth1", "n"},
					Usage:   "Shuts down the main ETH1 client and prunes its database, freeing up disk space, then restarts it when it's done.",
					Action: func(c *cli.Context) error {
						// Validate args
						if err := input.ValidateArgCount(c, 0); err != nil {
							return err
						}

						// Run command
						return pruneExecutionClient(c)
					},
				},
			*/
			{
				Name:    "check-cpu-features",
				Aliases: []string{"ccf"},
				Usage:   "Checks if your CPU supports all of the features required by the \"modern\" version of certain client images. If not, it prints what features are missing.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return checkCpuFeatures()
				},
			},

			{
				Name:  "get-config-yaml",
				Usage: "Generate YAML that shows the current configuration schema, including all of the parameters and their descriptions",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return getConfigYaml(c)
				},
			},
			/*
				{
					Name:      "export-ec-data",
					Aliases:   []string{"export-eth1-data"},
					Usage:     "Exports the execution client (eth1) chain data to an external folder. Use this if you want to back up your chain data before switching execution clients.",
					ArgsUsage: "target-folder",
					Flags: []cli.Flag{
						exportEcDataForceFlag,
						exportEcDataDirtyFlag,
						utils.YesFlag,
					},
					Action: func(c *cli.Context) error {
						// Validate args
						if err := input.ValidateArgCount(c, 1); err != nil {
							return err
						}
						targetDir := c.Args().Get(0)

						// Run command
						return exportEcData(c, targetDir)
					},
				},

				{
					Name:      "import-ec-data",
					Aliases:   []string{"import-eth1-data"},
					Usage:     "Imports execution client (eth1) chain data from an external folder. Use this if you want to restore the data from an execution client that you previously backed up.",
					ArgsUsage: "source-folder",
					Action: func(c *cli.Context) error {
						// Validate args
						if err := input.ValidateArgCount(c, 1); err != nil {
							return err
						}
						sourceDir := c.Args().Get(0)

						// Run command
						return importEcData(c, sourceDir)
					},
				},
			*/
			{
				Name:    "resync-ec",
				Aliases: []string{"resync-eth1"},
				Usage:   fmt.Sprintf("%sDeletes the main Execution client's chain data and resyncs it from scratch. Only use this as a last resort!%s", terminal.ColorRed, terminal.ColorReset),
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return resyncExecutionClient(c)
				},
			},

			{
				Name:    "resync-bn",
				Aliases: []string{"resync-eth2"},
				Usage:   fmt.Sprintf("%sDeletes the Beacon Node's chain data and resyncs it from scratch. Only use this as a last resort!%s", terminal.ColorRed, terminal.ColorReset),
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return resyncBeaconNode(c)
				},
			},

			{
				Name:    "terminate",
				Aliases: []string{"t"},
				Usage:   fmt.Sprintf("%sDeletes all of the Hyperdrive Docker containers and volumes, including your Execution Client and Beacon Node chain data and your Prometheus database (if metrics are enabled). Also removes your entire `.hyperdrive` configuration folder, including your wallet, password, and validator keys. Only use this if you are cleaning up Hyperdrive and want to start over!%s", terminal.ColorRed, terminal.ColorReset),
				Flags: []cli.Flag{
					utils.YesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return terminateService(c)
				},
			},

			// Used by the package installers only
			{
				Name:    "safe-start-after-install",
				Aliases: []string{"ssaf"},
				Usage:   "Install the Hyperdrive service",
				Hidden:  true,
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 1)
					systemDir := c.Args().Get(0)

					// Run command
					safeStartAfterInstall(systemDir)
					return nil
				},
			},
		},
	})
}

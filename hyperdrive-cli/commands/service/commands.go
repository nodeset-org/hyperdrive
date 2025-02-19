package service

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

var (
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
			/*
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
			*/
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
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Start the Hyperdrive service",
				Flags: []cli.Flag{
					/*
						ignoreSlashTimerFlag,
						nodeset.RegisterEmailFlag,
						wallet.PasswordFlag,
						wallet.SavePasswordFlag,
					*/
					utils.YesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run command
					return startService(c)
				},
			},

			/*
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
						BashComplete: func(c *cli.Context) {
							// Run bash completion
							daemonLogs_BashCompletion(c)
						},
					},
			*/
			/*
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
			*/
		},
	})
}

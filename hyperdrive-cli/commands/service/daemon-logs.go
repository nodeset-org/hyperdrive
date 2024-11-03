package service

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// View the Hyperdrive daemon logs
func daemonLogs(c *cli.Context, serviceNames ...string) error {
	lines := c.String(tailFlag.Name)
	lineArg := "--lines="
	if lines == "all" {
		lineArg += "+0"
	} else {
		lineArg += lines
	}

	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading Hyperdrive configuration: %w", err)
	}

	// Get the module log file arg names => log file names
	moduleLogLookup := map[string]string{}
	argNames := []string{"api", "tasks"}
	for _, mod := range cfg.GetAllModuleConfigs() {
		if !mod.IsEnabled() {
			continue
		}
		modName := mod.GetModuleName()
		shortModName := mod.GetShortName()
		logNames := mod.GetLogNames()

		for _, logFileName := range logNames {
			ext := filepath.Ext(logFileName)
			argName := shortModName + "-" + strings.TrimSuffix(logFileName, ext)
			absLogFilePath := cfg.Hyperdrive.GetModuleLogFilePath(modName, logFileName)
			moduleLogLookup[argName] = absLogFilePath
			argNames = append(argNames, argName)
		}
	}

	// Print available options if there are no service names
	if len(serviceNames) == 0 {
		fmt.Println("Available service logs:")
		for _, name := range argNames {
			fmt.Println("\t" + name)
		}
		return nil
	}

	// TODO: Get log paths from service names
	logPaths := []string{}
	for _, service := range serviceNames {
		switch service {
		// Vanilla
		case "api", "a":
			logPaths = append(logPaths, cfg.Hyperdrive.GetApiLogFilePath())
		case "tasks", "t":
			logPaths = append(logPaths, cfg.Hyperdrive.GetTasksLogFilePath())

		// Modules
		default:
			logPath, exists := moduleLogLookup[service]
			if !exists {
				return fmt.Errorf("unknown service name: %s", service)
			}
			logPaths = append(logPaths, logPath)
		}
	}

	// Print service logs
	return hd.PrintDaemonLogs(getComposeFiles(c), lineArg, logPaths...)
}

// Bash completion for the daemon logs command - prints all available log file names based on enabled modules
func daemonLogs_BashCompletion(c *cli.Context) {
	argNames := []string{"api", "tasks"}
	err := utils.BootstrapCliForBashCompletion(c)
	if err != nil {
		for _, name := range argNames {
			fmt.Println(name)
		}
		return
	}

	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return
	}
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return
	}

	// Get the module log file arg names => log file names
	moduleLogLookup := map[string]string{}
	for _, mod := range cfg.GetAllModuleConfigs() {
		if !mod.IsEnabled() {
			continue
		}
		shortModName := mod.GetShortName()
		logNames := mod.GetLogNames()

		for _, logFileName := range logNames {
			ext := filepath.Ext(logFileName)
			argName := shortModName + "-" + strings.TrimSuffix(logFileName, ext)
			moduleLogLookup[argName] = logFileName
			argNames = append(argNames, argName)
		}
	}

	// Print available options
	for _, name := range argNames {
		fmt.Println(name)
	}
}

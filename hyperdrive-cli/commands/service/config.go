package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	tuiconfig "github.com/nodeset-org/hyperdrive/hyperdrive-cli/tui/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

var (
	configUpdateDefaultsFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "update-defaults",
		Aliases: []string{"u"},
		Usage:   "Certain configuration values are reset when Hyperdrive is updated, such as Docker container tags; use this flag to force that reset, even if Hyperdrive hasn't been updated",
	}
)

// Configure the service
func configureService(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Load the modules
	modLoadResults, err := hd.LoadModules()
	if err != nil {
		return fmt.Errorf("error loading modules: %w", err)
	}

	// Load the config, checking to see if it's new (hasn't been installed before)
	hdCfg := hd.GetHyperdriveConfiguration()
	settings, isNew, err := hd.LoadMainSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}

	// Check if this is an update
	oldVersion := strings.TrimPrefix(settings.Version, "v")
	currentVersion := strings.TrimPrefix(shared.HyperdriveVersion, "v")
	isUpdate := c.Bool(configUpdateDefaultsFlag.Name) || (oldVersion != currentVersion)

	// Create default settings for modules that are installed but haven't been configured yet
	waitForOk := false
	for _, result := range modLoadResults {
		if result.LoadError != nil {
			fmt.Printf("Skipping module %s because it failed to load: %s\n", result.Info.Descriptor.GetFullyQualifiedModuleName(), result.LoadError.Error())
			waitForOk = true
			continue
		}

		// Check for an existing config
		fqmn := result.Info.Descriptor.GetFullyQualifiedModuleName()
		_, exists := settings.Modules[fqmn]
		if exists {
			continue
		}

		// Create a new default instance for any missing modules
		info := hdCfg.Modules[fqmn]
		defaultSettings := modconfig.CreateModuleSettings(info.Configuration)
		settings.Modules[fqmn] = &modconfig.ModuleInstance{
			Enabled:  false,
			Version:  info.Descriptor.Version.String(),
			Settings: defaultSettings.SerializeToMap(),
		}
	}
	if waitForOk {
		fmt.Println("The above modules will be disabled until their load errors are resolved.")
		fmt.Println("Press any key to continue.")
		_, _ = fmt.Scanln()
	}

	// For upgrades, move the config to the old one and create a new upgraded copy
	var oldSettings *config.HyperdriveSettings
	if isUpdate {
		oldSettings = settings
		settings = settings.CreateCopy()
		err = hd.UpdateDefaults(settings)
		if err != nil {
			return fmt.Errorf("error updating defaults: %w", err)
		}
	}

	// Save the config and exit in headless mode
	/*
		if c.NumFlags() > 0 {
			return fmt.Errorf("NYI")
			// TODO: HEADLESS MODE
				err := configureHeadless(c, cfg)
				if err != nil {
					return fmt.Errorf("error updating config from provided arguments: %w", err)
				}
				return hd.SaveConfig(cfg)
		}
	*/

	// Run the TUI
	app := tview.NewApplication()
	cfg := hd.GetHyperdriveConfiguration()
	modMgr := hd.GetModuleManager()
	md := tuiconfig.NewMainDisplay(app, modMgr, cfg, oldSettings, settings, isNew, isUpdate)
	err = app.Run()
	if err != nil {
		return err
	}
	if !md.ShouldSave {
		fmt.Println("Your changes have not been saved. Your Hyperdrive configuration is the same as it was before.")
		return nil
	}

	// Save the config
	err = md.UpdateSettingsFromTuiSelections()
	if err != nil {
		return fmt.Errorf("error updating settings from TUI selections: %w", err)
	}
	err = hd.SavePrimarySettings(settings, true)
	if err != nil {
		return fmt.Errorf("error saving config settings: %w", err)
	}
	fmt.Println("Settings saved successfully. Starting project adapters...")

	// Start all of the project module adapters
	modInfos, moduleSettingsMap, err := createModuleSettingsArtifacts(hdCfg, settings)
	if err != nil {
		return fmt.Errorf("error creating module settings: %w", err)
	}
	err = deployModules(modMgr, hd.Context.ModulesDir(), settings, moduleSettingsMap, modInfos)
	if err != nil {
		return fmt.Errorf("error deploying modules: %w", err)
	}
	err = startProjectAdapters(hd.Context.UserDirPath, settings.ProjectName, modInfos)
	if err != nil {
		return fmt.Errorf("error starting project adapters: %w", err)
	}

	// Set the settings for each module
	fmt.Println("Saving settings for each module...")
	for _, info := range modInfos {
		pac, err := modMgr.GetProjectAdapterClient(settings.ProjectName, info.Descriptor.GetFullyQualifiedModuleName())
		if err != nil {
			return fmt.Errorf("error getting project adapter client for module [%s]: %w", info.Descriptor.Name, err)
		}
		err = pac.SetSettings(context.Background(), settings)
		if err != nil {
			return fmt.Errorf("error saving settings for module [%s]: %w", info.Descriptor.Name, err)
		}
	}
	fmt.Println("Module settings saved successfully.")

	// Start the modules
	// TODO: ignore this if there are no changes
	if !utils.Confirm("Would you like to restart the services automatically now to apply the changes?") {
		fmt.Println("Please run `hyperdrive service start` when you are ready to apply the changes.")
		return nil
	}
	err = startModules(modMgr, settings, modInfos)
	if err != nil {
		return fmt.Errorf("error starting modules: %w", err)
	}
	return err

	/*
		// Handle network changes
		prefix := fmt.Sprint(md.PreviousConfig.Hyperdrive.ProjectName.Value)
		if md.ChangeNetworks {
			// Remove the checkpoint sync provider
			md.Config.Hyperdrive.LocalBeaconClient.CheckpointSyncProvider.Value = ""
			err = hd.SaveConfig(md.Config)
			if err != nil {
				return fmt.Errorf("error saving config: %w", err)
			}

			fmt.Printf("%sWARNING: You have requested to change networks.\n\nAll of your existing chain data, your node wallet, and your validator keys will be removed. If you had a Checkpoint Sync URL provided for your Beacon Node, it will be removed and you will need to specify a different one that supports the new network.\n\nPlease confirm you have backed up everything you want to keep, because it will be deleted if you answer `y` to the prompt below.\n\n%s", terminal.ColorYellow, terminal.ColorReset)

			if !utils.Confirm("Would you like Hyperdrive to automatically switch networks for you? This will destroy and rebuild your `data` folder and all of Hyperdrive's Docker containers.") {
				fmt.Println("Please clean up the data folder manually before proceeding.")
				return nil
			}

			err = changeNetworks(c)
			if err != nil {
				fmt.Printf("%s%s%s\nHyperdrive could not automatically change networks for you, so you will have to remove your old data folder manually.\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			}
			return nil
		}
	*/

	// Query for service start if this is a new installation
	/*
		if isNew {
			if !utils.Confirm("Would you like to start the Hyperdrive services automatically now?") {
				fmt.Println("Please run `hyperdrive service start` when you are ready to launch.")
				return nil
			}
			return startService(c, true)
		}

		// Query for service start if this is old and there are containers to change
		if len(md.ContainersToRestart) > 0 {
			fmt.Println("The following containers must be restarted for the changes to take effect:")
			for _, container := range md.ContainersToRestart {
				fmt.Printf("\t%s_%s\n", prefix, container)
			}
			if !utils.Confirm("Would you like to restart them automatically now?") {
				fmt.Println("Please run `hyperdrive service start` when you are ready to apply the changes.")
				return nil
			}

			runningContainers, err := hd.GetRunningContainers(prefix)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: couldn't check running containers: %s\n", err.Error())
				runningContainers = map[string]bool{}
			}
			for _, container := range md.ContainersToRestart {
				fullName := fmt.Sprintf("%s_%s", prefix, container)
				if !runningContainers[fullName] {
					fmt.Printf("%s is not currently running.\n", fullName)
				} else {
					fmt.Printf("Stopping %s... ", fullName)
					err := hd.StopContainer(fullName)
					if err != nil {
						fmt.Println("error!")
						fmt.Fprintf(os.Stderr, "Error stopping container %s: %s\n", fullName, err.Error())
						continue
					}
					fmt.Println("done!")
				}
			}

			fmt.Println()
			fmt.Println("Applying changes and restarting containers...")
			return startService(c, true)
		}
	*/

	return err
}

// TODO: HEADLESS MODE
/*
// Updates a configuration from the provided CLI arguments headlessly
func configureHeadless(c *cli.Context, cfg *hdconfig.HyperdriveConfig) error {
	// Root params
	for _, param := range cfg.GetParameters() {
		err := updateConfigParamFromCliArg(c, "", param, cfg)
		if err != nil {
			return err
		}
	}

	// Subconfigs
	for sectionName, subconfig := range cfg.GetSubconfigs() {
		for _, param := range subconfig.GetParameters() {
			err := updateConfigParamFromCliArg(c, sectionName, param, cfg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Update a config section from the CLI flags
func configureSection(c *cli.Context, section config.IConfigSection) error {
	// Update the parameters
	for _, param := range section.GetParameters() {
		err := updateConfigParamFromCliArg(c, "", param)
		if err != nil {
			return err
		}
	}
}

// Updates a config parameter from a CLI flag
func updateConfigParamFromCliArg(c *cli.Context, sectionName string, param config.IParameter) error {
	var paramName string
	if sectionName == "" {
		paramName = param.GetCommon().ID
	} else {
		paramName = fmt.Sprintf("%s-%s", sectionName, param.GetCommon().ID)
	}

	if c.IsSet(paramName) {

		switch param.Type {
		case config.ParameterType_Bool:
			param.Value = c.Bool(paramName)
		case cfgconfig.ParameterType_Int:
			param.Value = c.Int(paramName)
		case cfgconfig.ParameterType_Float:
			param.Value = c.Float64(paramName)
		case cfgconfig.ParameterType_String:
			setting := c.String(paramName)
			if param.MaxLength > 0 && len(setting) > param.MaxLength {
				return fmt.Errorf("error setting value for %s: [%s] is too long (max length %d)", paramName, setting, param.MaxLength)
			}
			param.Value = c.String(paramName)
		case cfgconfig.ParameterType_Uint:
			param.Value = c.Uint(paramName)
		case cfgconfig.ParameterType_Uint16:
			param.Value = uint16(c.Uint(paramName))
		case cfgconfig.ParameterType_Choice:
			selection := c.String(paramName)
			found := false
			for _, option := range param.Options {
				if fmt.Sprint(option.Value) == selection {
					param.Value = option.Value
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("error setting value for %s: [%s] is not one of the valid options", paramName, selection)
			}
		}
	}

	return nil
}
*/

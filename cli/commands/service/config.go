package service

import (
	"fmt"
	"strings"

	tuiconfig "github.com/nodeset-org/hyperdrive/cli/tui/config"
	"github.com/nodeset-org/hyperdrive/cli/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive/management"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
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
	hd, err := utils.NewHyperdriveManagerFromCtx(c)
	if err != nil {
		return err
	}

	// Load the modules
	err = hd.LoadModules()
	if err != nil {
		return fmt.Errorf("error loading modules: %w", err)
	}

	// Warn about broken modules
	for _, result := range hd.BrokenModules {
		if result.ConfigurationLoadError != nil {
			fmt.Printf("Skipping module %s because it failed to load: %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.ConfigurationLoadError)
		} else if result.GlobalAdapterContainerStatus != management.ContainerStatus_Running {
			fmt.Printf("Skipping module %s because its global adapter container could not start\n", result.Descriptor.GetFullyQualifiedModuleName())
		} else if result.GlobalAdapterRuntimeFileError != nil {
			fmt.Printf("Skipping module %s because its global adapter container file could not be instantiated: %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.GlobalAdapterRuntimeFileError)
		} else if result.DescriptorLoadError != nil {
			fmt.Printf("Skipping module %s because its descriptor could not be loaded: %s\n", result.Descriptor.GetFullyQualifiedModuleName(), result.DescriptorLoadError)
		} else {
			fmt.Printf("Skipping module %s because it could not be loaded for an unknown reason\n", result.Descriptor.GetFullyQualifiedModuleName())
		}
	}
	if len(hd.BrokenModules) > 0 {
		fmt.Println("The above modules will be disabled until their load errors are resolved.")
		if !c.Bool(utils.YesFlag.Name) {
			fmt.Println("Press any key to continue.")
			_, _ = fmt.Scanln()
		}
	}

	// Load the config, checking to see if it's new (hasn't been installed before)
	hdCfg := hd.GetHyperdriveConfiguration()
	currentSettings, isNew, err := hd.LoadMainSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}
	var pendingSettings *hdconfig.HyperdriveSettings
	isUpdate := false
	if isNew {
		// Create default settings for modules that are installed but haven't been configured yet
		for _, result := range hd.HealthyModules {
			// Check for an existing config
			fqmn := result.Descriptor.GetFullyQualifiedModuleName()
			_, exists := currentSettings.Modules[fqmn]
			if exists {
				continue
			}

			// Create a new default instance for any missing modules
			info := hdCfg.Modules[fqmn]
			defaultSettings := modconfig.CreateModuleSettings(info.Configuration)
			currentSettings.Modules[fqmn] = &modconfig.ModuleInstance{
				Enabled:  false,
				Version:  info.Descriptor.Version.String(),
				Settings: defaultSettings.SerializeToMap(),
			}
		}
		pendingSettings = currentSettings.CreateCopy()
	} else {
		pendingSettings, noPendingSettings, err := hd.LoadPendingSettingsFile()
		if err != nil {
			return fmt.Errorf("error loading pending settings: %w", err)
		}
		if noPendingSettings {
			pendingSettings = currentSettings.CreateCopy()
		}

		// Check if this is an update
		oldVersion := strings.TrimPrefix(pendingSettings.Version, "v")
		currentVersion := strings.TrimPrefix(shared.HyperdriveVersion, "v")
		isUpdate = c.Bool(configUpdateDefaultsFlag.Name) || (oldVersion != currentVersion)

		// Create default settings for modules that are installed but haven't been configured yet
		for _, result := range hd.HealthyModules {
			// Check for an existing config
			fqmn := result.Descriptor.GetFullyQualifiedModuleName()
			_, exists := pendingSettings.Modules[fqmn]
			if exists {
				continue
			}

			// Create a new default instance for any missing modules
			info := hdCfg.Modules[fqmn]
			defaultSettings := modconfig.CreateModuleSettings(info.Configuration)
			pendingSettings.Modules[fqmn] = &modconfig.ModuleInstance{
				Enabled:  false,
				Version:  info.Descriptor.Version.String(),
				Settings: defaultSettings.SerializeToMap(),
			}
		}

		// For upgrades, move the config to the old one and create a new upgraded copy
		if isUpdate {
			err = hd.UpdateDefaults(pendingSettings)
			if err != nil {
				return fmt.Errorf("error updating defaults: %w", err)
			}
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
	md := tuiconfig.NewMainDisplay(app, modMgr, cfg, currentSettings, pendingSettings, isNew, isUpdate)
	err = app.Run()
	if err != nil {
		return err
	}
	if !isNew {
		if !md.ShouldSave {
			fmt.Println("Your changes have not been saved. Your Hyperdrive configuration is the same as it was before.")
			return nil
		}
		if !md.HasChanges {
			fmt.Println("No changes were made to the configuration.")
			return nil
		}
	}

	// Save the config
	err = md.UpdateSettingsFromTuiSelections()
	if err != nil {
		return fmt.Errorf("error updating settings from TUI selections: %w", err)
	}
	err = hd.SavePendingSettings(pendingSettings)
	if err != nil {
		return fmt.Errorf("error saving config settings: %w", err)
	}
	fmt.Println("Settings saved successfully and are now 'pending'.")

	// Prompt for service start
	if isNew {
		if !utils.Confirm("Would you like to start the Hyperdrive services automatically now?") {
			fmt.Println("Please run `hyperdrive service start` when you are ready to launch.")
			return nil
		}
	} else if !utils.Confirm("To apply the changes, you must restart some services. Would you like to apply them now?") {
		fmt.Println("Please run `hyperdrive service start` when you are ready to apply the changes.")
		return nil
	}

	// Start the service, stopping any module services that need to be restarted
	return hd.StartService(currentSettings, pendingSettings)
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

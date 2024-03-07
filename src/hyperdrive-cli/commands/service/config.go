package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	cliconfig "github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/service/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

// Configure the service
func configureService(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Make sure the config directory exists first
	err := os.MkdirAll(hd.Context.ConfigPath, 0700)
	if err != nil {
		fmt.Printf("%sYour Hyperdrive user configuration directory of [%s] could not be created:%s.%s\n", terminal.ColorYellow, hd.Context.ConfigPath, err.Error(), terminal.ColorReset)
		return nil
	}

	// Load the config, checking to see if it's new (hasn't been installed before)
	var oldCfg *client.GlobalConfig
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}

	// Check if this is an update
	oldVersion := strings.TrimPrefix(cfg.Hyperdrive.Version, "v")
	currentVersion := strings.TrimPrefix(shared.HyperdriveVersion, "v")
	isUpdate := c.Bool(installUpdateDefaultsFlag.Name) || (oldVersion != currentVersion)

	// For upgrades, move the config to the old one and create a new upgraded copy
	if isUpdate {
		oldCfg = cfg
		cfg = cfg.CreateCopy()
		cfg.UpdateDefaults()
	}

	// Save the config and exit in headless mode
	if c.NumFlags() > 0 {
		return fmt.Errorf("NYI")
		// TODO: HEADLESS MODE
		/*
			err := configureHeadless(c, cfg)
			if err != nil {
				return fmt.Errorf("error updating config from provided arguments: %w", err)
			}
			return hd.SaveConfig(cfg)
		*/
	}

	app := tview.NewApplication()
	md := cliconfig.NewMainDisplay(app, oldCfg, cfg, isNew, isUpdate)
	err = app.Run()
	if err != nil {
		return err
	}

	// Deal with saving the config and printing the changes
	if md.ShouldSave {
		// Save the config
		err = hd.SaveConfig(md.Config)
		if err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}
		fmt.Println("Your changes have been saved!")

		// Handle network changes
		prefix := fmt.Sprint(md.PreviousConfig.Hyperdrive.ProjectName.Value)
		if md.ChangeNetworks {
			// Remove the checkpoint sync provider
			md.Config.Hyperdrive.LocalBeaconConfig.CheckpointSyncProvider.Value = ""
			err = hd.SaveConfig(md.Config)
			if err != nil {
				return fmt.Errorf("error saving config: %w", err)
			}

			fmt.Printf("%sWARNING: You have requested to change networks.\n\nAll of your existing chain data, your node wallet, and your validator keys will be removed. If you had a Checkpoint Sync URL provided for your Beacon Node, it will be removed and you will need to specify a different one that supports the new network.\n\nPlease confirm you have backed up everything you want to keep, because it will be deleted if you answer `y` to the prompt below.\n\n%s", terminal.ColorYellow, terminal.ColorReset)

			if !utils.Confirm("Would you like Hyperdrive to automatically switch networks for you? This will destroy and rebuild your `data` folder and all of Hyperdrive's Docker containers.") {
				fmt.Println("Please clean up the data folder manually before proceeding.")
				return nil
			}

			err = changeNetworks(c, hd)
			if err != nil {
				fmt.Printf("%s%s%s\nHyperdrive could not automatically change networks for you, so you will have to remove your old data folder manually.\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			}
			return nil
		}

		// Query for service start if this is a new installation
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

			fmt.Println()
			for _, container := range md.ContainersToRestart {
				fullName := fmt.Sprintf("%s_%s", prefix, container)
				fmt.Printf("Stopping %s... ", fullName)
				hd.StopContainer(fullName)
				fmt.Print("done!\n")
			}

			fmt.Println()
			fmt.Println("Applying changes and restarting containers...")
			return startService(c, true)
		}
	} else {
		fmt.Println("Your changes have not been saved. Your Hyperdrive configuration is the same as it was before.")
		return nil
	}

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

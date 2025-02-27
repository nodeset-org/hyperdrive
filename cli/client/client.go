package client

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/nodeset-org/hyperdrive/cli/context"
	"github.com/nodeset-org/hyperdrive/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/urfave/cli/v2"
)

type HyperdriveClient struct {
	Context *context.HyperdriveContext

	cfgDir    string
	systemDir string
	modMgr    *utils.ModuleManager
	cfgMgr    *utils.ConfigurationManager
}

// Create a new Hyperdrive client from the CLI context
func NewHyperdriveClientFromCtx(c *cli.Context) (*HyperdriveClient, error) {
	hdCtx := context.GetHyperdriveContext(c)
	if hdCtx == nil {
		return nil, fmt.Errorf("Hyperdrive CLI context has not been created")
	}
	cfgDir := hdCtx.UserDirPath
	systemDir := hdCtx.SystemDirPath

	// Config manager
	cfgMgr := utils.NewConfigurationManager(cfgDir, systemDir)

	// Module manager
	//adapterKeyPath := shared.GetAdapterKeyPath(cfgDir)
	modulesDir := shared.GetModulesDirectoryPath(systemDir)
	gacDir := shared.GetGlobalAdapterDirectoryPath(systemDir)
	modMgr, err := utils.NewModuleManager(modulesDir, gacDir, cfgDir)
	if err != nil {
		return nil, fmt.Errorf("error creating module manager: %w", err)
	}

	return &HyperdriveClient{
		cfgDir:    cfgDir,
		systemDir: systemDir,
		modMgr:    modMgr,
		cfgMgr:    cfgMgr,
		Context:   hdCtx,
	}, nil
}

// Get the current Hyperdrive configuration
func (c *HyperdriveClient) GetHyperdriveConfiguration() *config.HyperdriveConfig {
	return c.cfgMgr.HyperdriveConfiguration
}

// Get the module manager
func (c *HyperdriveClient) GetModuleManager() *utils.ModuleManager {
	return c.modMgr
}

// Load the main (currently applied) settings
func (c *HyperdriveClient) LoadMainSettingsFile() (*config.HyperdriveSettings, bool, error) {
	primaryConfigPath := filepath.Join(c.cfgDir, config.SettingsFilename)
	return c.LoadSettingsFile(primaryConfigPath)
}

// Load the pending (not yet applied) settings
func (c *HyperdriveClient) LoadPendingSettingsFile() (*config.HyperdriveSettings, bool, error) {
	pendingConfigPath := filepath.Join(c.cfgDir, config.PendingSettingsFilename)
	return c.LoadSettingsFile(pendingConfigPath)
}

// Load the config settings
func (c *HyperdriveClient) LoadSettingsFile(path string) (*config.HyperdriveSettings, bool, error) {
	cfg, err := c.cfgMgr.HyperdriveConfiguration.LoadSettingsFromFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("error loading config settings file [%s]: %w", path, err)
	}

	if cfg == nil {
		defaults := modconfig.CreateModuleSettings(c.cfgMgr.HyperdriveConfiguration)
		cfg = config.NewHyperdriveSettings()
		err = defaults.ConvertToKnownType(cfg)
		if err != nil {
			return nil, false, fmt.Errorf("error creating default hyperdrive settings: %w", err)
		}
		return cfg, true, nil
	}
	return cfg, false, nil
}

// Load all of the module info and settings
func (c *HyperdriveClient) LoadModules() ([]*utils.ModuleInfoLoadResult, error) {
	results, err := c.modMgr.LoadModuleInfo(true)
	if err != nil {
		return nil, fmt.Errorf("error loading module info: %w", err)
	}
	for _, result := range results {
		if result.LoadError == nil {
			name := result.Info.Descriptor.GetFullyQualifiedModuleName()
			c.cfgMgr.HyperdriveConfiguration.Modules[name] = result.Info.ModuleInfo
		}
	}
	return results, nil
}

// Initializes a new installation of hyperdrive
func (c *HyperdriveClient) InitializeNewInstallation() error {
	// Make sure the config directory exists
	err := c.createDirectory(c.cfgDir, 0755)
	if err != nil {
		return err
	}

	// Make the subdirectories
	subdirs := []string{
		config.BackupConfigFolder,
		shared.LogsDir,
		shared.SecretsDir,
		shared.OverrideDir,
		shared.RuntimeDir,
		shared.MetricsDir,
	}
	for _, subdir := range subdirs {
		err = c.createDirectory(filepath.Join(c.cfgDir, subdir), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create a directory if it doesn't already exist
func (c HyperdriveClient) createDirectory(path string, mode os.FileMode) error {
	stat, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(path, mode)
		if err != nil {
			return fmt.Errorf("error creating directory [%s]: %w", c.cfgDir, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking directory [%s]: %w", c.cfgDir, err)
	} else if !stat.IsDir() {
		return fmt.Errorf("directory [%s] is already in use but is not a directory", c.cfgDir)
	}

	return nil
}

// Update the defaults for the Hyperdrive settings after upgrading.
// Assumes the modules have been loaded already.
func (c *HyperdriveClient) UpdateDefaults(hdSettings *config.HyperdriveSettings) error {
	hdCfg := c.GetHyperdriveConfiguration()

	// Make an instance
	instance := modconfig.CreateModuleSettings(hdCfg)
	err := instance.CopySettingsFromKnownType(hdSettings)
	if err != nil {
		return fmt.Errorf("error creating settings instance: %w", err)
	}

	// Update the defaults
	err = modconfig.UpdateDefaults(hdCfg, instance)
	if err != nil {
		return fmt.Errorf("error updating defaults: %w", err)
	}

	// Set the settings
	err = instance.ConvertToKnownType(hdSettings)
	if err != nil {
		return fmt.Errorf("error converting settings to known type: %w", err)
	}

	// Update the modules
	for fqmn, module := range hdCfg.Modules {
		// Get the instance
		instance, exists := hdSettings.Modules[fqmn]
		if !exists {
			continue
		}

		// Ignore modules that are up to date
		if instance.Version == module.Descriptor.Version.String() {
			continue
		}

		// Update the defaults
		modCfg := module.Configuration
		modSettings := modconfig.CreateModuleSettings(modCfg)
		err := modconfig.UpdateDefaults(modCfg, modSettings)
		if err != nil {
			return fmt.Errorf("error updating defaults for module [%s]: %w", fqmn, err)
		}
		instance.Settings = modSettings.SerializeToMap()
	}

	return nil
}

// Save the settings to a file, optionally saving the current config as a backup
func (c *HyperdriveClient) SavePendingSettings(settings *config.HyperdriveSettings) error {
	err := c.InitializeNewInstallation()
	if err != nil {
		return fmt.Errorf("error initializing hyperdrive user directory [%s]: %w", c.cfgDir, err)
	}

	// Save the settings
	settingsPath := filepath.Join(c.cfgDir, config.PendingSettingsFilename)
	err = settings.SaveToFile(settingsPath)
	if err != nil {
		return fmt.Errorf("error writing Hyperdrive pending settings file [%s]: %w", settingsPath, err)
	}

	return nil
}

// Commit the pending settings to the primary settings file.
// If there is no pending settings file, this will do nothing.
func (c *HyperdriveClient) CommitPendingSettings(backupOldSettings bool) error {
	// Check if the pending settings file exists
	pendingSettingsPath := filepath.Join(c.cfgDir, config.PendingSettingsFilename)
	_, err := os.Stat(pendingSettingsPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("error checking for Hyperdrive pending settings file [%s]: %w", pendingSettingsPath, err)
	}

	// Backup the old settings if requested
	if backupOldSettings {
		err = c.backupPrimarySettings()
		if err != nil {
			return fmt.Errorf("error backing up Hyperdrive settings file: %w", err)
		}
	}

	// Save the settings
	primaryConfigPath := filepath.Join(c.cfgDir, config.SettingsFilename)
	err = os.Rename(pendingSettingsPath, primaryConfigPath)
	if err != nil {
		return fmt.Errorf("error committing Hyperdrive settings file: %w", err)
	}

	return nil
}

// Backup the primary settings file
func (c *HyperdriveClient) backupPrimarySettings() error {
	// Check if the old settings file exists
	primarySettingsPath := filepath.Join(c.cfgDir, config.SettingsFilename)
	_, err := os.Stat(primarySettingsPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("error checking for Hyperdrive settings file [%s]: %w", primarySettingsPath, err)
	}

	// Backup the old settings file
	now := time.Now().Local()
	backupFilename := fmt.Sprintf(config.BackupConfigFilenamePattern, now.Format("2006-01-02_15-04-05"))
	backupSettingsPath := filepath.Join(c.cfgDir, config.BackupConfigFolder, backupFilename)
	err = os.Rename(primarySettingsPath, backupSettingsPath)
	if err != nil {
		return fmt.Errorf("error backing up Hyperdrive settings file go [%s]: %w", backupSettingsPath, err)
	}
	return nil
}

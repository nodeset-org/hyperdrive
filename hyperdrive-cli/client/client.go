package client

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/context"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

type HyperdriveClient struct {
	Context *context.HyperdriveContext

	cfgDir    string
	systemDir string
	modMgr    *shared.ModuleManager
	cfgMgr    *config.ConfigurationManager
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
	cfgMgr := config.NewConfigurationManager(cfgDir, systemDir)

	// Module manager
	//adapterKeyPath := shared.GetAdapterKeyPath(cfgDir)
	modulesDir := shared.GetModulesDirectoryPath(systemDir)
	modMgr, err := shared.NewModuleManager(modulesDir)
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

// Load the config settings
func (c *HyperdriveClient) LoadMainSettingsFile() (*config.HyperdriveSettings, bool, error) {
	primaryConfigPath := filepath.Join(c.cfgDir, config.SettingsFilename)
	cfg, err := c.cfgMgr.HyperdriveConfiguration.LoadSettingsFromFile(primaryConfigPath)
	if err != nil {
		return nil, false, fmt.Errorf("error loading main config settings file: %w", err)
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
func (c *HyperdriveClient) LoadModules() ([]*shared.ModuleInfoLoadResult, error) {
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
func (c *HyperdriveClient) SavePrimarySettings(settings *config.HyperdriveSettings, backupOldSettings bool) error {
	err := c.InitializeNewInstallation()
	if err != nil {
		return fmt.Errorf("error initializing hyperdrive user directory [%s]: %w", c.cfgDir, err)
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
	err = settings.SaveToFile(primaryConfigPath)
	if err != nil {
		return fmt.Errorf("error writing Hyperdrive settings file at %s: %w", primaryConfigPath, err)
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

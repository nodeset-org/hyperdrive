package management

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	hdconfig "github.com/nodeset-org/hyperdrive/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/utils"
)

type HyperdriveManager struct {
	Context *HyperdriveContext

	cfgDir    string
	systemDir string
	modMgr    *utils.ModuleManager
	cfgMgr    *utils.ConfigurationManager
}

// Create a new Hyperdrive client from the CLI context
func NewHyperdriveManager(hdCtx *HyperdriveContext) (*HyperdriveManager, error) {
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

	return &HyperdriveManager{
		cfgDir:    cfgDir,
		systemDir: systemDir,
		modMgr:    modMgr,
		cfgMgr:    cfgMgr,
		Context:   hdCtx,
	}, nil
}

// Get the current Hyperdrive configuration
func (m *HyperdriveManager) GetHyperdriveConfiguration() *hdconfig.HyperdriveConfig {
	return m.cfgMgr.HyperdriveConfiguration
}

// Get the module manager
func (m *HyperdriveManager) GetModuleManager() *utils.ModuleManager {
	return m.modMgr
}

// Load the main (currently applied) settings
func (m *HyperdriveManager) LoadMainSettingsFile() (*hdconfig.HyperdriveSettings, bool, error) {
	primaryConfigPath := filepath.Join(m.cfgDir, hdconfig.SettingsFilename)
	return m.LoadSettingsFile(primaryConfigPath)
}

// Load the pending (not yet applied) settings
func (m *HyperdriveManager) LoadPendingSettingsFile() (*hdconfig.HyperdriveSettings, bool, error) {
	pendingConfigPath := filepath.Join(m.cfgDir, hdconfig.PendingSettingsFilename)
	return m.LoadSettingsFile(pendingConfigPath)
}

// Load the config settings
func (m *HyperdriveManager) LoadSettingsFile(path string) (*hdconfig.HyperdriveSettings, bool, error) {
	cfg, err := m.cfgMgr.HyperdriveConfiguration.LoadSettingsFromFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("error loading config settings file [%s]: %w", path, err)
	}

	if cfg == nil {
		defaults := modconfig.CreateModuleSettings(m.cfgMgr.HyperdriveConfiguration)
		cfg = hdconfig.NewHyperdriveSettings()
		err = defaults.ConvertToKnownType(cfg)
		if err != nil {
			return nil, false, fmt.Errorf("error creating default hyperdrive settings: %w", err)
		}
		return cfg, true, nil
	}
	return cfg, false, nil
}

// Load all of the module info and settings
func (m *HyperdriveManager) LoadModules() ([]*utils.ModuleInfoLoadResult, error) {
	results, err := m.modMgr.LoadModuleInfo(true)
	if err != nil {
		return nil, fmt.Errorf("error loading module info: %w", err)
	}
	for _, result := range results {
		if result.LoadError == nil {
			name := result.Info.Descriptor.GetFullyQualifiedModuleName()
			m.cfgMgr.HyperdriveConfiguration.Modules[name] = result.Info.ModuleInfo
		}
	}
	return results, nil
}

// Initializes a new installation of hyperdrive
func (m *HyperdriveManager) InitializeNewInstallation() error {
	// Make sure the config directory exists
	err := m.createDirectory(m.cfgDir, 0755)
	if err != nil {
		return err
	}

	// Make the subdirectories
	subdirs := []string{
		hdconfig.BackupConfigFolder,
		shared.LogsDir,
		shared.SecretsDir,
		shared.OverrideDir,
		shared.RuntimeDir,
		shared.MetricsDir,
	}
	for _, subdir := range subdirs {
		err = m.createDirectory(filepath.Join(m.cfgDir, subdir), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create a directory if it doesn't already exist
func (m HyperdriveManager) createDirectory(path string, mode os.FileMode) error {
	stat, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(path, mode)
		if err != nil {
			return fmt.Errorf("error creating directory [%s]: %w", m.cfgDir, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking directory [%s]: %w", m.cfgDir, err)
	} else if !stat.IsDir() {
		return fmt.Errorf("directory [%s] is already in use but is not a directory", m.cfgDir)
	}

	return nil
}

// Update the defaults for the Hyperdrive settings after upgrading.
// Assumes the modules have been loaded already.
func (m *HyperdriveManager) UpdateDefaults(hdSettings *hdconfig.HyperdriveSettings) error {
	hdCfg := m.GetHyperdriveConfiguration()

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
func (m *HyperdriveManager) SavePendingSettings(settings *hdconfig.HyperdriveSettings) error {
	err := m.InitializeNewInstallation()
	if err != nil {
		return fmt.Errorf("error initializing hyperdrive user directory [%s]: %w", m.cfgDir, err)
	}

	// Save the settings
	settingsPath := filepath.Join(m.cfgDir, hdconfig.PendingSettingsFilename)
	err = settings.SaveToFile(settingsPath)
	if err != nil {
		return fmt.Errorf("error writing Hyperdrive pending settings file [%s]: %w", settingsPath, err)
	}

	return nil
}

// Commit the pending settings to the primary settings file.
// If there is no pending settings file, this will do nothing.
func (m *HyperdriveManager) CommitPendingSettings(backupOldSettings bool) error {
	// Check if the pending settings file exists
	pendingSettingsPath := filepath.Join(m.cfgDir, hdconfig.PendingSettingsFilename)
	_, err := os.Stat(pendingSettingsPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("error checking for Hyperdrive pending settings file [%s]: %w", pendingSettingsPath, err)
	}

	// Backup the old settings if requested
	if backupOldSettings {
		err = m.backupPrimarySettings()
		if err != nil {
			return fmt.Errorf("error backing up Hyperdrive settings file: %w", err)
		}
	}

	// Save the settings
	primaryConfigPath := filepath.Join(m.cfgDir, hdconfig.SettingsFilename)
	err = os.Rename(pendingSettingsPath, primaryConfigPath)
	if err != nil {
		return fmt.Errorf("error committing Hyperdrive settings file: %w", err)
	}

	return nil
}

// Backup the primary settings file
func (m *HyperdriveManager) backupPrimarySettings() error {
	// Check if the old settings file exists
	primarySettingsPath := filepath.Join(m.cfgDir, hdconfig.SettingsFilename)
	_, err := os.Stat(primarySettingsPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("error checking for Hyperdrive settings file [%s]: %w", primarySettingsPath, err)
	}

	// Backup the old settings file
	now := time.Now().Local()
	backupFilename := fmt.Sprintf(hdconfig.BackupConfigFilenamePattern, now.Format("2006-01-02_15-04-05"))
	backupSettingsPath := filepath.Join(m.cfgDir, hdconfig.BackupConfigFolder, backupFilename)
	err = os.Rename(primarySettingsPath, backupSettingsPath)
	if err != nil {
		return fmt.Errorf("error backing up Hyperdrive settings file go [%s]: %w", backupSettingsPath, err)
	}
	return nil
}

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
)

type HyperdriveManager struct {
	Context        *HyperdriveContext
	BrokenModules  []*ModuleInstallation
	HealthyModules []*ModuleInstallation

	cfgDir        string
	systemDir     string
	modMgr        *ModuleManager
	cfgMgr        *ConfigurationManager
	modulesLoaded bool
}

// Create a new Hyperdrive client from the CLI context
func NewHyperdriveManager(hdCtx *HyperdriveContext) (*HyperdriveManager, error) {
	cfgDir := hdCtx.UserDirPath
	systemDir := hdCtx.SystemDirPath

	// Config manager
	cfgMgr := NewConfigurationManager(cfgDir, systemDir)

	// Module manager
	//adapterKeyPath := shared.GetAdapterKeyPath(cfgDir)
	modulesDir := shared.GetModulesDirectoryPath(systemDir)
	gacDir := shared.GetGlobalAdapterDirectoryPath(systemDir)
	modMgr, err := NewModuleManager(modulesDir, gacDir, cfgDir)
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
func (m *HyperdriveManager) GetModuleManager() *ModuleManager {
	return m.modMgr
}

// Load the modules installed on the system, sorting them by status
func (m *HyperdriveManager) LoadModules() error {
	brokenModules, eligibleToStartModules, healthyModules, err := m.reloadModuleInstallations()
	if err != nil {
		return err
	}

	// Try to start the modules that are eligible to start but not broken (include the healthy ones too so Docker Compose has a complete picture)
	if len(eligibleToStartModules) > 0 {
		globalAdapterNames := []string{}
		modulesToStart := map[string]*ModuleInstallation{}
		filesToStart := []string{}
		for _, module := range append(eligibleToStartModules, healthyModules...) {
			globalAdapterNames = append(globalAdapterNames, module.GlobalAdapterContainerName)
			modulesToStart[module.GlobalAdapterContainerName] = module
			filesToStart = append(filesToStart, module.GlobalAdapterRuntimeFilePath)
		}
		err = StartProject(shared.GlobalAdapterProjectName, filesToStart)
		if err != nil {
			return fmt.Errorf("error starting global adapters: %w", err)
		}

		// Reload the modules after starting them
		brokenModules, eligibleToStartModules, healthyModules, err = m.reloadModuleInstallations()
		if err != nil {
			return fmt.Errorf("error reloading modules after starting global adapters: %w", err)
		}
	}

	// Any eligible modules that remain need to be considered broken
	for _, module := range eligibleToStartModules {
		brokenModules = append(brokenModules, module)
	}

	// Load the module configs
	m.cfgMgr.HyperdriveConfiguration.Modules = map[string]*modconfig.ModuleInfo{}
	for _, module := range healthyModules {
		m.cfgMgr.LoadModuleConfiguration(m.modMgr, module)

		// Check the config
		if module.ConfigurationLoadError != nil {
			// Modules that couldn't load their config probably need to be reinstalled
			brokenModules = append(brokenModules, module)
			continue
		}
		if module.Configuration == nil {
			// This should never happen
			brokenModules = append(brokenModules, module)
			module.ConfigurationLoadError = fmt.Errorf("module [%s] has a nil configuration", module.Descriptor.GetFullyQualifiedModuleName())
			continue
		}
	}

	// Remove the healthy modules with broken configs
	newHealthyModules := []*ModuleInstallation{}
	for _, module := range healthyModules {
		if module.ConfigurationLoadError != nil {
			continue
		}
		newHealthyModules = append(newHealthyModules, module)
	}

	// Update the module lists
	m.BrokenModules = brokenModules
	m.HealthyModules = healthyModules
	m.modulesLoaded = true
	return nil
}

// Reload the installed modules, sorting them by status
func (m *HyperdriveManager) reloadModuleInstallations() (
	brokenModules []*ModuleInstallation,
	eligibleToStartModules []*ModuleInstallation,
	healthyModules []*ModuleInstallation,
	err error,
) {
	// Load the installed modules and check on their status
	err = m.modMgr.LoadModules()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error loading modules: %w", err)
	}

	// Sort the modules by their status
	brokenModules = []*ModuleInstallation{}
	eligibleToStartModules = []*ModuleInstallation{}
	healthyModules = []*ModuleInstallation{}
	for _, module := range m.modMgr.InstalledModules {
		if module.DescriptorLoadError != nil {
			// Modules that couldn't load their descriptor probably need to be reinstalled
			brokenModules = append(brokenModules, module)
			continue
		}

		if module.GlobalAdapterRuntimeFileError != nil {
			// Modules missing a deployed global adapter file definitely need to be reinstalled
			brokenModules = append(brokenModules, module)
			continue
		}

		switch module.GlobalAdapterContainerStatus {
		case ContainerStatus_Running:
			// Candidate for checking the config, nothing to do
			break

		case ContainerStatus_Stopped:
			// The adapter isn't running but it exists, so it just needs to be started
			eligibleToStartModules = append(eligibleToStartModules, module)
			continue

		case ContainerStatus_Missing:
			// The adapter doesn't exist, which is weird, but since the runtime file exists we can try to just start it
			eligibleToStartModules = append(eligibleToStartModules, module)
			continue
		default:
			// This should never happen
			return nil, nil, nil, fmt.Errorf("module [%s] has an unknown global adapter container status: %s", module.Descriptor.GetFullyQualifiedModuleName(), module.GlobalAdapterContainerStatus)
		}

		// Everything checks out
		healthyModules = append(healthyModules, module)
	}
	return brokenModules, eligibleToStartModules, healthyModules, nil
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
		cfg, err = hdconfig.CreateDefaultHyperdriveSettingsFromConfiguration(m.cfgMgr.HyperdriveConfiguration)
		if err != nil {
			return nil, false, fmt.Errorf("error creating default settings: %w", err)
		}
		return cfg, true, nil
	}
	return cfg, false, nil
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

// Purge all data associated with the Hyperdrive installation, including data saved by modules.
// Services should probably be stopped before this unless you have a good reason to do a "live purge".
// Note this will require root privileges, as that's what the module containers run as when saving files.
func (m *HyperdriveManager) PurgeData(settings *hdconfig.HyperdriveSettings) error {
	err := os.RemoveAll(settings.UserDataPath)
	if err != nil {
		return fmt.Errorf("error deleting data directory [%s]: %w", settings.UserDataPath, err)
	}

	return nil
}

// Delete the user configuration directory.
// This should only be done when the user is uninstalling Hyperdrive.
// Note this *may* require root privileges to run.
func (m *HyperdriveManager) DeleteUserFolder(ctx *HyperdriveContext) error {
	err := os.RemoveAll(ctx.UserDirPath)
	if err != nil {
		return fmt.Errorf("error purging data: %w", err)
	}
	return nil
}

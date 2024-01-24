package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/shared/types"
	"gopkg.in/yaml.v2"
)

const (
	ConfigFilename string = "user-settings.yml"
)

type ConfigManager struct {
	installDir string
}

func NewConfigManager(installDir string) *ConfigManager {
	return &ConfigManager{
		installDir: installDir,
	}
}

func (m *ConfigManager) LoadOrCreateConfig(isDaemon bool) (*HyperdriveConfig, bool, error) {
	// Make sure HD has been installed
	info, err := os.Stat(m.installDir)
	if os.IsNotExist(err) {
		return nil, false, fmt.Errorf("installation directory [%s] doesn't exist", m.installDir)
	}
	if !os.IsExist(err) && err != nil {
		return nil, false, fmt.Errorf("error accessing installation directory [%s]: %w", m.installDir, err)
	}
	if !info.IsDir() {
		return nil, false, fmt.Errorf("installation path [%s] is not a directory", m.installDir)
	}

	// Prep the config file path
	configFilepath := filepath.Join(m.installDir, ConfigFilename)
	info, err = os.Stat(configFilepath)
	if os.IsExist(err) {
		if info.IsDir() {
			return nil, false, fmt.Errorf("config file [%s] is a directory, not a file", configFilepath)
		}

		// Load it
		cfg, err := LoadFromFile(configFilepath)
		return cfg, true, err
	}
	if !os.IsNotExist(err) && err != nil {
		return nil, false, fmt.Errorf("error checking config file [%s]: %w", configFilepath, err)
	}

	// Make the config
	cfg := NewHyperdriveConfig(m.installDir)
	return cfg, false, nil
}

func (m *ConfigManager) SaveConfig(cfg *HyperdriveConfig) error {
	serializedMap := cfg.Serialize()
	bytes, err := yaml.Marshal(serializedMap)
	if err != nil {
		return fmt.Errorf("error creating config: %w", err)
	}

	// Save it
	configFilepath := filepath.Join(m.installDir, ConfigFilename)
	err = os.WriteFile(configFilepath, bytes, 0664)
	if err != nil {
		return fmt.Errorf("error saving config file: %w", err)
	}
	return nil
}

// Get all of the changed settings between an old and new config
func getChangedSettingsMap(oldConfig *HyperdriveConfig, newConfig *HyperdriveConfig) []types.ChangedSetting {
	changedSettings := []types.ChangedSetting{}

	// Root settings
	oldRootParams := oldConfig.GetParameters()
	newRootParams := newConfig.GetParameters()
	changedSettings = getChangedSettings(oldRootParams, newRootParams, newConfig)

	return changedSettings
}

// Get all of the settings that have changed between the given parameter lists.
// Assumes the parameter lists represent identical parameters (e.g. they have the same number of elements and
// each element has the same ID).
func getChangedSettings(oldParams []types.IParameter, newParams []types.IParameter, newConfig *HyperdriveConfig) []types.ChangedSetting {
	changedSettings := []types.ChangedSetting{}

	for i, param := range newParams {
		oldValString := oldParams[i].GetValueAsString()
		newValString := param.GetValueAsString()
		if oldValString != newValString {
			changedSettings = append(changedSettings, types.ChangedSetting{
				Name:               param.GetCommon().Name,
				OldValue:           oldValString,
				NewValue:           newValString,
				AffectedContainers: getAffectedContainers(param, newConfig),
			})
		}
	}

	return changedSettings
}

// Get a list of containers that will be need to be restarted after this change is applied
func getAffectedContainers(param types.IParameter, cfg *HyperdriveConfig) map[types.ContainerID]bool {
	affectedContainers := map[types.ContainerID]bool{}
	for _, container := range param.GetCommon().AffectsContainers {
		affectedContainers[container] = true
	}
	return affectedContainers
}

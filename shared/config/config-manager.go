package config

import (
	"fmt"
	"os"
	"path/filepath"

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

func (m *ConfigManager) LoadOrCreateConfig() (*HyperdriveConfig, bool, error) {
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
		var cfg HyperdriveConfig
		bytes, err := os.ReadFile(configFilepath)
		if err != nil {
			return nil, false, fmt.Errorf("error reading config file [%s]: %w", configFilepath, err)
		}
		err = yaml.Unmarshal(bytes, &cfg)
		if err != nil {
			return nil, false, fmt.Errorf("error deserializing config file [%s]: %w", configFilepath, err)
		}
		return &cfg, true, nil
	}
	if !os.IsNotExist(err) && err != nil {
		return nil, false, fmt.Errorf("error checking config file [%s]: %w", configFilepath, err)
	}

	// Make the config file
	cfg := NewHyperdriveConfig(m.installDir)
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, false, fmt.Errorf("error creating config: %w", err)
	}

	// Save it
	err = os.WriteFile(configFilepath, bytes, 0664)
	if err != nil {
		return nil, false, fmt.Errorf("error saving config file: %w", err)
	}

	return cfg, false, nil
}

package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
)

const (
	// Mode for the Hyperdrive configuration file
	ConfigFileMode os.FileMode = 0644
)

// Result of processing a module's config
type ModuleConfigProcessResult struct {
	// An error that occurred at the system level while trying to process the module config
	ProcessError error

	// A list of errors or issues with the module's config that need to be addressed prior to saving
	Issues []string

	// A list of ports that the module will expose on the host machine
	Ports map[string]uint16
}

// A manager for Hyperdrive's configuration and all of its modules
type ConfigurationManager struct {
	// The configuration metadata for the module
	HyperdriveConfiguration *HyperdriveConfig

	// The name of the Hyperdrive project
	projectName string
}

// Create a new configuration manager
func NewConfigurationManager(hyperdriveDir string, systemDir string) *ConfigurationManager {
	modulesDir := filepath.Join(systemDir, shared.ModulesDir)
	cfg := NewHyperdriveConfig(hyperdriveDir, modulesDir)
	loader := &ConfigurationManager{
		HyperdriveConfiguration: cfg,
	}
	return loader
}

// Process the configuration settings for each module without saving them. This will validate them and collect any extra information about how they will impact the system.
func (m *ConfigurationManager) ProcessModuleSettings(modMgr *shared.ModuleManager, settings *HyperdriveSettings) (map[*config.ModuleInstance]*ModuleConfigProcessResult, error) {
	results := map[*config.ModuleInstance]*ModuleConfigProcessResult{}
	hdConfigMap := settings.SerializeToMap()

	// Make sure all of the module settings have been created
	for fqmn, module := range settings.Modules {
		modInfo := m.HyperdriveConfiguration.Modules[fqmn]
		if modInfo == nil {
			continue
		}
		if !module.Enabled {
			continue
		}

		result := &ModuleConfigProcessResult{}
		results[module] = result

		// Get the adapter client
		client, err := modMgr.GetGlobalAdapterClient(fqmn)
		if err != nil {
			result.ProcessError = err
			continue
		}

		// Process the config
		response, err := client.ProcessSettings(context.Background(), hdConfigMap)
		if err != nil {
			result.ProcessError = err
			continue
		}
		result.Issues = response.Errors
		result.Ports = response.Ports
	}
	return results, nil
}

// Set the configurations for each module. Provide a list of modules you want to set the configuration for here; any modules that are loaded but not in the map will be skipped. If the map has modules that the manager doesn't know about, they will be ignored.
func (m *ConfigurationManager) SetModuleConfigs(modMgr *shared.ModuleManager, settings *HyperdriveSettings) (map[*config.ModuleInstance]error, error) {
	results := map[*config.ModuleInstance]error{}
	for fqmn, module := range settings.Modules {
		modInfo := m.HyperdriveConfiguration.Modules[fqmn]
		if modInfo == nil {
			continue
		}
		if !module.Enabled {
			continue
		}

		// Get the adapter client
		client, err := modMgr.GetGlobalAdapterClient(fqmn)
		if err != nil {
			results[module] = err
			continue
		}

		// Process the config
		err = client.SetSettings(context.Background(), settings.SerializeToMap())
		if err != nil {
			results[module] = err
			continue
		}
	}
	return results, nil
}

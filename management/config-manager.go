package management

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	hdconfig "github.com/nodeset-org/hyperdrive/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
)

const (
	// Mode for the Hyperdrive configuration file
	ConfigFileMode os.FileMode = 0644
)

// Result of processing a module's config settings
type ModuleProcessSettingsResult struct {
	// An error that occurred at the system level while trying to process the module settings
	ProcessError error

	// A list of errors or issues with the module's config that need to be addressed prior to saving
	Issues []string

	// A list of ports that the module will expose on the host machine
	Ports map[string]uint16

	// A list of services that need to be restarted as a result of the new settings
	ServicesToRestart []string
}

// A manager for Hyperdrive's configuration and all of its modules
type ConfigurationManager struct {
	// The configuration metadata for the module
	HyperdriveConfiguration *hdconfig.HyperdriveConfig
}

// Create a new configuration manager
func NewConfigurationManager(hyperdriveDir string, systemDir string) *ConfigurationManager {
	modulesDir := filepath.Join(systemDir, shared.ModulesDir)
	cfg := hdconfig.NewHyperdriveConfig(hyperdriveDir, modulesDir)
	loader := &ConfigurationManager{
		HyperdriveConfiguration: cfg,
	}
	return loader
}

// Load the configuration for a module installation, noting any failures in the installation's config load error
func (m *ConfigurationManager) LoadModuleConfiguration(modMgr *ModuleManager, installInfo *ModuleInstallation) {
	if installInfo.GlobalAdapterContainerStatus != ContainerStatus_Running {
		return
	}
	client, err := modMgr.GetGlobalAdapterClient(installInfo.Descriptor.GetFullyQualifiedModuleName())
	if err != nil {
		installInfo.ConfigurationLoadError = fmt.Errorf("error creating adapter client: %w", err)
		return
	}
	cfg, err := client.GetConfigMetadata(context.Background())
	if err != nil {
		installInfo.ConfigurationLoadError = fmt.Errorf("error getting config metadata: %w", err)
		return
	}
	installInfo.Configuration = cfg
	fqmn := installInfo.Descriptor.GetFullyQualifiedModuleName()
	m.HyperdriveConfiguration.Modules[fqmn] = &modconfig.ModuleInfo{
		Descriptor:    installInfo.Descriptor,
		Configuration: cfg,
	}
}

// Process the configuration settings for each module without saving them. This will validate them and collect any extra information about how they will impact the system.
func (m *ConfigurationManager) ProcessModuleSettings(modMgr *ModuleManager, oldSettings *hdconfig.HyperdriveSettings, newSettings *hdconfig.HyperdriveSettings) (map[*modconfig.ModuleInstance]*ModuleProcessSettingsResult, error) {
	results := map[*modconfig.ModuleInstance]*ModuleProcessSettingsResult{}

	// Make sure all of the module settings have been created
	for fqmn, module := range newSettings.Modules {
		modInfo := m.HyperdriveConfiguration.Modules[fqmn]
		if modInfo == nil {
			continue
		}
		if !module.Enabled {
			continue
		}

		result := &ModuleProcessSettingsResult{}
		results[module] = result

		// Get the adapter client
		client, err := modMgr.GetGlobalAdapterClient(fqmn)
		if err != nil {
			result.ProcessError = err
			continue
		}

		// Process the config
		response, err := client.ProcessSettings(context.Background(), oldSettings, newSettings)
		if err != nil {
			result.ProcessError = err
			continue
		}
		result.Issues = response.Errors
		result.Ports = response.Ports
		result.ServicesToRestart = response.ServicesToRestart
	}
	return results, nil
}

// Set the configurations for each module. Provide a list of modules you want to set the configuration for here; any modules that are loaded but not in the map will be skipped. If the map has modules that the manager doesn't know about, they will be ignored.
func (m *ConfigurationManager) SetModuleConfigs(modMgr *ModuleManager, settings *hdconfig.HyperdriveSettings) (map[*modconfig.ModuleInstance]error, error) {
	results := map[*modconfig.ModuleInstance]error{}
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

		// Set the settings
		err = client.SetSettings(context.Background(), settings)
		if err != nil {
			results[module] = err
			continue
		}
	}
	return results, nil
}

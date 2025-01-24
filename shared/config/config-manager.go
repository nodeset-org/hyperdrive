package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared/adapter"
)

const (
	// Mode for the Hyperdrive configuration file
	ConfigFileMode os.FileMode = 0644
)

// A manager for Hyperdrive's configuration and all of its modules
type ConfigurationManager struct {
	// The configuration metadata for the module
	HyperdriveConfiguration *HyperdriveConfig

	// The name of the Hyperdrive project
	projectName string

	// The path to the directory containing the installed modules
	modulePath string

	// The module adapter key file
	adapterKeyPath string

	// The module adapter key
	adapterKey string

	// A cache of adapter clients for each module
	adapterClients map[string]*adapter.AdapterClient
}

// Create a new configuration manager
func NewConfigurationManager(hyperdriveDir string, systemDir string) *ConfigurationManager {
	modulesDir := filepath.Join(systemDir, ModulesDir)
	cfg := NewHyperdriveConfig(hyperdriveDir, modulesDir)
	loader := &ConfigurationManager{
		HyperdriveConfiguration: cfg,
		modulePath:              modulesDir,
		adapterKeyPath:          cfg.GetAdapterKeyPath(),
		adapterClients:          map[string]*adapter.AdapterClient{},
	}
	return loader
}

// Load the module info and configs based on what's installed. The Hyperdrive project name should be provided in case it changes after the manager has been created.
func (m *ConfigurationManager) LoadModuleInfo(projectName string) ([]*ModuleInfoLoadResult, error) {
	m.projectName = projectName
	err := m.loadAdapterKey()
	if err != nil {
		return nil, err
	}

	// Enumerate the installed modules
	entries, err := os.ReadDir(m.modulePath)
	if err != nil {
		return nil, fmt.Errorf("error reading module directory: %w", err)
	}

	// Go through each module
	moduleInfos := map[string]*ModuleInfo{}
	loadResults := []*ModuleInfoLoadResult{}
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		loadResult := &ModuleInfoLoadResult{
			Name: entry.Name(),
		}
		loadResults = append(loadResults, loadResult)

		// Get the descriptor
		var descriptor modules.ModuleDescriptor
		moduleDir := filepath.Join(m.modulePath, entry.Name())
		descriptorPath := filepath.Join(moduleDir, modules.DescriptorFilename)
		bytes, err := os.ReadFile(descriptorPath)
		if errors.Is(err, fs.ErrNotExist) {
			loadResult.LoadError = ErrNoDescriptor
			continue
		}
		if err != nil {
			loadResult.LoadError = fmt.Errorf("error reading descriptor file: %w", err)
			continue
		}
		err = json.Unmarshal(bytes, &descriptor)
		if err != nil {
			loadResult.LoadError = fmt.Errorf("error unmarshalling descriptor: %w", err)
			continue
		}

		// Get the config
		client, err := m.getAdapterClient(descriptor)
		if err != nil {
			loadResult.LoadError = fmt.Errorf("error creating adapter client: %w", err)
			continue
		}
		cfg, err := client.GetConfigMetadata(context.Background())
		if err != nil {
			loadResult.LoadError = NewModuleInfoLoadError(err)
			continue
		}
		loadResult.Name = descriptor.GetFullyQualifiedModuleName()
		moduleInfo := &ModuleInfo{
			Descriptor:    descriptor,
			Configuration: cfg,
		}
		moduleInfos[descriptor.GetFullyQualifiedModuleName()] = moduleInfo
	}
	m.HyperdriveConfiguration.Modules = moduleInfos
	return loadResults, nil
}

// Process the configurations for each module without saving them. Provide a list of modules you want to set the configuration for here; any modules that are loaded but not in the map will be skipped. If the map has modules that the manager doesn't know about, they will be ignored.
func (m *ConfigurationManager) ProcessModuleConfigurations(hdConfig *HyperdriveConfigInstance) (map[*config.ModuleInstance]*ModuleConfigProcessResult, error) {
	results := map[*config.ModuleInstance]*ModuleConfigProcessResult{}
	err := m.loadAdapterKey()
	if err != nil {
		return nil, err
	}
	hdConfigMap := hdConfig.SerializeToMap()

	// Make sure all of the module settings have been created
	for fqmn, module := range hdConfig.Modules {
		modInfo := m.HyperdriveConfiguration.Modules[fqmn]
		if modInfo == nil {
			continue
		}
		descriptor := modInfo.Descriptor
		if !module.Enabled {
			continue
		}

		result := &ModuleConfigProcessResult{}
		results[module] = result

		// Get the adapter client
		client, err := m.getAdapterClient(descriptor)
		if err != nil {
			result.ProcessError = err
			continue
		}

		// Process the config
		response, err := client.ProcessConfig(context.Background(), hdConfigMap)
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
func (m *ConfigurationManager) SetModuleConfigs(hdConfig *HyperdriveConfigInstance) (map[*config.ModuleInstance]error, error) {
	results := map[*config.ModuleInstance]error{}
	err := m.loadAdapterKey()
	if err != nil {
		return nil, err
	}

	for fqmn, module := range hdConfig.Modules {
		modInfo := m.HyperdriveConfiguration.Modules[fqmn]
		if modInfo == nil {
			continue
		}
		descriptor := modInfo.Descriptor
		if !module.Enabled {
			continue
		}

		// Get the adapter client
		client, err := m.getAdapterClient(descriptor)
		if err != nil {
			results[module] = err
			continue
		}

		// Process the config
		err = client.SetConfig(context.Background(), hdConfig.SerializeToMap())
		if err != nil {
			results[module] = err
			continue
		}
	}
	return results, nil
}

// Load the adapter key file if it's not already loaded
func (m *ConfigurationManager) loadAdapterKey() error {
	if m.adapterKey != "" {
		return nil
	}
	key, err := os.ReadFile(m.adapterKeyPath)
	if err != nil {
		return fmt.Errorf("error reading adapter key file: %w", err)
	}
	m.adapterKey = string(key)
	return nil
}

// Get an adapter client for a module
func (m *ConfigurationManager) getAdapterClient(descriptor modules.ModuleDescriptor) (*adapter.AdapterClient, error) {
	if client, exists := m.adapterClients[descriptor.GetFullyQualifiedModuleName()]; exists {
		return client, nil
	}
	containerName := GetModuleAdapterContainerName(descriptor, m.projectName)
	client, err := adapter.NewAdapterClient(containerName, m.adapterKey)
	if err != nil {
		return nil, err
	}
	m.adapterClients[descriptor.GetFullyQualifiedModuleName()] = client
	return client, nil
}

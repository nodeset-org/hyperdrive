package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"al.essio.dev/pkg/shellescape"
	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/adapter"
	"gopkg.in/yaml.v3"
)

const (
	// Mode for the Hyperdrive configuration file
	ConfigFileMode os.FileMode = 0644
)

// A manager for Hyperdrive's configuration and all of its modules
type ConfigurationManager struct {
	// The configuration metadata for the module
	HyperdriveConfiguration *HyperdriveConfig

	// The configurations for each installed module
	ModuleInfos map[string]*HyperdriveModuleInfo

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
func NewConfigurationManager(hyperdriveDir string, modulesDir string) *ConfigurationManager {
	cfg := NewHyperdriveConfig(hyperdriveDir, modulesDir)
	loader := &ConfigurationManager{
		HyperdriveConfiguration: cfg,
		ModuleInfos:             map[string]*HyperdriveModuleInfo{},
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
	moduleInfos := map[string]*HyperdriveModuleInfo{}
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
		var descriptor modules.HyperdriveModuleDescriptor
		moduleInfo := &HyperdriveModuleInfo{}
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
		moduleInfo.Descriptor = descriptor
		moduleInfo.Configuration = cfg
		moduleInfos[moduleInfo.Descriptor.GetFullyQualifiedModuleName()] = moduleInfo
	}
	m.ModuleInfos = moduleInfos
	m.HyperdriveConfiguration.ModuleInfo = moduleInfos
	return loadResults, nil
}

// Process the configurations for each module without saving them. Provide a list of modules you want to set the configuration for here; any modules that are loaded but not in the map will be skipped. If the map has modules that the manager doesn't know about, they will be ignored.
func (m *ConfigurationManager) ProcessModuleConfigurations(configs []*HyperdriveModuleInstanceInfo) (map[*HyperdriveModuleInstanceInfo]*ModuleConfigProcessResult, error) {
	results := map[*HyperdriveModuleInstanceInfo]*ModuleConfigProcessResult{}
	err := m.loadAdapterKey()
	if err != nil {
		return nil, err
	}

	for _, module := range configs {
		modInfo := module.ModuleInfo
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
		response, err := client.ProcessConfig(context.Background(), &module.Configuration)
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
func (m *ConfigurationManager) SetModuleConfigs(configs map[*HyperdriveModuleInfo]*config.ModuleConfigurationInstance) (map[*HyperdriveModuleInfo]error, error) {
	results := map[*HyperdriveModuleInfo]error{}
	err := m.loadAdapterKey()
	if err != nil {
		return nil, err
	}

	for _, module := range m.ModuleInfos {
		cfg, exists := configs[module]
		if !exists {
			continue
		}

		// Get the adapter client
		client, err := m.getAdapterClient(module.Descriptor)
		if err != nil {
			results[module] = err
			continue
		}

		// Process the config
		err = client.SetConfig(context.Background(), cfg)
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
func (m *ConfigurationManager) getAdapterClient(descriptor modules.HyperdriveModuleDescriptor) (*adapter.AdapterClient, error) {
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

// Load the Hyperdrive configuration from a file; the Hyperdrive user directory will be set to the directory containing the config file
func (m *ConfigurationManager) LoadInstanceFromFile(configFilePath string, systemDir string) (*HyperdriveConfigInstance, error) {
	// Return nil if the file doesn't exist
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	// Read the file
	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read Hyperdrive settings file at %s: %w", shellescape.Quote(configFilePath), err)
	}

	// Attempt to parse it out into a config instance
	var cfg *HyperdriveConfigInstance
	if err := yaml.Unmarshal(configBytes, cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	// Link all of the modules to the module info
	for name, module := range cfg.Modules {
		if info, exists := m.ModuleInfos[name]; exists {
			module.ModuleInfo = info
		}
	}
	return cfg, nil
}

// Save an instance to a file, updating the version to be the current version of Hyperdrive
func (m *ConfigurationManager) SaveInstanceToFile(configFilePath string, instance *HyperdriveConfigInstance) error {
	// Serialize the instance
	instance.Version = shared.HyperdriveVersion

	// Serialize the instance
	configBytes, err := yaml.Marshal(instance)
	if err != nil {
		return fmt.Errorf("could not serialize config instance: %w", err)
	}

	// Write the file
	if err := os.WriteFile(configFilePath, configBytes, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}
	return nil
}

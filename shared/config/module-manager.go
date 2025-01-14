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
	// The name of the descriptor file for a module
	DescriptorFileName string = "descriptor.json"
)

var (
	// The module descriptor was missing
	ErrNoDescriptor error = errors.New("descriptor file is missing")

	// The module config was not loaded because the module adapter container is missing
	ErrNoAdapterContainer error = errors.New("adapter container is missing")

	// The module config was not loaded because the module adapter is not running
	ErrAdapterContainerOffline error = errors.New("adapter container is offline")
)

// An error that occurs when loading a module configuration
type ModuleConfigLoadError struct {
	// The error thrown by the adapter container while getting the module config
	internalError error
}

// Create a new module config load error
func NewModuleConfigLoadError(err error) ModuleConfigLoadError {
	return ModuleConfigLoadError{internalError: err}
}

// Get the error message for a module config load error
func (e ModuleConfigLoadError) Error() string {
	return "error loading module config: " + e.internalError.Error()
}

// The configuration for a module, along with some module metadata
type ModuleConfig struct {
	// The module's descriptor
	Descriptor modules.HyperdriveModuleDescriptor

	// Whether or not the module is currently enabled
	Enabled bool

	// The configuration metadata for the module
	Config config.IConfiguration

	// An error that occurred while loading the module config from its adapter, if any
	LoadError error
}

type ModuleConfigProcessResult struct {
	// An error that occurred at the system level while trying to process the module config
	ProcessError error

	// A list of errors or issues with the module's config that need to be addressed prior to saving
	Issues []string

	// A list of ports that the module will expose on the host machine
	Ports map[string]uint16
}

// A manager for module configurations
type ModuleConfigManager struct {
	ModuleConfigs []*ModuleConfig

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

// Create a new module config manager
func NewModuleConfigManager(modulePath string, adapterKeyPath string) *ModuleConfigManager {
	loader := &ModuleConfigManager{
		ModuleConfigs:  []*ModuleConfig{},
		modulePath:     modulePath,
		adapterKeyPath: adapterKeyPath,
		adapterClients: map[string]*adapter.AdapterClient{},
	}
	return loader
}

// Load the module configs based on what's installed. The Hyperdrive project name should be provided in case it changes after the module manager has been created.
func (m *ModuleConfigManager) LoadModuleConfigs(projectName string) error {
	m.projectName = projectName
	err := m.loadAdapterKey()
	if err != nil {
		return err
	}

	// Enumerate the installed modules
	entries, err := os.ReadDir(m.modulePath)
	if err != nil {
		return fmt.Errorf("error reading module directory: %w", err)
	}

	// Go through each module
	moduleConfigs := []*ModuleConfig{}
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		// Get the descriptor
		var descriptor modules.HyperdriveModuleDescriptor
		moduleConfig := &ModuleConfig{}
		moduleConfigs = append(moduleConfigs, moduleConfig)
		moduleDir := filepath.Join(m.modulePath, entry.Name())
		descriptorPath := filepath.Join(moduleDir, DescriptorFileName)
		bytes, err := os.ReadFile(descriptorPath)
		if errors.Is(err, fs.ErrNotExist) {
			moduleConfig.LoadError = ErrNoDescriptor
			continue
		}
		if err != nil {
			moduleConfig.LoadError = fmt.Errorf("error reading descriptor file: %w", err)
			continue
		}
		err = json.Unmarshal(bytes, &descriptor)
		if err != nil {
			moduleConfig.LoadError = fmt.Errorf("error unmarshalling descriptor: %w", err)
			continue
		}

		// Get the config
		client, err := m.getAdapterClient(descriptor)
		if err != nil {
			moduleConfig.LoadError = fmt.Errorf("error creating adapter client: %w", err)
			continue
		}
		cfg, err := client.GetConfigMetadata(context.Background())
		if err != nil {
			moduleConfig.LoadError = NewModuleConfigLoadError(err)
			continue
		}
		moduleConfig.Descriptor = descriptor
		moduleConfig.Config = cfg
	}
	m.ModuleConfigs = moduleConfigs
	return nil
}

// Process the module configs
func (m *ModuleConfigManager) ProcessModuleConfigs() (map[*ModuleConfig]*ModuleConfigProcessResult, error) {
	results := map[*ModuleConfig]*ModuleConfigProcessResult{}
	err := m.loadAdapterKey()
	if err != nil {
		return nil, err
	}

	for _, module := range m.ModuleConfigs {
		result := &ModuleConfigProcessResult{}
		results[module] = result

		// Get the adapter client
		client, err := m.getAdapterClient(module.Descriptor)
		if err != nil {
			result.ProcessError = err
			continue
		}

		// Process the config
		response, err := client.ProcessConfig(context.Background(), module.Config)
		if err != nil {
			result.ProcessError = err
			continue
		}
		result.Issues = response.Errors
		result.Ports = response.Ports
	}
	return results, nil
}

// Save the module configs
func (m *ModuleConfigManager) SaveModuleConfigs() (map[*ModuleConfig]error, error) {
	results := map[*ModuleConfig]error{}
	err := m.loadAdapterKey()
	if err != nil {
		return nil, err
	}

	for _, module := range m.ModuleConfigs {
		// Get the adapter client
		client, err := m.getAdapterClient(module.Descriptor)
		if err != nil {
			results[module] = err
			continue
		}

		// Process the config
		err = client.SetConfig(context.Background(), module.Config)
		if err != nil {
			results[module] = err
			continue
		}
	}
	return results, nil
}

// Load the adapter key file if it's not already loaded
func (m *ModuleConfigManager) loadAdapterKey() error {
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
func (m *ModuleConfigManager) getAdapterClient(descriptor modules.HyperdriveModuleDescriptor) (*adapter.AdapterClient, error) {
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

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
	ConfigLoadError error
}

// A loader for module configurations
type ModuleConfigLoader struct {
	// The name of the Hyperdrive project
	projectName string

	// The path to the directory containing the installed modules
	modulePath string

	// The module adapter key file
	adapterKey string
}

// Create a new module config loader
func NewModuleConfigLoader(projectName string, modulePath string, adapterKeyPath string) (*ModuleConfigLoader, error) {
	key, err := os.ReadFile(adapterKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading adapter key file: %w", err)
	}

	loader := &ModuleConfigLoader{
		projectName: projectName,
		modulePath:  modulePath,
		adapterKey:  string(key),
	}
	return loader, nil
}

// Load the module configs based on what's installed
func (l *ModuleConfigLoader) LoadModuleConfigs() ([]*ModuleConfig, error) {
	// Enumerate the installed modules
	entries, err := os.ReadDir(l.modulePath)
	if err != nil {
		return nil, fmt.Errorf("error reading module directory: %w", err)
	}

	// Go through each module
	moduleConfigs := []*ModuleConfig{}
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		// Load the module config
		modDir := filepath.Join(l.modulePath, entry.Name())
		moduleConfig := l.loadModuleConfig(modDir)
		moduleConfigs = append(moduleConfigs, moduleConfig)
		if err != nil {
			return nil, fmt.Errorf("error loading module config for %s: %w", entry.Name(), err)
		}
	}
	return moduleConfigs, nil
}

// Load the configuration for a module
func (l *ModuleConfigLoader) loadModuleConfig(moduleDir string) *ModuleConfig {
	// Get the descriptor
	var descriptor modules.HyperdriveModuleDescriptor
	descriptorPath := filepath.Join(moduleDir, DescriptorFileName)
	bytes, err := os.ReadFile(descriptorPath)
	if errors.Is(err, fs.ErrNotExist) {
		return &ModuleConfig{
			ConfigLoadError: ErrNoDescriptor,
		}
	}
	if err != nil {
		return &ModuleConfig{
			ConfigLoadError: fmt.Errorf("error reading descriptor file: %w", err),
		}
	}
	err = json.Unmarshal(bytes, &descriptor)
	if err != nil {
		return &ModuleConfig{
			ConfigLoadError: fmt.Errorf("error unmarshalling descriptor: %w", err),
		}
	}

	// Create the adapter client
	containerName := GetModuleAdapterContainerName(descriptor, l.projectName)
	client, err := adapter.NewAdapterClient(containerName, l.adapterKey)
	if err != nil {
		return &ModuleConfig{
			ConfigLoadError: fmt.Errorf("error creating adapter client: %w", err),
		}
	}

	// Get the config
	cfg, err := client.GetConfigMetadata(context.Background())
	if err != nil {
		return &ModuleConfig{
			ConfigLoadError: NewModuleConfigLoadError(err),
		}
	}
	return &ModuleConfig{
		Descriptor:      descriptor,
		Config:          cfg,
		ConfigLoadError: nil,
	}
}

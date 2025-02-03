package shared

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared/adapter"
	"github.com/nodeset-org/hyperdrive/shared/auth"
)

const (
	// 384 bit keys by default
	AdapterKeySizeBytes int = 48

	// Mode for the adapter key file
	AdapterKeyMode os.FileMode = 0600
)

// Result of loading a module's info
type ModuleInfoLoadResult struct {
	// An error that occurred while loading the module's info
	LoadError error

	// The module's info
	Info *ModuleInstallationInfo
}

type ProjectInfo struct {
	// The name of the Hyperdrive project
	projectName string

	// The module adapter key file
	adapterKeyPath string

	// The module adapter key
	adapterKey string
}

// A manager for loading and interacting with Hyperdrive's modules, including their adapters
type ModuleManager struct {
	// The path to the directory containing the installed modules
	modulePath string

	// Docker client for interacting with the Docker API
	docker client.APIClient

	// A cache of adapter clients for each module
	adapterClients map[string]*adapter.AdapterClient

	// A cache of loaded module info
	moduleInfos map[string]*config.ModuleInfo
}

// Info about the module, including its installation
type ModuleInstallationInfo struct {
	*config.ModuleInfo

	// The full path of the directory the module is installed in
	DirectoryPath string

	// The name of the global adapter container
	GlobalAdapterContainerName string
}

// Create a new module manager
func NewModuleManager(modulesDir string) (*ModuleManager, error) {
	docker, err := client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %w", err)
	}

	loader := &ModuleManager{
		modulePath:     modulesDir,
		docker:         docker,
		adapterClients: map[string]*adapter.AdapterClient{},
		moduleInfos:    map[string]*config.ModuleInfo{},
	}
	return loader, nil
}

// Load the module descriptors to determine what's installed
func (m *ModuleManager) LoadModuleInfo(ensureModuleStart bool) ([]*ModuleInfoLoadResult, error) {
	// Enumerate the installed modules
	entries, err := os.ReadDir(m.modulePath)
	if err != nil {
		return nil, fmt.Errorf("error reading module directory: %w", err)
	}

	// Find the modules
	loadResults := []*ModuleInfoLoadResult{}
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		moduleDir := filepath.Join(m.modulePath, entry.Name())
		info := &ModuleInstallationInfo{
			DirectoryPath: moduleDir,
			ModuleInfo:    &config.ModuleInfo{},
		}
		loadResult := &ModuleInfoLoadResult{
			Info: info,
		}
		loadResults = append(loadResults, loadResult)

		// Check if the descriptor exists - this is the key for modules
		var descriptor modules.ModuleDescriptor
		descriptorPath := filepath.Join(moduleDir, modules.DescriptorFilename)
		bytes, err := os.ReadFile(descriptorPath)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			loadResult.LoadError = fmt.Errorf("error reading descriptor file [%s]: %w", descriptorPath, err)
			continue
		}

		// Load the descriptor
		err = json.Unmarshal(bytes, &descriptor)
		if err != nil {
			loadResult.LoadError = fmt.Errorf("error unmarshalling descriptor: %w", err)
			continue
		}
		info.Descriptor = descriptor
		info.GlobalAdapterContainerName = GetGlobalAdapterContainerName(descriptor)
	}

	if ensureModuleStart {
		// Start all of the adapters for modules with descriptors
		adapterComposefiles := []string{}
		for _, result := range loadResults {
			if result.LoadError != nil {
				continue
			}
			adapterFile := filepath.Join(result.Info.DirectoryPath, modules.AdapterComposeFilename)

			// Check if the file exists
			stat, err := os.Stat(adapterFile)
			if errors.Is(err, fs.ErrNotExist) {
				result.LoadError = ErrNoAdapterComposeFile
				continue
			}
			if err != nil {
				result.LoadError = fmt.Errorf("error checking adapter compose file [%s]: %w", adapterFile, err)
				continue
			}
			if stat.IsDir() {
				result.LoadError = fmt.Errorf("adapter compose file [%s] is a directory, not a file", adapterFile)
				continue
			}

			// Add the file to the list
			adapterComposefiles = append(adapterComposefiles, adapterFile)
		}
		if len(adapterComposefiles) == 0 {
			return loadResults, nil
		}

		// Run the docker compose command
		cmd := exec.Command("docker-compose", "-f", strings.Join(adapterComposefiles, " "), "up", "-d", "--quiet-pull")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("error starting global adapters: %w", err)
		}
	}

	// Check to see if any adapters aren't started yet
	containers, err := m.docker.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %w", err)
	}
	for _, result := range loadResults {
		if result.LoadError != nil {
			continue
		}

		// Check if the container exists
		id := m.getContainerID(result.Info.GlobalAdapterContainerName, containers)
		if id == "" {
			result.LoadError = ErrNoAdapterContainer
		}

		// Get the container info
		containerInfo, err := m.docker.ContainerInspect(context.Background(), id)
		if err != nil {
			result.LoadError = fmt.Errorf("error inspecting global adapter container [%s]: %w", result.Info.GlobalAdapterContainerName, err)
			continue
		}

		// Check if it's started
		if !containerInfo.State.Running {
			result.LoadError = ErrAdapterContainerOffline
			continue
		}
	}

	// Get the configs from each adapter
	for _, result := range loadResults {
		if result.LoadError != nil {
			continue
		}

		// Get the adapter client
		client, err := m.getGlobalAdapterClientFromDescriptor(result.Info.Descriptor)
		if err != nil {
			result.LoadError = fmt.Errorf("error creating adapter client: %w", err)
			continue
		}

		// Get the config
		cfg, err := client.GetConfigMetadata(context.Background())
		if err != nil {
			result.LoadError = NewModuleInfoLoadError(err)
			continue
		}
		result.Info.Configuration = cfg
	}

	// Set the module info in the cache for loaded modules
	for _, result := range loadResults {
		if result.LoadError == nil {
			m.moduleInfos[result.Info.Descriptor.GetFullyQualifiedModuleName()] = result.Info.ModuleInfo
		}
	}
	return loadResults, nil
}

// Get an adapter client for the global adapter container of a module
func (m *ModuleManager) GetGlobalAdapterClient(fqmn string) (*adapter.AdapterClient, error) {
	if client, exists := m.adapterClients[fqmn]; exists {
		return client, nil
	}
	modInfo, exists := m.moduleInfos[fqmn]
	if !exists {
		return nil, fmt.Errorf("module info not found for %s", fqmn)
	}
	containerName := GetGlobalAdapterContainerName(modInfo.Descriptor)
	client, err := adapter.NewAdapterClient(containerName, "")
	if err != nil {
		return nil, err
	}
	m.adapterClients[fqmn] = client
	return client, nil
}

// Get an adapter client for the global adapter container of a module
func (m *ModuleManager) getGlobalAdapterClientFromDescriptor(descriptor modules.ModuleDescriptor) (*adapter.AdapterClient, error) {
	fqmn := descriptor.GetFullyQualifiedModuleName()
	if client, exists := m.adapterClients[fqmn]; exists {
		return client, nil
	}
	containerName := GetGlobalAdapterContainerName(descriptor)
	client, err := adapter.NewAdapterClient(containerName, "")
	if err != nil {
		return nil, err
	}
	m.adapterClients[fqmn] = client
	return client, nil
}

// TEMPORARY placeholder function that exists solely to facilitate development
func GetGlobalAdapterContainerName(descriptor modules.ModuleDescriptor) string {
	return "hd_" + string(descriptor.Shortcut) + "_adapter"
}

// TEMPORARY placeholder function that exists solely to facilitate development
func GetProjectAdapterContainerName(descriptor *modules.ModuleDescriptor, projectName string) string {
	return "hd_" + projectName + "_" + string(descriptor.Shortcut) + "_adapter"
}

// Load the adapter key file if it's not already loaded
func (m *ModuleManager) loadAdapterKey(project *ProjectInfo) error {
	if project.adapterKey != "" {
		return nil
	}
	key, err := os.ReadFile(project.adapterKeyPath)
	if errors.Is(err, fs.ErrNotExist) {
		return project.createAdapterKey()
	}
	if err != nil {
		return fmt.Errorf("error reading adapter key file: %w", err)
	}
	project.adapterKey = string(key)
	return nil
}

// Creates a new secret adapter key and saves it to the module's adapter key file path if one doesn't already exist
func (p *ProjectInfo) createAdapterKey() error {
	// Check if the key file already exists
	if _, err := os.Stat(p.adapterKeyPath); !errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	// Generate a random key
	key, err := auth.GenerateAuthKey(AdapterKeySizeBytes)
	if err != nil {
		return fmt.Errorf("error generating adapter key: %w", err)
	}

	// Save the key to the file
	err = os.WriteFile(p.adapterKeyPath, []byte(key), AdapterKeyMode)
	if err != nil {
		return fmt.Errorf("error saving adapter key to [%s]: %w", p.adapterKeyPath, err)
	}
	p.adapterKey = key
	return nil
}

// Gets the ID of the container with the given name
func (m *ModuleManager) getContainerID(name string, containers []types.Container) string {
	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == "/"+name {
				return container.ID
			}
		}
	}
	return ""
}

// ==============
// === Errors ===
// ==============

var (
	// The module descriptor was missing
	ErrNoDescriptor error = errors.New("descriptor file is missing")

	// The module config was not loaded because the module adapter container is missing
	ErrNoAdapterContainer error = errors.New("adapter container is missing")

	// The module config was not loaded because the module adapter compose file is missing
	ErrNoAdapterComposeFile error = errors.New("adapter compose file is missing")

	// The module config was not loaded because the module adapter is not running
	ErrAdapterContainerOffline error = errors.New("adapter container is offline")
)

// An error that occurs when loading module information
type ModuleInfoLoadError struct {
	// The error thrown by the adapter container while getting the module config
	internalError error
}

// Create a new module info load error
func NewModuleInfoLoadError(err error) ModuleInfoLoadError {
	return ModuleInfoLoadError{internalError: err}
}

// Get the error message for a module config load error
func (e ModuleInfoLoadError) Error() string {
	return "error loading module info: " + e.internalError.Error()
}

// An error that occurs when loading a module descriptor
type ModuleDescriptorLoadError struct {
	// The error thrown while getting the descriptor
	internalError error
}

// Create a new module descriptor load error
func NewModuleDescriptorLoadError(err error) ModuleDescriptorLoadError {
	return ModuleDescriptorLoadError{internalError: err}
}

// Get the error message for a module descriptor load error
func (e ModuleDescriptorLoadError) Error() string {
	return "error loading module descriptor: " + e.internalError.Error()
}

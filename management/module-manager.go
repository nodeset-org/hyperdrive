package management

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/goccy/go-json"
	"github.com/klauspost/compress/zip"
	hdconfig "github.com/nodeset-org/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/adapter"
	"github.com/nodeset-org/hyperdrive/shared/auth"
	"github.com/nodeset-org/hyperdrive/shared/logging"
	"github.com/nodeset-org/hyperdrive/shared/templates"
	"github.com/nodeset-org/hyperdrive/shared/utils"
)

const (
	// 384 bit keys by default
	AdapterKeySizeBytes int = 48

	// Mode for the adapter key file
	AdapterKeyMode os.FileMode = 0600

	// Mode for module directories
	ModuleDirMode os.FileMode = 0755

	// Mode for module files
	ModuleFileMode os.FileMode = 0644
)

// Result of loading a module's info
type ModuleInfoLoadResult struct {
	// An error that occurred while loading the module's info
	LoadError error

	// The module's info
	Info *ModuleInstallation
}

// A manager for loading and interacting with Hyperdrive's modules, including their adapters
type ModuleManager struct {
	// Information about the modules that have been installed on the system. Run LoadModuleInfo() to populate this.
	InstalledModules []*ModuleInstallation

	// === Internal fields ===

	// Logger for recording messages
	logger *slog.Logger

	// The path to the directory containing the installed modules
	modulePath string

	// The path to the directory for global adapter compose files
	globalAdapterDir string

	// The path to the user's project directory
	userDir string

	// Docker client for interacting with the Docker API
	docker client.APIClient

	// A cache of global adapter clients for each module
	globalAdapterClients map[string]*adapter.AdapterClient

	// A cache of adapter clients for each project
	projectAdapterClients map[string]map[string]*adapter.AdapterClient
}

// Create a new module manager
func NewModuleManager(modulesDir string, globalAdapterDir string, userDir string) (*ModuleManager, error) {
	docker, err := client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %w", err)
	}

	loader := &ModuleManager{
		modulePath:            modulesDir,
		globalAdapterDir:      globalAdapterDir,
		userDir:               userDir,
		docker:                docker,
		globalAdapterClients:  map[string]*adapter.AdapterClient{},
		projectAdapterClients: map[string]map[string]*adapter.AdapterClient{},
	}
	return loader, nil
}

// The system directory for module installation
func (m *ModuleManager) GetModuleSystemDir() string {
	return m.modulePath
}

// Load information about all of the modules installed on the system.
func (m *ModuleManager) LoadModules() error {
	// Enumerate the installed modules
	entries, err := os.ReadDir(m.modulePath)
	if err != nil {
		return fmt.Errorf("error reading module directory \"%s\": %w", m.modulePath, err)
	}

	// Go through the entries in the module directory
	installationInfos := []*ModuleInstallation{}
	for _, entry := range entries {
		logging.SafeDebug(m.logger, "Found candidate in module directory", "dir", entry.Name())
		// Skip non-directories
		if !entry.IsDir() {
			logging.SafeDebug(m.logger, "Entry is not a directory, skipping")
			continue
		}
		moduleDir := filepath.Join(m.modulePath, entry.Name())
		info, err := m.beginLoadModule(moduleDir)
		if err != nil {
			return fmt.Errorf("error loading module \"%s\": %w", moduleDir, err)
		}
		if info != nil {
			installationInfos = append(installationInfos, info)
		}
	}

	// Check the global adapter status for each module
	containers, err := m.docker.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("error listing docker containers: %w", err)
	}
	for _, installInfo := range installationInfos {
		containerName := installInfo.GlobalAdapterContainerName
		id := m.getContainerID(containerName, containers)
		if id == "" {
			logging.SafeDebug(m.logger, "Global adapter container not found", "name", containerName)
			installInfo.GlobalAdapterContainerStatus = ContainerStatus_Missing
		} else {
			containerInfo, err := m.docker.ContainerInspect(context.Background(), id)
			if err != nil {
				return fmt.Errorf("error inspecting container \"%s\": %w", containerName, err)
			}
			if containerInfo.State.Running {
				logging.SafeDebug(m.logger, "Global adapter container is running", "name", containerName)
				installInfo.GlobalAdapterContainerStatus = ContainerStatus_Running
			} else {
				logging.SafeDebug(m.logger, "Global adapter container is stopped", "name", containerName, "status", containerInfo.State.Status)
				installInfo.GlobalAdapterContainerStatus = ContainerStatus_Stopped
			}
		}
	}

	m.InstalledModules = installationInfos
	return nil
}

// Starts the process of loading single module that's been installed to the provided path.
// This checks the descriptor and runtime global adapter compose file.
func (m *ModuleManager) beginLoadModule(moduleDir string) (*ModuleInstallation, error) {
	installInfo := &ModuleInstallation{
		InstallationPath: moduleDir,
	}

	// Check if the descriptor exists - this is the key for modules
	descriptorPath := filepath.Join(moduleDir, modules.DescriptorFilename)
	installInfo.DescriptorPath = descriptorPath
	var descriptor modules.ModuleDescriptor
	bytes, err := os.ReadFile(descriptorPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logging.SafeDebug(m.logger, "Descriptor file not found, skipping")
			return nil, nil
		}
		logging.SafeDebug(m.logger, "Error reading descriptor file", "path", descriptorPath, "error", err)
		installInfo.DescriptorLoadError = fmt.Errorf("error reading descriptor file \"%s\": %w", descriptorPath, err)
		return installInfo, nil
	}

	// Load the descriptor
	err = json.Unmarshal(bytes, &descriptor)
	if err != nil {
		logging.SafeDebug(m.logger, "Error unmarshalling descriptor", "path", descriptorPath, "error", err)
		installInfo.DescriptorLoadError = fmt.Errorf("error unmarshalling descriptor: %w", err)
		return installInfo, nil
	}
	installInfo.Descriptor = &descriptor

	// Check if the global adapter compose file exists
	globalAdapterPath := m.GetGlobalAdapterComposeFilePath(&descriptor)
	installInfo.GlobalAdapterRuntimeFilePath = globalAdapterPath
	fileInfo, err := os.Stat(globalAdapterPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logging.SafeDebug(m.logger, "Global adapter file not found", "path", globalAdapterPath)
			installInfo.GlobalAdapterRuntimeFileError = fs.ErrNotExist
		} else {
			logging.SafeDebug(m.logger, "Error checking global adapter file", "path", globalAdapterPath, "error", err)
			installInfo.GlobalAdapterRuntimeFileError = err
		}
	}

	// Check if the compose file is a regular file
	if installInfo.GlobalAdapterRuntimeFileError == nil {
		if fileInfo.IsDir() {
			logging.SafeDebug(m.logger, "Global adapter file is a directory", "path", globalAdapterPath)
			installInfo.GlobalAdapterRuntimeFileError = fmt.Errorf("global adapter file \"%s\" is a directory, not a file", globalAdapterPath)
		} else if !fileInfo.Mode().IsRegular() {
			logging.SafeDebug(m.logger, "Global adapter file is not a regular file", "path", globalAdapterPath)
			installInfo.GlobalAdapterRuntimeFileError = fmt.Errorf("global adapter file \"%s\" is not a regular file", globalAdapterPath)
		}
	}

	// Try to parse the compose file
	if installInfo.GlobalAdapterRuntimeFileError == nil {
		_, err = ParseComposeFile(shared.GlobalAdapterProjectName, globalAdapterPath)
		if err != nil {
			logging.SafeDebug(m.logger, "Error parsing global adapter file", "path", globalAdapterPath, "error", err)
			installInfo.GlobalAdapterRuntimeFileError = fmt.Errorf("error parsing global adapter file \"%s\": %w", globalAdapterPath, err)
		}
	}

	// Get the global adapter container name
	containerName := utils.GetGlobalAdapterContainerName(&descriptor)
	installInfo.GlobalAdapterContainerName = containerName
	return installInfo, nil
}

// Start all of the global adapters that can be started for the loaded modules.
func (m *ModuleManager) StartGlobalAdapters() error {
	gacFiles := []string{}
	for _, info := range m.InstalledModules {
		if info.DescriptorLoadError != nil {
			m.logger.Debug("Can't start global adapter for module because the descriptor wasn't loaded", "path", filepath.Base(info.InstallationPath))
			continue
		}
		if info.GlobalAdapterRuntimeFileError != nil {
			m.logger.Debug("Can't start global adapter for module because the global adapter file encountered errors", "module", info.Descriptor.GetFullyQualifiedModuleName())
			continue
		}
		gacFiles = append(gacFiles, info.GlobalAdapterRuntimeFilePath)
	}

	if len(gacFiles) == 0 {
		m.logger.Debug("No global adapters to start")
		return nil
	}
	err := StartProject(shared.GlobalAdapterProjectName, gacFiles)
	if err != nil {
		return fmt.Errorf("error starting global adapters: %w", err)
	}
	return nil
}

// Stop all of the global adapters
func (m *ModuleManager) StopGlobalAdapters() error {
	err := StopProject(shared.GlobalAdapterProjectName, nil)
	if err != nil {
		return fmt.Errorf("error stopping global adapters: %w", err)
	}
	return nil
}

// Stop and remove all of the global adapters
func (m *ModuleManager) DownGlobalAdapters() error {
	err := DownProject(shared.GlobalAdapterProjectName, true)
	if err != nil {
		return fmt.Errorf("error removing global adapters: %w", err)
	}
	return nil
}

// Get an adapter client for the global adapter container of a module
func (m *ModuleManager) GetGlobalAdapterClient(fqmn string) (*adapter.AdapterClient, error) {
	if client, exists := m.globalAdapterClients[fqmn]; exists {
		isStale, err := client.CheckIfStale()
		if err != nil {
			return nil, fmt.Errorf("error checking if adapter client is stale: %w", err)
		}
		if !isStale {
			return client, nil
		}
		// If it's stale, delete the old one and regenerate
		delete(m.globalAdapterClients, fqmn)
	}
	var installInfo *ModuleInstallation
	for _, modInfo := range m.InstalledModules {
		if modInfo.Descriptor.GetFullyQualifiedModuleName() == fqmn {
			installInfo = modInfo
			break
		}
	}
	if installInfo == nil {
		return nil, fmt.Errorf("module info not found for %s", fqmn)
	}
	containerName := utils.GetGlobalAdapterContainerName(installInfo.Descriptor)
	client, err := adapter.NewAdapterClient(containerName, "")
	if err != nil {
		return nil, err
	}
	m.globalAdapterClients[fqmn] = client
	return client, nil
}

// Get and adapter client for the project adapter container of a module
func (m *ModuleManager) GetProjectAdapterClient(projectName string, fqmn string) (*adapter.AdapterClient, error) {
	// Check if the cache exists
	adapterClients, exists := m.projectAdapterClients[projectName]
	if !exists {
		adapterClients = map[string]*adapter.AdapterClient{}
	}
	if client, exists := adapterClients[fqmn]; exists {
		isStale, err := client.CheckIfStale()
		if err != nil {
			return nil, fmt.Errorf("error checking if adapter client is stale: %w", err)
		}
		if !isStale {
			return client, nil
		}
		// If it's stale, delete the old one and regenerate
		delete(m.projectAdapterClients, projectName)
	}
	var installInfo *ModuleInstallation
	for _, modInfo := range m.InstalledModules {
		if modInfo.Descriptor.GetFullyQualifiedModuleName() == fqmn {
			installInfo = modInfo
			break
		}
	}
	if installInfo == nil {
		return nil, fmt.Errorf("module info not found for %s", fqmn)
	}

	// Load the adapter key
	moduleDir := filepath.Join(m.userDir, shared.ModulesDir, string(installInfo.Descriptor.Name))
	adapterKeyPath := filepath.Join(moduleDir, shared.SecretsDir, shared.AdapterKeyFile)
	bytes, err := os.ReadFile(adapterKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading adapter key file \"%s\": %w", adapterKeyPath, err)
	}

	// Create the adapter client
	containerName := utils.GetProjectAdapterContainerName(projectName, installInfo.Descriptor)
	client, err := adapter.NewAdapterClient(containerName, string(bytes))
	if err != nil {
		return nil, err
	}
	adapterClients[fqmn] = client
	m.projectAdapterClients[projectName] = adapterClients
	return client, nil
}

// Installs a module to the system directory
func (m *ModuleManager) InstallModule(
	moduleFile string,
) error {
	// Unpack the module
	pkgReader, err := zip.OpenReader(moduleFile)
	if err != nil {
		return fmt.Errorf("error opening module package \"%s\": %w", moduleFile, err)
	}
	defer pkgReader.Close()

	// Get the descriptor
	descriptorBuffer := new(bytes.Buffer)
	descriptorFile, err := pkgReader.Open(modules.DescriptorFilename)
	if err != nil {
		return fmt.Errorf("error opening descriptor file in module package \"%s\": %w", moduleFile, err)
	}
	defer descriptorFile.Close()
	_, err = descriptorBuffer.ReadFrom(descriptorFile)
	if err != nil {
		return fmt.Errorf("error reading descriptor file in module package \"%s\": %w", moduleFile, err)
	}
	var descriptor modules.ModuleDescriptor
	err = json.Unmarshal(descriptorBuffer.Bytes(), &descriptor)
	if err != nil {
		return fmt.Errorf("error unmarshalling descriptor file in module package \"%s\": %w", moduleFile, err)
	}

	// Create the module dir
	moduleDir := filepath.Join(m.modulePath, string(descriptor.Name))
	err = os.MkdirAll(moduleDir, ModuleDirMode)
	if err != nil {
		return fmt.Errorf("error creating module directory \"%s\": %w", moduleDir, err)
	}

	// Extract the files
	descriptorPath := filepath.Join(moduleDir, modules.DescriptorFilename)
	err = os.WriteFile(descriptorPath, descriptorBuffer.Bytes(), ModuleFileMode)
	if err != nil {
		return fmt.Errorf("error writing descriptor file \"%s\": %w", descriptorPath, err)
	}
	for _, file := range pkgReader.File {
		// Skip the descriptor since we already wrote it
		if file.Name == modules.DescriptorFilename {
			continue
		}

		// Handle directories
		if file.FileInfo().IsDir() {
			fileDir := filepath.Join(moduleDir, file.Name)
			err = os.MkdirAll(fileDir, ModuleDirMode)
			if err != nil {
				return fmt.Errorf("error creating directory \"%s\": %w", fileDir, err)
			}
			continue
		}

		// Open the file
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening file \"%s\" in module package \"%s\": %w", file.Name, moduleFile, err)
		}
		defer fileReader.Close()

		// Write the file
		filePath := filepath.Join(moduleDir, file.Name)
		buffer := new(bytes.Buffer)
		_, err = buffer.ReadFrom(fileReader)
		if err != nil {
			return fmt.Errorf("error reading file \"%s\" in module package \"%s\": %w", file.Name, moduleFile, err)
		}
		err = os.WriteFile(filePath, buffer.Bytes(), ModuleFileMode)
		if err != nil {
			return fmt.Errorf("error writing file \"%s\": %w", filePath, err)
		}

		// Instantiate the adapter template in global mode
		if file.Name == modules.AdapterComposeTemplateFilename {
			globalAdapterPath := m.GetGlobalAdapterComposeFilePath(&descriptor)
			globalAdapterDir := filepath.Dir(globalAdapterPath)
			err := os.MkdirAll(globalAdapterDir, ModuleDirMode)
			if err != nil {
				return fmt.Errorf("error creating global adapter directory \"%s\": %w", globalAdapterDir, err)
			}

			adapterTemplate, err := template.New("adapter").Parse(buffer.String())
			if err != nil {
				return fmt.Errorf("error parsing adapter template \"%s\": %w", filePath, err)
			}
			adapterFile, err := os.OpenFile(globalAdapterPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, ModuleFileMode)
			if err != nil {
				return fmt.Errorf("error creating adapter compose file \"%s\": %w", globalAdapterPath, err)
			}
			adapterSrc := templates.NewGlobalAdapterDataSource(&descriptor)
			err = adapterTemplate.Execute(adapterFile, adapterSrc)
			if err != nil {
				return fmt.Errorf("error executing adapter template: %w", err)
			}
		}
	}

	return nil
}

// Deploy a module to the user's directory
func (m *ModuleManager) DeployModule(
	moduleInstallDir string,
	hdSettings *hdconfig.HyperdriveSettings,
	moduleSettingsMap map[string]*modconfig.ModuleSettings,
	info *modconfig.ModuleInfo,
) error {
	hdProjectName := hdSettings.ProjectName

	// Create the adapter data source for project mode
	adapterSrc := templates.NewProjectAdapterDataSource(m.userDir, hdProjectName, info.Descriptor)

	// Build the module directory structure
	err := errors.Join(
		os.MkdirAll(adapterSrc.ModuleConfigDir, ModuleDirMode),
		os.MkdirAll(adapterSrc.ModuleLogDir, ModuleDirMode),
		os.MkdirAll(adapterSrc.ModuleComposeDir, ModuleDirMode),
		os.MkdirAll(adapterSrc.ModuleOverrideDir, ModuleDirMode),
		os.MkdirAll(adapterSrc.ModuleMetricsDir, ModuleDirMode),
		os.MkdirAll(filepath.Dir(adapterSrc.AdapterKeyFile), auth.KeyDirPermissions),
	)
	if err != nil {
		return fmt.Errorf("error creating module directories: %w", err)
	}

	// Build the module data dir
	moduleDataDir := filepath.Join(hdSettings.UserDataPath, shared.ModulesDir, string(info.Descriptor.Name))
	err = os.MkdirAll(moduleDataDir, ModuleDirMode)
	if err != nil {
		return fmt.Errorf("error creating module data directory \"%s\": %w", moduleDataDir, err)
	}

	// Create the adapter key
	adapterKeyPath := filepath.Join(adapterSrc.ModuleDir, shared.SecretsDir, shared.AdapterKeyFile)
	_, err = auth.CreateOrLoadKeyFile(adapterKeyPath, AdapterKeySizeBytes)
	if err != nil {
		return fmt.Errorf("error creating adapter key file \"%s\": %w", adapterKeyPath, err)
	}

	// Instantiate the adapter template
	adapterTemplatePath := filepath.Join(moduleInstallDir, string(info.Descriptor.Name), modules.AdapterComposeTemplateFilename)
	adapterTemplate, err := template.ParseFiles(adapterTemplatePath)
	if err != nil {
		return fmt.Errorf("error parsing adapter template \"%s\": %w", adapterTemplatePath, err)
	}
	adapterRuntimePath := filepath.Join(adapterSrc.ModuleComposeDir, modules.AdapterComposeFilename)
	adapterFile, err := os.OpenFile(adapterRuntimePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, ModuleFileMode)
	if err != nil {
		return fmt.Errorf("error creating adapter compose file \"%s\": %w", adapterRuntimePath, err)
	}
	err = adapterTemplate.Execute(adapterFile, adapterSrc)
	if err != nil {
		return fmt.Errorf("error executing adapter template: %w", err)
	}

	// Instantiate the service templates
	hdDynamicSettings := hdSettings.CreateModuleSettings()
	serviceSrc := templates.NewServiceDataSource(hdDynamicSettings, moduleSettingsMap, info, adapterSrc)
	serviceTemplatePath := filepath.Join(moduleInstallDir, string(info.Descriptor.Name), shared.TemplatesDir)
	entries, err := os.ReadDir(serviceTemplatePath)
	if err != nil {
		return fmt.Errorf("error reading service template directory \"%s\": %w", serviceTemplatePath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		entryName := entry.Name()
		templateSuffix := ".tmpl"
		if filepath.Ext(entryName) != templateSuffix {
			// Skip non-template files
			continue
		}
		templatePath := filepath.Join(serviceTemplatePath, entryName)
		template, err := template.ParseFiles(templatePath)
		if err != nil {
			return fmt.Errorf("error parsing service template \"%s\": %w", templatePath, err)
		}
		serviceName := strings.TrimSuffix(filepath.Base(entryName), templateSuffix)
		serviceName += ".yml"
		serviceRuntimePath := filepath.Join(adapterSrc.ModuleComposeDir, serviceName)
		serviceFile, err := os.OpenFile(serviceRuntimePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, ModuleFileMode)
		if err != nil {
			return fmt.Errorf("error creating service compose file \"%s\": %w", serviceRuntimePath, err)
		}
		err = template.Execute(serviceFile, serviceSrc)
		if err != nil {
			return fmt.Errorf("error executing service template: %w", err)
		}
	}
	return nil
}

// Get the path to the instantiated global adapter compose file for a module
func (m *ModuleManager) GetGlobalAdapterComposeFilePath(descriptor *modules.ModuleDescriptor) string {
	return filepath.Join(m.globalAdapterDir, string(descriptor.Name), modules.AdapterComposeFilename)
}

// Gets the ID of the container with the given name.
// Returns an empty string if the container doesn't exist.
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

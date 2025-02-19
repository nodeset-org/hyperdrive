package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/goccy/go-json"
	"github.com/klauspost/compress/zip"
	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/adapter"
	"github.com/nodeset-org/hyperdrive/shared/auth"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/templates"
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

	// A cache of loaded module info
	moduleInfos map[string]*modconfig.ModuleInfo
}

// Info about the module, including its installation
type ModuleInstallationInfo struct {
	*modconfig.ModuleInfo

	// The full path of the directory the module is installed in
	DirectoryPath string

	// The name of the global adapter container
	GlobalAdapterContainerName string
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
		moduleInfos:           map[string]*modconfig.ModuleInfo{},
	}
	return loader, nil
}

// The system directory for module installation
func (m *ModuleManager) GetModuleSystemDir() string {
	return m.modulePath
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
			ModuleInfo:    &modconfig.ModuleInfo{},
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
			adapterFile := filepath.Join(m.globalAdapterDir, string(result.Info.Descriptor.Name), modules.AdapterComposeFilename)

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
		args := []string{
			"compose",
		}
		for _, file := range adapterComposefiles {
			args = append(args, "-f", file)
		}
		args = append(args,
			"up",
			"-d",
			"--quiet-pull",
		)
		cmd := exec.Command("docker", args...)
		cmd.Env = append(cmd.Env, "COMPOSE_PROJECT_NAME="+shared.GlobalAdapterProjectName)
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
	if client, exists := m.globalAdapterClients[fqmn]; exists {
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
		return client, nil
	}
	modInfo, exists := m.moduleInfos[fqmn]
	if !exists {
		return nil, fmt.Errorf("module info not found for %s", fqmn)
	}

	// Load the adapter key
	moduleDir := filepath.Join(m.userDir, shared.ModulesDir, string(modInfo.Descriptor.Name))
	adapterKeyPath := filepath.Join(moduleDir, shared.SecretsDir, shared.AdapterKeyFile)
	bytes, err := os.ReadFile(adapterKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading adapter key file [%s]: %w", adapterKeyPath, err)
	}

	// Create the adapter client
	containerName := GetProjectAdapterContainerName(&modInfo.Descriptor, projectName)
	client, err := adapter.NewAdapterClient(containerName, string(bytes))
	if err != nil {
		return nil, err
	}
	adapterClients[fqmn] = client
	m.projectAdapterClients[projectName] = adapterClients
	return client, nil
}

// Get an adapter client for the global adapter container of a module
func (m *ModuleManager) getGlobalAdapterClientFromDescriptor(descriptor modules.ModuleDescriptor) (*adapter.AdapterClient, error) {
	fqmn := descriptor.GetFullyQualifiedModuleName()
	if client, exists := m.globalAdapterClients[fqmn]; exists {
		return client, nil
	}
	containerName := GetGlobalAdapterContainerName(descriptor)
	client, err := adapter.NewAdapterClient(containerName, "")
	if err != nil {
		return nil, err
	}
	m.globalAdapterClients[fqmn] = client
	return client, nil
}

// TEMPORARY placeholder function that exists solely to facilitate development
func GetGlobalAdapterContainerName(descriptor modules.ModuleDescriptor) string {
	return shared.GlobalAdapterProjectName + "-" + string(descriptor.Shortcut) + "_adapter"
}

// TEMPORARY placeholder function that exists solely to facilitate development
func GetProjectAdapterContainerName(descriptor *modules.ModuleDescriptor, projectName string) string {
	return projectName + "-" + string(descriptor.Shortcut) + "_adapter"
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

// Installs a module to the system directory
func (m *ModuleManager) InstallModule(
	moduleFile string,
) error {
	// Unpack the module
	pkgReader, err := zip.OpenReader(moduleFile)
	if err != nil {
		return fmt.Errorf("error opening module package [%s]: %w", moduleFile, err)
	}
	defer pkgReader.Close()

	// Get the descriptor
	descriptorBuffer := new(bytes.Buffer)
	descriptorFile, err := pkgReader.Open(modules.DescriptorFilename)
	if err != nil {
		return fmt.Errorf("error opening descriptor file in module package [%s]: %w", moduleFile, err)
	}
	defer descriptorFile.Close()
	_, err = descriptorBuffer.ReadFrom(descriptorFile)
	if err != nil {
		return fmt.Errorf("error reading descriptor file in module package [%s]: %w", moduleFile, err)
	}
	var descriptor modules.ModuleDescriptor
	err = json.Unmarshal(descriptorBuffer.Bytes(), &descriptor)
	if err != nil {
		return fmt.Errorf("error unmarshalling descriptor file in module package [%s]: %w", moduleFile, err)
	}

	// Create the module dir
	moduleDir := filepath.Join(m.modulePath, string(descriptor.Name))
	err = os.MkdirAll(moduleDir, ModuleDirMode)
	if err != nil {
		return fmt.Errorf("error creating module directory [%s]: %w", moduleDir, err)
	}

	// Extract the files
	descriptorPath := filepath.Join(moduleDir, modules.DescriptorFilename)
	err = os.WriteFile(descriptorPath, descriptorBuffer.Bytes(), ModuleFileMode)
	if err != nil {
		return fmt.Errorf("error writing descriptor file [%s]: %w", descriptorPath, err)
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
				return fmt.Errorf("error creating directory [%s]: %w", fileDir, err)
			}
			continue
		}

		// Open the file
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening file [%s] in module package [%s]: %w", file.Name, moduleFile, err)
		}
		defer fileReader.Close()

		// Write the file
		filePath := filepath.Join(moduleDir, file.Name)
		buffer := new(bytes.Buffer)
		_, err = buffer.ReadFrom(fileReader)
		if err != nil {
			return fmt.Errorf("error reading file [%s] in module package [%s]: %w", file.Name, moduleFile, err)
		}
		err = os.WriteFile(filePath, buffer.Bytes(), ModuleFileMode)
		if err != nil {
			return fmt.Errorf("error writing file [%s]: %w", filePath, err)
		}

		// Instantiate the adapter template in global mode
		if file.Name == modules.AdapterComposeTemplateFilename {
			globalAdapterDir := filepath.Join(m.globalAdapterDir, string(descriptor.Name))
			globalAdapterPath := filepath.Join(globalAdapterDir, modules.AdapterComposeFilename)
			err := os.MkdirAll(globalAdapterDir, ModuleDirMode)
			if err != nil {
				return fmt.Errorf("error creating global adapter directory [%s]: %w", globalAdapterDir, err)
			}

			adapterTemplate, err := template.New("adapter").Parse(buffer.String())
			if err != nil {
				return fmt.Errorf("error parsing adapter template [%s]: %w", filePath, err)
			}
			adapterFile, err := os.OpenFile(globalAdapterPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, ModuleFileMode)
			if err != nil {
				return fmt.Errorf("error creating adapter compose file [%s]: %w", globalAdapterPath, err)
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
	adapterSrc := templates.NewProjectAdapterDataSource(m.userDir, hdProjectName, &info.Descriptor)

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
		return fmt.Errorf("error creating module data directory [%s]: %w", moduleDataDir, err)
	}

	// Create the adapter key
	adapterKeyPath := filepath.Join(adapterSrc.ModuleDir, shared.SecretsDir, shared.AdapterKeyFile)
	err = auth.CreateKeyFile(adapterKeyPath, AdapterKeySizeBytes)
	if err != nil {
		return fmt.Errorf("error creating adapter key file [%s]: %w", adapterKeyPath, err)
	}

	// Instantiate the adapter template
	adapterTemplatePath := filepath.Join(moduleInstallDir, string(info.Descriptor.Name), "adapter.tmpl")
	adapterTemplate, err := template.ParseFiles(adapterTemplatePath)
	if err != nil {
		return fmt.Errorf("error parsing adapter template [%s]: %w", adapterTemplatePath, err)
	}
	adapterRuntimePath := filepath.Join(adapterSrc.ModuleComposeDir, "adapter.yml")
	adapterFile, err := os.OpenFile(adapterRuntimePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, ModuleFileMode)
	if err != nil {
		return fmt.Errorf("error creating adapter compose file [%s]: %w", adapterRuntimePath, err)
	}
	err = adapterTemplate.Execute(adapterFile, adapterSrc)
	if err != nil {
		return fmt.Errorf("error executing adapter template: %w", err)
	}

	// Instantiate the service templates
	serviceSrc := templates.NewServiceDataSource(hdSettings, moduleSettingsMap, info, adapterSrc)
	serviceTemplatePath := filepath.Join(moduleInstallDir, string(info.Descriptor.Name), shared.TemplatesDir)
	entries, err := os.ReadDir(serviceTemplatePath)
	if err != nil {
		return fmt.Errorf("error reading service template directory [%s]: %w", serviceTemplatePath, err)
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
			return fmt.Errorf("error parsing service template [%s]: %w", templatePath, err)
		}
		serviceName := strings.TrimSuffix(filepath.Base(entryName), templateSuffix)
		serviceName += ".yml"
		serviceRuntimePath := filepath.Join(adapterSrc.ModuleComposeDir, serviceName)
		serviceFile, err := os.OpenFile(serviceRuntimePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, ModuleFileMode)
		if err != nil {
			return fmt.Errorf("error creating service compose file [%s]: %w", serviceRuntimePath, err)
		}
		err = template.Execute(serviceFile, serviceSrc)
		if err != nil {
			return fmt.Errorf("error executing service template: %w", err)
		}
	}
	return nil
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

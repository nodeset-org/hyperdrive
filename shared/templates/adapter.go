package templates

import (
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/shared"
)

// Enum type for an adapter mode.
type AdapterMode string

const (
	// Global mode, which is for read-only metadata about the module outside of a project context
	AdapterMode_Global AdapterMode = "global"

	// Project mode, which is for interacting with the adapter when it governs the fully instantiated module within a project
	AdapterMode_Project AdapterMode = "project"
)

// Environment variables
const (
	// The environment variable for the adapter mode
	AdapterModeEnvVarName string = "HD_ADAPTER_MODE"

	// The environment variable for the module config directory
	ModuleConfigDirEnvVarName string = "HD_CONFIG_DIR"

	// The environment variable for the module log directory
	ModuleLogDirEnvVarName string = "HD_LOG_DIR"

	// The environment variable for the module data directory
	ModuleDataDirEnvVarName string = "HD_DATA_DIR"

	// The environment variable for the adapter key file
	AdapterKeyFileEnvVarName string = "HD_KEY_FILE"

	// The environment variable for the directory of runtime Docker compose files
	ComposeDirEnvVarName string = "HD_COMPOSE_DIR"

	// The environment variable for the Docker compose project name
	ComposeProjectEnvVarName string = "HD_COMPOSE_PROJECT"
)

// Struct to pass into the template engine containing all necessary data and methods for populating a template.
type AdapterDataSource struct {
	AdapterContainerName        string
	AdapterEnvironmentVariables []string
	AdapterMode                 AdapterMode
	ModuleDir                   string
	ModuleConfigDir             string
	ModuleLogDir                string
	ModuleOverrideDir           string
	ModuleMetricsDir            string
	AdapterKeyFile              string
	ModuleComposeDir            string
	AdapterVolumes              []string
	ModuleNetwork               string
	ModuleComposeProject        string
}

func NewGlobalAdapterDataSource(descriptor *modules.ModuleDescriptor) *AdapterDataSource {
	moduleProjectName := shared.GlobalAdapterProjectName + "-" + string(descriptor.Shortcut)
	src := &AdapterDataSource{
		AdapterMode:          AdapterMode_Global,
		AdapterContainerName: moduleProjectName + "_adapter",
		AdapterEnvironmentVariables: []string{
			AdapterModeEnvVarName + "=" + string(AdapterMode_Global),
		},
		AdapterVolumes: []string{},
	}
	return src
}

func NewProjectAdapterDataSource(
	userDir string,
	hdProjectName string,
	descriptor *modules.ModuleDescriptor,
) *AdapterDataSource {
	moduleDir := filepath.Join(userDir, shared.ModulesDir, string(descriptor.Name))
	moduleComposeProject := hdProjectName + "-" + string(descriptor.Shortcut)
	moduleConfigDir := filepath.Join(moduleDir, shared.ConfigDir)
	moduleLogDir := filepath.Join(moduleDir, shared.LogsDir)
	moduleRuntimeDir := filepath.Join(moduleDir, shared.RuntimeDir)
	moduleOverrideDir := filepath.Join(moduleDir, shared.OverrideDir)
	moduleMetricsDir := filepath.Join(moduleDir, shared.MetricsDir)
	adapterKeyPath := filepath.Join(moduleDir, shared.SecretsDir, shared.AdapterKeyFile)

	src := &AdapterDataSource{
		ModuleDir:            moduleDir,
		ModuleComposeProject: moduleComposeProject,
		AdapterContainerName: moduleComposeProject + "_adapter",
		AdapterMode:          AdapterMode_Project,
		ModuleConfigDir:      moduleConfigDir,
		ModuleLogDir:         moduleLogDir,
		ModuleComposeDir:     moduleRuntimeDir,
		ModuleOverrideDir:    moduleOverrideDir,
		ModuleMetricsDir:     moduleMetricsDir,

		AdapterKeyFile: adapterKeyPath,
		ModuleNetwork:  hdProjectName + "_net",
		AdapterEnvironmentVariables: []string{
			AdapterModeEnvVarName + "=" + string(AdapterMode_Project),
			ModuleConfigDirEnvVarName + "=" + moduleConfigDir,
			ModuleLogDirEnvVarName + "=" + moduleLogDir,
			AdapterKeyFileEnvVarName + "=" + adapterKeyPath,
			ComposeDirEnvVarName + "=" + moduleRuntimeDir,
			ComposeProjectEnvVarName + "=" + moduleComposeProject,
		},
		AdapterVolumes: []string{
			moduleConfigDir + ":" + moduleConfigDir,
			moduleLogDir + ":" + moduleLogDir,
			moduleRuntimeDir + ":" + moduleRuntimeDir,
			moduleOverrideDir + ":" + moduleOverrideDir,
			moduleMetricsDir + ":" + moduleMetricsDir,
			adapterKeyPath + ":" + adapterKeyPath,
			"/var/run/docker.sock:/var/run/docker.sock",  // Access to the Docker socket
			"/usr/bin/docker:/usr/bin/docker:ro",         // Read-only access to the Docker binary
			"/usr/libexec/docker:/usr/libexec/docker:ro", // Read-only access to Docker Compose
		},
	}
	return src
}

package shared

import "path/filepath"

const (

	// The version of the CLI and Daemon
	HyperdriveVersion string = "2.0.0-dev"

	// The name of the directory for log files
	LogsDir string = "logs"

	// The name of the module directory, under the system directory
	ModulesDir string = "modules"

	// The name of the directory for global adapter compose files
	AdaptersDir string = "adapters"

	// The name of the config directory for module configs
	ConfigDir string = "config"

	// The name of the secrets directory, under the user config directory
	SecretsDir string = "secrets"

	// The name of the directory for compose file overrides
	OverrideDir string = "override"

	// The name of the directory for instantiated compose files
	RuntimeDir string = "runtime"

	// The name of the directory for metrics
	MetricsDir string = "metrics"

	// The name of the directory for the service templates
	TemplatesDir string = "templates"

	// The name of the Adapter's secret key file
	AdapterKeyFile string = "adapter.key"

	// The name of the Docker compose project for global adapter containers
	GlobalAdapterProjectName string = "hd"
)

// Get the full path of the modules directory
func GetModulesDirectoryPath(systemDir string) string {
	return filepath.Join(systemDir, ModulesDir)
}

// Get the full path of the directory for global adapter compose files
func GetGlobalAdapterDirectoryPath(systemDir string) string {
	return filepath.Join(systemDir, AdaptersDir)
}

// Get the full path of the secret key file for Adapters
func GetAdapterKeyPath(userConfigDir string) string {
	return filepath.Join(userConfigDir, SecretsDir, AdapterKeyFile)
}

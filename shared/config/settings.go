package config

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/modules/config"
	modutils "github.com/nodeset-org/hyperdrive/modules/utils"
	"github.com/nodeset-org/hyperdrive/shared"
	"gopkg.in/yaml.v3"
)

const (
	// The filename for Hyperdrive settings files
	SettingsFilename string = "user-settings.yml"

	// The folder for Hyperdrive config settings backups
	BackupConfigFolder string = "backups"

	// The pattern for Hyperdrive config settings backup filenames
	BackupConfigFilenamePattern string = "user-settings-%s.yml"

	// API base route for daemon requests
	HyperdriveDaemonRoute string = "hyperdrive"

	// API version for daemon requests
	HyperdriveApiVersion string = "1"

	// Complete API route for client requests
	HyperdriveApiClientRoute string = HyperdriveDaemonRoute + "/api/v" + HyperdriveApiVersion
)

// The instance of the Hyperdrive configuration
type HyperdriveSettings struct {
	Version                  string                            `json:"version" yaml:"version"`
	ProjectName              string                            `json:"projectName" yaml:"projectName"`
	ApiPort                  uint64                            `json:"apiPort" yaml:"apiPort"`
	EnableIPv6               bool                              `json:"enableIPv6" yaml:"enableIPv6"`
	UserDataPath             string                            `json:"userDataPath" yaml:"userDataPath"`
	AdditionalDockerNetworks string                            `json:"additionalDockerNetworks" yaml:"additionalDockerNetworks"`
	ClientTimeout            uint64                            `json:"clientTimeout" yaml:"clientTimeout"`
	ContainerTag             string                            `json:"containerTag" yaml:"containerTag"`
	Logging                  *LoggingConfigInstance            `json:"logging" yaml:"logging"`
	Modules                  map[string]*config.ModuleInstance `json:"modules" yaml:"modules"`
}

// Create a new Hyperdrive configuration instance with all of its settings set to the default values
func NewHyperdriveSettings() *HyperdriveSettings {
	cfg := NewHyperdriveConfig("", "")
	settings := config.CreateModuleSettings(cfg)

	var typedSettings HyperdriveSettings
	err := settings.ConvertToKnownType(&typedSettings)
	if err != nil {
		panic(fmt.Errorf("error converting Hyperdrive config to known type: %w", err))
	}
	typedSettings.Modules = map[string]*config.ModuleInstance{}
	return &typedSettings
}

// Create a copy of the Hyperdrive configuration instance
func (i HyperdriveSettings) CreateCopy() *HyperdriveSettings {
	// Serialize the instance to JSON
	bytes, err := json.Marshal(i)
	if err != nil {
		panic(fmt.Errorf("error serializing Hyperdrive config instance: %w", err))
	}

	// Deserialize the JSON back into a new instance
	newInstance := NewHyperdriveSettings()
	if err := json.Unmarshal(bytes, newInstance); err != nil {
		panic(fmt.Errorf("error deserializing Hyperdrive config instance: %w", err))
	}
	return newInstance
}

// Save an instance to a file, updating the version to be the current version of Hyperdrive
func (m HyperdriveSettings) SaveToFile(configFilePath string) error {
	// Serialize the module settings
	m.Version = shared.HyperdriveVersion

	// Serialize the instance
	configBytes, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not serialize config instance: %w", err)
	}

	// Write the file
	err = modutils.SafeSaveFile(configBytes, configFilePath, 0644)
	if err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}
	return nil
}

// Serialize the instance to a map, suitable for JSON serialization.
func (m HyperdriveSettings) SerializeToMap() map[string]any {
	// Serialize the module settings
	m.Version = shared.HyperdriveVersion

	// Serialize the instance
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("error serializing Hyperdrive config instance: %w", err))
	}

	// Deserialize the JSON to a map
	var instanceMap map[string]any
	if err := json.Unmarshal(bytes, &instanceMap); err != nil {
		panic(fmt.Errorf("error deserializing Hyperdrive config instance: %w", err))
	}
	return instanceMap
}

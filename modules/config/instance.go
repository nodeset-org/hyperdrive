package config

import (
	"fmt"

	"github.com/goccy/go-json"
)

// ModuleInstance represents an instance of a module that Hyperdrive is managing.
type ModuleInstance struct {
	// Whether or not the module is currently enabled
	Enabled bool `json:"enabled" yaml:"enabled"`

	// The module's settings (instance of its configuration)
	Settings ModuleInstanceSettingsContainer `json:"settings" yaml:"settings"`
}

type ModuleInstanceSettingsContainer struct {
	// The raw module's settings as it appears on disk without any conversion into a formal configuration instance
	rawSettings map[string]any

	// The module's configuration
	settings *ModuleConfigurationInstance
}

// Marshal the module info to JSON
func (i ModuleInstanceSettingsContainer) MarshalJSON() ([]byte, error) {
	if i.settings == nil {
		return json.Marshal(i.rawSettings)
	}

	return json.Marshal(i.settings.SerializeToMap())
}

// Marshal the module info to YAML
func (i ModuleInstanceSettingsContainer) MarshalYAML() (interface{}, error) {
	if i.settings == nil {
		return i.rawSettings, nil
	}
	return i.settings.SerializeToMap(), nil
}

func (i *ModuleInstanceSettingsContainer) UnmarshalJSON(data []byte) error {
	i.rawSettings = map[string]any{}
	err := json.Unmarshal(data, &i.rawSettings)
	if err != nil {
		return fmt.Errorf("error unmarshalling module settings: %w", err)
	}

	if i.settings != nil {
		return i.settings.DeserializeFromMap(i.rawSettings)
	}
	return nil
}

// Unmarshal the module info to YAML
func (i *ModuleInstanceSettingsContainer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	i.rawSettings = map[string]any{}
	err := unmarshal(&i.rawSettings)
	if err != nil {
		return fmt.Errorf("error unmarshalling module settings: %w", err)
	}

	if i.settings != nil {
		return i.settings.DeserializeFromMap(i.rawSettings)
	}
	return nil
}

// Loads the module's settings into a strongly-typed instance of the module's configuration based on its metadata. If the settings have already been loaded, this will overwrite it with a new instance.
func (i *ModuleInstanceSettingsContainer) CreateSettingsFromMetadata(metadata IModuleConfiguration) (*ModuleConfigurationInstance, error) {
	i.settings = CreateModuleConfigurationInstance(metadata)
	err := i.settings.DeserializeFromMap(i.rawSettings)
	if err != nil {
		return nil, fmt.Errorf("error deserializing module settings into instance: %w", err)
	}
	return i.settings, nil
}

// Gets the raw settings as a map without any type safety or validation. This is useful for modules that don't explicitly have a struct definition for the module's configuration, but want to explore it anyway.
func (i ModuleInstanceSettingsContainer) GetRawSettings() map[string]any {
	return i.rawSettings
}

// Gets the settings as a strongly typed instance of the module's configuration. This is useful for modules that have a struct definition for the module's configuration.
// If the settings haven't been loaded yet with CreateSettingsFromMetadata, this will return nil.
func (i ModuleInstanceSettingsContainer) GetSettings() *ModuleConfigurationInstance {
	return i.settings
}

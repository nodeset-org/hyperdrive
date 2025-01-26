package config

import (
	"fmt"

	"github.com/goccy/go-json"
)

// ModuleInstance represents an instance of a module that Hyperdrive is managing.
type ModuleInstance struct {
	// Whether or not the module is currently enabled
	Enabled bool `json:"enabled" yaml:"enabled"`

	// The version of the module that was used to create this instance
	Version string `json:"version" yaml:"version"`

	// The module's raw settings (instance of its configuration). Use the utility methods for ModuleInstance to convert it to a type-safe instance if needed.
	Settings map[string]any `json:"settings" yaml:"settings"`
}

// Converts the instance to a map, suitable for JSON serialization.
func (i ModuleInstance) SerializeToMap() map[string]any {
	instanceMap := map[string]any{
		"enabled":  i.Enabled,
		"version":  i.Version,
		"settings": i.Settings,
	}
	return instanceMap
}

// Creates a strongly-typed settings wrapper of the module's configuration based on its metadata, and loads the instance's settings into it.
func (i *ModuleInstance) CreateSettingsFromMetadata(metadata IModuleConfiguration) (*ModuleSettings, error) {
	settings := CreateModuleSettings(metadata)
	err := settings.DeserializeFromMap(i.Settings)
	if err != nil {
		return nil, fmt.Errorf("error deserializing module settings into instance: %w", err)
	}
	return settings, nil
}

// Loads the instance's settings directly into a known type. This uses JSON serialization to convert between the two types, so the known type must have the same JSON signature as the settings.
func (i *ModuleInstance) DeserializeSettingsIntoKnownType(knownType any) error {
	// Serialize the settings to JSON
	jsonBytes, err := json.Marshal(i.Settings)
	if err != nil {
		return fmt.Errorf("error serializing module settings to JSON: %w", err)
	}

	err = json.Unmarshal(jsonBytes, knownType)
	if err != nil {
		return fmt.Errorf("error deserializing module settings into known type: %w", err)
	}
	return nil
}

// Sets the instance's raw settings from a strongly-typed settings wrapper of the module's configuration.
func (i *ModuleInstance) SetSettings(settings *ModuleSettings) {
	i.Settings = settings.SerializeToMap()
}

// Sets the instance's raw settings from a known type. This uses JSON serialization to convert between the two types, so the known type must have the same JSON signature as the settings.
func (i *ModuleInstance) SetSettingsFromKnownType(knownType any) error {
	bytes, err := json.Marshal(knownType)
	if err != nil {
		return fmt.Errorf("error serializing module settings to JSON: %w", err)
	}

	var settings map[string]any
	if err := json.Unmarshal(bytes, &settings); err != nil {
		return fmt.Errorf("error deserializing module settings from JSON: %w", err)
	}
	i.Settings = settings
	return nil
}

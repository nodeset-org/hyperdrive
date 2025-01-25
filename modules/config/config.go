package config

import (
	"fmt"
)

// Top-level object of a module configuration
type IModuleConfiguration interface {
	IMetadataContainer
}

type moduleConfiguration struct {
	IMetadataContainer
}

// Marshal a configuration metadata object to a map, suitable for JSON serialization
func MarshalConfigurationToMap(cfg IModuleConfiguration) map[string]any {
	containerMap := map[string]any{}
	serializeContainerToMap(cfg, containerMap)
	return containerMap
}

// Unmarshal a configuration metadata object from a map, typically from a JSON response
func UnmarshalConfigurationFromMap(cfgMap map[string]any) (IModuleConfiguration, error) {
	container, err := deserializeContainerFromMap(cfgMap)
	if err != nil {
		return nil, fmt.Errorf("error deserializing configuration: %w", err)
	}
	moduleConfiguration := moduleConfiguration{
		IMetadataContainer: container,
	}
	return moduleConfiguration, nil
}

// Create a new instance of a module configuration's settings. All of the settings will be assigned their default values based on the configuration metadata.
func CreateModuleSettings(metadata IModuleConfiguration) *ModuleSettings {
	instance := &ModuleSettings{
		metadata:   metadata,
		parameters: map[Identifier]IParameterSetting{},
		sections:   map[Identifier]*SettingsSection{},
	}

	// Create the parameter instances
	for _, parameter := range metadata.GetParameters() {
		instance.parameters[parameter.GetID()] = parameter.CreateSetting()
	}

	// Create the subsection instances
	for _, subsection := range metadata.GetSections() {
		instance.sections[subsection.GetID()] = CreateSettingsSection(subsection)
	}

	return instance
}

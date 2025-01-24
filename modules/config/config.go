package config

import (
	"fmt"
)

const (
	VersionKey string = "version"
)

// Top-level object of a module configuration
type IModuleConfiguration interface {
	IMetadataContainer

	GetVersion() string
}

type moduleConfiguration struct {
	IMetadataContainer

	Version string
}

func (m moduleConfiguration) GetVersion() string {
	return m.Version
}

// Marshal a configuration metadata object to a map, suitable for JSON serialization
func MarshalConfigurationToMap(cfg IModuleConfiguration, version string) map[string]any {
	containerMap := map[string]any{}
	serializeContainerToMap(cfg, containerMap)
	containerMap[VersionKey] = version
	return containerMap
}

// Unmarshal a configuration metadata object from a map, typically from a JSON response
func UnmarshalConfigurationFromMap(cfgMap map[string]any) (IModuleConfiguration, error) {
	container, err := deserializeContainerFromMap(cfgMap)
	if err != nil {
		return nil, fmt.Errorf("error deserializing configuration: %w", err)
	}
	version, exists := cfgMap[VersionKey]
	if !exists {
		return nil, fmt.Errorf("missing version key")
	}
	versionString, ok := version.(string)
	if !ok {
		return nil, fmt.Errorf("version is not a string, it's a %T", version)
	}
	moduleConfiguration := moduleConfiguration{
		IMetadataContainer: container,
		Version:            versionString,
	}
	return moduleConfiguration, nil
}

// Create a new instance of a module configuration
func CreateModuleConfigurationInstance(metadata IModuleConfiguration) *ModuleSettings {
	instance := &ModuleSettings{
		metadata:   metadata,
		version:    metadata.GetVersion(),
		parameters: map[Identifier]IParameterInstance{},
		sections:   map[Identifier]*SectionInstance{},
	}

	// Create the parameter instances
	for _, parameter := range metadata.GetParameters() {
		instance.parameters[parameter.GetID()] = parameter.CreateInstance()
	}

	// Create the subsection instances
	for _, subsection := range metadata.GetSections() {
		instance.sections[subsection.GetID()] = CreateSectionInstance(subsection)
	}

	return instance
}

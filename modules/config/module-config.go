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

// Instance of a module configuration that will be serialized with the Hyperdrive config
type ModuleConfigurationInstance struct {
	// The metadata for this configuration
	metadata IModuleConfiguration

	// The version of the module that produced this instance
	version string

	// The parameters in this section
	parameters map[Identifier]IParameterInstance

	// The sections under this configuration
	sections map[Identifier]*SectionInstance
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
func CreateModuleConfigurationInstance(metadata IModuleConfiguration) *ModuleConfigurationInstance {
	instance := &ModuleConfigurationInstance{
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

// Get the metadata for this configuration instance
func (m ModuleConfigurationInstance) GetConfigurationMetadata() IModuleConfiguration {
	return m.metadata
}

// Get the version of the module that produced this instance
func (m ModuleConfigurationInstance) GetVersion() string {
	return m.version
}

// Get the parameter instance with the given ID
func (m ModuleConfigurationInstance) GetParameter(id Identifier) (IParameterInstance, error) {
	param, exists := m.parameters[id]
	if !exists {
		return nil, NewErrorNotFound(id, EntryType_Parameter)
	}
	return param, nil
}

// Get the section instance with the given ID
func (m ModuleConfigurationInstance) GetSection(id Identifier) (*SectionInstance, error) {
	section, exists := m.sections[id]
	if !exists {
		return nil, NewErrorNotFound(id, EntryType_Section)
	}
	return section, nil
}

// Internal method to get the parameters in this configuration instance
func (m ModuleConfigurationInstance) getParameters() map[Identifier]IParameterInstance {
	return m.parameters
}

// Internal method to get the sections in this configuration instance
func (m ModuleConfigurationInstance) getSections() map[Identifier]*SectionInstance {
	return m.sections
}

// Create a map of the configuration instance, suitable for marshalling
func (m ModuleConfigurationInstance) SerializeToMap() map[string]any {
	instanceMap := serializeContainerInstance(m)
	instanceMap[VersionKey] = m.GetVersion()
	return instanceMap
}

// Set the values from a map into the configuration instance
func (m *ModuleConfigurationInstance) DeserializeFromMap(instance map[string]any) error {
	// Get the version
	version, exists := instance[VersionKey]
	if !exists {
		return fmt.Errorf("missing version key")
	}
	versionString, ok := version.(string)
	if !ok {
		return fmt.Errorf("version is not a string, it's a %T", version)
	}
	m.version = versionString
	return deserializeContainerInstance(m, instance)
}

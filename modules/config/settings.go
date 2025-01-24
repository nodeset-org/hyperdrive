package config

import (
	"fmt"

	"github.com/goccy/go-json"
)

// Instance of a module configuration that will be serialized with the Hyperdrive config
type ModuleSettings struct {
	// The metadata for this configuration
	metadata IModuleConfiguration

	// The version of the module that produced this instance
	version string

	// The parameters in this section
	parameters map[Identifier]IParameterInstance

	// The sections under this configuration
	sections map[Identifier]*SectionInstance
}

// Get the metadata for this configuration instance
func (m ModuleSettings) GetConfigurationMetadata() IModuleConfiguration {
	return m.metadata
}

// Get the version of the module that produced this instance
func (m ModuleSettings) GetVersion() string {
	return m.version
}

// Get the parameter instance with the given ID
func (m ModuleSettings) GetParameter(id Identifier) (IParameterInstance, error) {
	param, exists := m.parameters[id]
	if !exists {
		return nil, NewErrorNotFound(id, EntryType_Parameter)
	}
	return param, nil
}

// Get the section instance with the given ID
func (m ModuleSettings) GetSection(id Identifier) (*SectionInstance, error) {
	section, exists := m.sections[id]
	if !exists {
		return nil, NewErrorNotFound(id, EntryType_Section)
	}
	return section, nil
}

// Internal method to get the parameters in this configuration instance
func (m ModuleSettings) getParameters() map[Identifier]IParameterInstance {
	return m.parameters
}

// Internal method to get the sections in this configuration instance
func (m ModuleSettings) getSections() map[Identifier]*SectionInstance {
	return m.sections
}

// Create a map of the configuration instance, suitable for marshalling
func (m ModuleSettings) SerializeToMap() map[string]any {
	instanceMap := serializeContainerInstance(m)
	instanceMap[VersionKey] = m.GetVersion()
	return instanceMap
}

// Set the values from a map into the configuration instance
func (m *ModuleSettings) DeserializeFromMap(instance map[string]any) error {
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

// Marshal the configuration instance to JSON
func (m ModuleSettings) MarshalJSON() ([]byte, error) {
	instanceMap := m.SerializeToMap()
	return json.Marshal(instanceMap)
}

// Unmarshal the configuration instance from JSON
func (m *ModuleSettings) UnmarshalJSON(data []byte) error {
	var instanceMap map[string]any
	if err := json.Unmarshal(data, &instanceMap); err != nil {
		return fmt.Errorf("error unmarshalling configuration instance: %w", err)
	}
	return m.DeserializeFromMap(instanceMap)
}

// Marshal the configuration instance to YAML
func (m ModuleSettings) MarshalYAML() (interface{}, error) {
	return m.SerializeToMap(), nil
}

// Unmarshal the configuration instance from YAML
func (m *ModuleSettings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var instanceMap map[string]any
	if err := unmarshal(&instanceMap); err != nil {
		return fmt.Errorf("error unmarshalling configuration instance: %w", err)
	}
	return m.DeserializeFromMap(instanceMap)
}

// Convert the generic configuration instance to a known struct type
func (m ModuleSettings) ConvertToKnownType(config any) error {
	// Serialize the instance to JSON
	bytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("error serializing configuration instance: %w", err)
	}

	// Deserialize the JSON back into the known type
	if err := json.Unmarshal(bytes, config); err != nil {
		return fmt.Errorf("error deserializing configuration instance: %w", err)
	}
	return nil
}

package config

import (
	"fmt"
)

// Top-level object of a module configuration
type IConfigurationMetadata interface {
	IMetadataContainer
}

// Marshal a configuration metadata object to a map, suitable for JSON serialization
func MarshalConfigurationMetadataToMap(metadata any) map[string]any {
	containerMap := map[string]any{}
	serializeContainerMetadataToMap(metadata, containerMap)
	return containerMap
}

// Unmarshal a configuration metadata object from a map, typically from a JSON response
func UnmarshalConfigurationMetadataFromMap(data map[string]any) (IConfigurationMetadata, error) {
	container, err := deserializeContainerMetadataFromMap(data)
	if err != nil {
		return nil, fmt.Errorf("error deserializing configuration: %w", err)
	}
	return container, nil
}

// Marshal a configuration instance to a map, suitable for JSON serialization
func CreateInstanceFromMetadata(metadata any) map[string]any {
	containerMap := map[string]any{}
	serializeContainerMetadataToInstance(metadata, containerMap)
	return containerMap
}

// Set the values from a configuration instance into a configuration metadata object
func UnmarshalConfigurationInstanceIntoMetadata(instance map[string]any, cfg any) error {
	return deserializeContainerInstanceToMetadata(instance, cfg)
}

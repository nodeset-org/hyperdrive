package config

import (
	"fmt"
)

// Top-level object of a module configuration
type IConfiguration interface {
	IMetadataContainer
}

// Marshal a configuration metadata object to a map, suitable for JSON serialization
func MarshalConfigurationToMap(cfg any) map[string]any {
	containerMap := map[string]any{}
	serializeContainerToMap(cfg, containerMap)
	return containerMap
}

// Unmarshal a configuration metadata object from a map, typically from a JSON response
func UnmarshalConfigurationFromMap(cfgMap map[string]any) (IConfiguration, error) {
	container, err := deserializeContainerFromMap(cfgMap)
	if err != nil {
		return nil, fmt.Errorf("error deserializing configuration: %w", err)
	}
	return container, nil
}

// Marshal a configuration instance to a map, suitable for JSON serialization
func CreateInstance(cfg any) map[string]any {
	containerMap := map[string]any{}
	serializeContainerToInstance(cfg, containerMap)
	return containerMap
}

// Set the values from a configuration instance into a configuration metadata object
func UnmarshalConfigurationInstance(instance map[string]any, cfg any) error {
	return deserializeContainerInstance(instance, cfg)
}

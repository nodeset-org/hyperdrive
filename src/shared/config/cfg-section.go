package config

import (
	"fmt"
	"reflect"
)

// Interface for describing config sections
type IConfigSection interface {
	// Get the name of the section (for display purposes)
	GetTitle() string

	// Get the list of parameters directly belonging to this section
	GetParameters() []IParameter

	// Get the sections underneath this one
	GetSubconfigs() map[string]IConfigSection
}

// Serialize a config section into a map
func Serialize(cfg IConfigSection) map[string]any {
	masterMap := map[string]any{}

	// Serialize parameters
	params := cfg.GetParameters()
	for _, param := range params {
		id := param.GetCommon().ID
		masterMap[id] = param.String()
	}

	// Serialize subconfigs
	subConfigs := cfg.GetSubconfigs()
	for name, subconfig := range subConfigs {
		masterMap[name] = Serialize(subconfig)
	}

	return masterMap
}

// Deserialize a config section
func Deserialize(cfg IConfigSection, serializedParams map[string]any, network Network) error {
	// Handle the parameters
	params := cfg.GetParameters()
	for _, param := range params {
		id := param.GetCommon().ID
		val, exists := serializedParams[id]
		if !exists {
			param.SetToDefault(network)
		} else {
			valString, isString := val.(string)
			if !isString {
				return fmt.Errorf("parameter [%s] is not a string but has a parameter ID name", id)
			}
			err := param.Deserialize(valString, network)
			if err != nil {
				return fmt.Errorf("error deserializing parameter [%s]: %w", id, err)
			}
		}
	}

	// Handle the subconfigs
	subconfigs := cfg.GetSubconfigs()
	for name, subconfig := range subconfigs {
		subParams, exists := serializedParams[name]
		if exists {
			submap, isMap := subParams.(map[string]any)
			if !isMap {
				return fmt.Errorf("subsection [%s] is not a map, it is %s", name, reflect.TypeOf(subParams))
			}
			err := Deserialize(subconfig, submap, network)
			if err != nil {
				return fmt.Errorf("error deserializing subsection [%s]: %w", name, err)
			}
		}
	}

	return nil
}

// Copy a section's settings into the corresponding section of a new config
func Clone(source IConfigSection, target IConfigSection, network Network) {
	// Handle the parameters
	targetParams := target.GetParameters()
	for i, sourceParam := range source.GetParameters() {
		targetParams[i].SetValue(sourceParam.GetValueAsAny())
		targetParams[i].GetCommon().UpdateDescription(network)
	}

	// Handle the subconfigs
	targetSubconfigs := target.GetSubconfigs()
	for i, sourceSubconfig := range source.GetSubconfigs() {
		Clone(sourceSubconfig, targetSubconfigs[i], network)
	}
}

// Change the active network for an entire configuration
func ChangeNetwork(cfg IConfigSection, oldNetwork Network, newNetwork Network) {
	// Update the master parameters
	params := cfg.GetParameters()
	for _, param := range params {
		param.ChangeNetwork(oldNetwork, newNetwork)
	}

	// Update all of the child config objects
	subconfigs := cfg.GetSubconfigs()
	for _, subconfig := range subconfigs {
		ChangeNetwork(subconfig, oldNetwork, newNetwork)
	}
}

// Update the default settings after a network change
func UpdateDefaults(cfg IConfigSection, newNetwork Network) {
	// Update the parameters
	for _, param := range cfg.GetParameters() {
		if param.GetCommon().OverwriteOnUpgrade {
			param.SetToDefault(newNetwork)
		}
	}

	// Update the subconfigs
	for _, subconfig := range cfg.GetSubconfigs() {
		UpdateDefaults(subconfig, newNetwork)
	}
}

// Apply the default settings for each parameter and subparameter
func ApplyDefaults(cfg IConfigSection, network Network) {
	// Update the parameters
	for _, param := range cfg.GetParameters() {
		param.SetToDefault(network)
	}

	// Update the subconfigs
	for _, subconfig := range cfg.GetSubconfigs() {
		ApplyDefaults(subconfig, network)
	}
}

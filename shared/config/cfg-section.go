package config

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Interface for describing config sections
type IConfigSection interface {
	// Get the name of the section (for display purposes)
	GetTitle() string

	// Get the list of parameters directly belonging to this section
	GetParameters() []types.IParameter

	// Get the sections underneath this one
	GetSubconfigs() map[string]IConfigSection
}

// Serialize a config section into a map
func serialize(cfg IConfigSection) map[string]any {
	masterMap := map[string]any{}

	// Serialize parameters
	params := cfg.GetParameters()
	for _, param := range params {
		id := param.GetCommon().ID
		masterMap[id] = param.GetValueAsString()
	}

	// Serialize subconfigs
	subConfigs := cfg.GetSubconfigs()
	for name, subconfig := range subConfigs {
		masterMap[name] = serialize(subconfig)
	}

	return masterMap
}

// Deserialize a config section
func deserialize(cfg IConfigSection, serializedParams map[string]any, network types.Network) error {
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
				return fmt.Errorf("subsection [%s] is not a map", name)
			}
			err := deserialize(subconfig, submap, network)
			if err != nil {
				return fmt.Errorf("error deserializing subsection [%s]: %w", name, err)
			}
		}
	}

	return nil
}

// Copy a section's settings into the corresponding section of a new config
func clone(source IConfigSection, target IConfigSection, network types.Network) {
	// Handle the parameters
	targetParams := target.GetParameters()
	for i, sourceParam := range source.GetParameters() {
		targetParams[i].SetValue(sourceParam.GetValueAsAny())
		targetParams[i].GetCommon().UpdateDescription(network)
	}

	// Handle the subconfigs
	targetSubconfigs := target.GetSubconfigs()
	for i, sourceSubconfig := range source.GetSubconfigs() {
		clone(sourceSubconfig, targetSubconfigs[i], network)
	}
}

// Change the active network for an entire configuration
func changeNetwork(cfg IConfigSection, oldNetwork types.Network, newNetwork types.Network) {
	// Update the master parameters
	params := cfg.GetParameters()
	for _, param := range params {
		param.ChangeNetwork(oldNetwork, newNetwork)
	}

	// Update all of the child config objects
	subconfigs := cfg.GetSubconfigs()
	for _, subconfig := range subconfigs {
		changeNetwork(subconfig, oldNetwork, newNetwork)
	}
}

// Update the default settings after a network change
func updateDefaults(cfg IConfigSection, newNetwork types.Network) {
	// Update the parameters
	for _, param := range cfg.GetParameters() {
		if param.GetCommon().OverwriteOnUpgrade {
			param.SetToDefault(newNetwork)
		}
	}

	// Update the subconfigs
	for _, subconfig := range cfg.GetSubconfigs() {
		updateDefaults(subconfig, newNetwork)
	}
}

// Apply the default settings for each parameter and subparameter
func applyDefaults(cfg IConfigSection, network types.Network) {
	// Update the parameters
	for _, param := range cfg.GetParameters() {
		param.SetToDefault(network)
	}

	// Update the subconfigs
	for _, subconfig := range cfg.GetSubconfigs() {
		applyDefaults(subconfig, network)
	}
}

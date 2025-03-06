package config

import (
	"fmt"
	"strings"

	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

// ===================================
// === configurationTemplateSource ===
// ===================================

// The base struct to use as the data source for configuration templates
type configurationTemplateSource struct {
	// The fully-qualified module name for the module that owns this parameter
	fqmn string

	// The Base Hyperdrive settings
	hdSettings *modconfig.ModuleSettings

	// The map of module FQMNs to their settings
	moduleSettingsMap map[string]*modconfig.ModuleSettings
}

// Get the value of a parameter setting from its fully-qualified path
func (s configurationTemplateSource) GetValue(fqpn string) (any, error) {
	fqmn := ""
	propertyPath := ""
	elements := strings.Split(fqpn, ":")
	if len(elements) == 1 {
		// This is a local property so use the module's fully qualified module name
		fqmn = s.fqmn
		propertyPath = fqpn
	} else {
		// TODO: Error out if there are more than 2 elements?
		fqmn = elements[0]
		propertyPath = elements[1]
	}

	// Get the module settings
	var settings *modconfig.ModuleSettings
	if fqmn == modules.HyperdriveFqmn {
		settings = s.hdSettings
	} else {
		var exists bool
		settings, exists = s.moduleSettingsMap[fqmn]
		if !exists {
			return "", fmt.Errorf("module settings not found for module [%s] in path [%s]", fqmn, propertyPath)
		}
	}
	return getModulePropertyValue(settings, propertyPath)
}

// Get the value of a parameter setting from its fully-qualified path, splitting it into an array using the delimiter
func (s configurationTemplateSource) GetValueArray(fqpn string, delimiter string) ([]string, error) {
	val, err := s.GetValue(fqpn)
	if err != nil {
		return nil, fmt.Errorf("error getting value for path [%s]: %w", fqpn, err)
	}
	valString, isString := val.(string)
	if !isString {
		return nil, fmt.Errorf("value for path [%s] is not a string", fqpn)
	}
	return strings.Split(valString, delimiter), nil
}

// ===============================
// === parameterTemplateSource ===
// ===============================

// The data source for dynamic property templates for parameters
type parameterTemplateSource struct {
	configurationTemplateSource

	// The parameter that owns the template
	parameter modconfig.IParameter
}

// Use the default value of the parameter setting
func (s parameterTemplateSource) UseDefault() any {
	return s.parameter.GetDefault()
}

// =============
// === Utils ===
// =============

// Get the value of a module settings property from its path
// TODO: unify with the function in shared/templates/service.go
func getModulePropertyValue(settings *modconfig.ModuleSettings, propertyPath string) (any, error) {
	// Split the param path into its components
	elements := strings.Split(propertyPath, "/")
	var container modconfig.IInstanceContainer = settings

	// Iterate through the sections
	level := 0
	for level < len(elements)-1 {
		elementString := elements[level]
		var id modconfig.Identifier
		err := id.UnmarshalText([]byte(elementString))
		if err != nil {
			return "", fmt.Errorf("error converting section [%s] in path [%s] to identifier: %w", elementString, propertyPath, err)
		}
		container, err = container.GetSection(id)
		if err != nil {
			return "", fmt.Errorf("error getting section [%s] in path [%s]: %w", elementString, propertyPath, err)
		}
		level++
	}

	// Get the parameter value
	elementString := elements[level]
	var id modconfig.Identifier
	err := id.UnmarshalText([]byte(elementString))
	if err != nil {
		return "", fmt.Errorf("error converting parameter [%s] in path [%s] to identifier: %w", elementString, propertyPath, err)
	}
	param, err := container.GetParameter(id)
	if err != nil {
		return "", fmt.Errorf("error getting parameter [%s] in path [%s]: %w", elementString, propertyPath, err)
	}
	return param.GetValue(), nil
}

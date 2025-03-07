package config

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/nodeset-org/hyperdrive/modules"
)

// Utility for processing dynamic property templates
type TemplateProcessor struct {
	// The base Hyperdrive settings
	hdSettings *ModuleSettings

	// The map of module FQMNs to their settings
	moduleSettingsMap map[string]*ModuleSettings
}

// Create a new template processor
func NewTemplateProcessor(hdSettings *ModuleSettings, moduleSettingsMap map[string]*ModuleSettings) *TemplateProcessor {
	return &TemplateProcessor{
		hdSettings:        hdSettings,
		moduleSettingsMap: moduleSettingsMap,
	}
}

// Execute a template owned by a parameter (such as the default or hidden properties).
func (p *TemplateProcessor) ExecuteTemplate(
	fqmn string,
	entity any,
	templateBody string,
) (string, error) {
	// Generate a template source for the entity
	var templateSource any
	switch entity.(type) {
	case IParameter:
		parameter := entity.(IParameter)
		templateSource = parameterTemplateSource{
			settingsTemplateSource: settingsTemplateSource{
				fqmn:              fqmn,
				hdSettings:        p.hdSettings,
				moduleSettingsMap: p.moduleSettingsMap,
			},
			parameter: parameter,
		}

	case IParameterOption:
		parameterOption := entity.(IParameterOption)
		parameter := parameterOption.GetOwner()
		templateSource = parameterTemplateSource{
			settingsTemplateSource: settingsTemplateSource{
				fqmn:              fqmn,
				hdSettings:        p.hdSettings,
				moduleSettingsMap: p.moduleSettingsMap,
			},
			parameter: parameter,
		}

	default:
		templateSource = settingsTemplateSource{
			fqmn:              fqmn,
			hdSettings:        p.hdSettings,
			moduleSettingsMap: p.moduleSettingsMap,
		}
	}

	// Execute the template
	template, err := template.New("").Parse(templateBody)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}
	result := &strings.Builder{}
	err = template.Execute(result, templateSource)
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return result.String(), nil
}

// Get the description of an entity, executing the template if provided or falling back to the default if not
func (p *TemplateProcessor) GetEntityDescription(
	fqmn string,
	entity IDescribable,
) (string, error) {
	description := entity.GetDescription()
	if description.Template == "" {
		return description.Default, nil
	}
	return p.ExecuteTemplate(fqmn, entity, description.Template)
}

// Get the hidden flag of an entity, executing the template if provided or falling back to the default if not
func (p *TemplateProcessor) IsEntityHidden(
	fqmn string,
	entity IHideable,
) (bool, error) {
	hidden := entity.GetHidden()
	if hidden.Template == "" {
		return hidden.Default, nil
	}
	result, err := p.ExecuteTemplate(fqmn, entity, hidden.Template)
	if err != nil {
		return false, err
	}
	isHidden, err := strconv.ParseBool(result)
	if err != nil {
		return false, fmt.Errorf("hidden template evaluated to [%s] which is not a boolean: %w", result, err)
	}
	return isHidden, nil
}

// Get the disabled flag of an entity, executing the template if provided or falling back to the default if not
func (p *TemplateProcessor) IsEntityDisabled(
	fqmn string,
	entity IDisableable,
) (bool, error) {
	disabled := entity.GetDisabled()
	if disabled.Template == "" {
		return disabled.Default, nil
	}
	result, err := p.ExecuteTemplate(fqmn, entity, disabled.Template)
	if err != nil {
		return false, err
	}
	isDisabled, err := strconv.ParseBool(result)
	if err != nil {
		return false, fmt.Errorf("disabled template evaluated to [%s] which is not a boolean: %w", result, err)
	}
	return isDisabled, nil
}

// ==============================
// === settingsTemplateSource ===
// ==============================

// The base struct to use as the data source for property templates
type settingsTemplateSource struct {
	// The fully-qualified module name for the module that owns this parameter
	fqmn string

	// The Base Hyperdrive settings
	hdSettings *ModuleSettings

	// The map of module FQMNs to their settings
	moduleSettingsMap map[string]*ModuleSettings
}

// Get the value of a parameter setting from its fully-qualified path
func (s settingsTemplateSource) GetValue(fqpn string) (any, error) {
	fqmn := ""
	propertyPath := ""
	elements := strings.Split(fqpn, ":")
	if len(elements) == 1 {
		// This is a local property so use the module's fully qualified module name
		fqmn = s.fqmn
		propertyPath = fqpn
	} else if len(elements) == 2 {
		fqmn = elements[0]
		propertyPath = elements[1]
	} else {
		return "", fmt.Errorf("invalid fully-qualified property name [%s]", fqpn)
	}

	// Get the module settings
	var settings *ModuleSettings
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
func (s settingsTemplateSource) GetValueArray(fqpn string, delimiter string) ([]string, error) {
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
	settingsTemplateSource

	// The parameter that owns the template
	parameter IParameter
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
func getModulePropertyValue(settings *ModuleSettings, propertyPath string) (any, error) {
	// Split the param path into its components
	elements := strings.Split(propertyPath, "/")
	var container IInstanceContainer = settings

	// Iterate through the sections
	level := 0
	for level < len(elements)-1 {
		elementString := elements[level]
		var id Identifier
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
	var id Identifier
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

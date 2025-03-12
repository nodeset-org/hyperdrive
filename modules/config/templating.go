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

// Get the description of an entity, executing the template if provided or falling back to the default if not
func (p *TemplateProcessor) GetEntityDescription(
	fqmn string,
	entity IDescribable,
) (string, error) {
	description := entity.GetDescription()
	return processDynamicProperty(&description, fqmn, p.hdSettings, p.moduleSettingsMap)
}

// Get the hidden flag of an entity, executing the template if provided or falling back to the default if not
func (p *TemplateProcessor) IsEntityHidden(
	fqmn string,
	entity IHideable,
) (bool, error) {
	hidden := entity.GetHidden()
	return processDynamicProperty(&hidden, fqmn, p.hdSettings, p.moduleSettingsMap)
}

// Get the disabled flag of an entity, executing the template if provided or falling back to the default if not
func (p *TemplateProcessor) IsEntityDisabled(
	fqmn string,
	entity IDisableable,
) (bool, error) {
	disabled := entity.GetDisabled()
	return processDynamicProperty(&disabled, fqmn, p.hdSettings, p.moduleSettingsMap)
}

// The base struct to use as the data source for property templates
type dynamicPropertyTemplateSource[Type any] struct {
	// The property that owns the template
	property *DynamicProperty[Type]

	// The fully-qualified module name for the module that owns this parameter
	fqmn string

	// The Base Hyperdrive settings
	hdSettings *ModuleSettings

	// The map of module FQMNs to their settings
	moduleSettingsMap map[string]*ModuleSettings
}

// Get the value of a parameter setting from its fully-qualified path
func (s *dynamicPropertyTemplateSource[Type]) GetValue(fqpn string) (any, error) {
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
		return "", fmt.Errorf("invalid fully-qualified property name \"%s\"", fqpn)
	}

	// Get the module settings
	var settings *ModuleSettings
	if fqmn == modules.HyperdriveFqmn {
		settings = s.hdSettings
	} else {
		var exists bool
		settings, exists = s.moduleSettingsMap[fqmn]
		if !exists {
			return "", fmt.Errorf("module settings not found for module \"%s\" in path \"%s\"", fqmn, propertyPath)
		}
	}
	return getModulePropertyValue(settings, propertyPath)
}

// Get the value of a parameter setting from its fully-qualified path, splitting it into an array using the delimiter
func (s *dynamicPropertyTemplateSource[Type]) GetValueArray(fqpn string, delimiter string) ([]string, error) {
	val, err := s.GetValue(fqpn)
	if err != nil {
		return nil, fmt.Errorf("error getting value for path \"%s\": %w", fqpn, err)
	}
	valString, isString := val.(string)
	if !isString {
		return nil, fmt.Errorf("value for path \"%s\" is not a string", fqpn)
	}
	return strings.Split(valString, delimiter), nil
}

// Get the default value of the property
func (s *dynamicPropertyTemplateSource[Type]) UseDefault() Type {
	return s.property.Default
}

// Execute a template owned by a parameter (such as the default or hidden properties).
func processDynamicProperty[Type any](
	property *DynamicProperty[Type],
	fqmn string,
	hdSettings *ModuleSettings,
	moduleSettingsMap map[string]*ModuleSettings,
) (Type, error) {
	if property.Template == "" {
		return property.Default, nil
	}
	templateSource := createDynamicPropertyTemplateSource(property, fqmn, hdSettings, moduleSettingsMap)

	// Execute the template
	var result Type
	template, err := template.New("").Parse(property.Template)
	if err != nil {
		return result, fmt.Errorf("error parsing template: %w", err)
	}
	resultBuilder := &strings.Builder{}
	err = template.Execute(resultBuilder, templateSource)
	if err != nil {
		return result, fmt.Errorf("error executing template: %w", err)
	}

	// Convert the result to the right type
	// TODO: probably use a formal marshaller for this instead of manually handling types
	err = nil
	switch resultPtr := any(&result).(type) {
	case *bool:
		*resultPtr, err = strconv.ParseBool(resultBuilder.String())
	case *string:
		*resultPtr = resultBuilder.String()
	case *int:
		*resultPtr, err = strconv.Atoi(resultBuilder.String())
	case *int8:
		var resultTyped int64
		resultTyped, err = strconv.ParseInt(resultBuilder.String(), 10, 8)
		*resultPtr = int8(resultTyped)
	case *int16:
		var resultTyped int64
		resultTyped, err = strconv.ParseInt(resultBuilder.String(), 10, 16)
		*resultPtr = int16(resultTyped)
	case *int32:
		var resultTyped int64
		resultTyped, err = strconv.ParseInt(resultBuilder.String(), 10, 32)
		*resultPtr = int32(resultTyped)
	case *int64:
		*resultPtr, err = strconv.ParseInt(resultBuilder.String(), 10, 64)
	case *uint:
		var resultTyped uint64
		resultTyped, err = strconv.ParseUint(resultBuilder.String(), 10, 0)
		*resultPtr = uint(resultTyped)
	case *uint8:
		var resultTyped uint64
		resultTyped, err = strconv.ParseUint(resultBuilder.String(), 10, 8)
		*resultPtr = uint8(resultTyped)
	case *uint16:
		var resultTyped uint64
		resultTyped, err = strconv.ParseUint(resultBuilder.String(), 10, 16)
		*resultPtr = uint16(resultTyped)
	case *uint32:
		var resultTyped uint64
		resultTyped, err = strconv.ParseUint(resultBuilder.String(), 10, 32)
		*resultPtr = uint32(resultTyped)
	case *uint64:
		*resultPtr, err = strconv.ParseUint(resultBuilder.String(), 10, 64)
	case *float32:
		var resultTyped float64
		resultTyped, err = strconv.ParseFloat(resultBuilder.String(), 32)
		*resultPtr = float32(resultTyped)
	case *float64:
		*resultPtr, err = strconv.ParseFloat(resultBuilder.String(), 64)
	default:
		err = fmt.Errorf("invalid type for dynamic property: %T", result)
	}
	return result, err
}

// Create a template source for a dynamic property (used for type inference)
func createDynamicPropertyTemplateSource[Type any](
	property *DynamicProperty[Type],
	fqmn string,
	hdSettings *ModuleSettings,
	moduleSettingsMap map[string]*ModuleSettings,
) *dynamicPropertyTemplateSource[Type] {
	return &dynamicPropertyTemplateSource[Type]{
		property:          property,
		fqmn:              fqmn,
		hdSettings:        hdSettings,
		moduleSettingsMap: moduleSettingsMap,
	}
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
			return "", fmt.Errorf("error converting section \"%s\" in path \"%s\" to identifier: %w", elementString, propertyPath, err)
		}
		container, err = container.GetSection(id)
		if err != nil {
			return "", fmt.Errorf("error getting section \"%s\" in path \"%s\": %w", elementString, propertyPath, err)
		}
		level++
	}

	// Get the parameter value
	elementString := elements[level]
	var id Identifier
	err := id.UnmarshalText([]byte(elementString))
	if err != nil {
		return "", fmt.Errorf("error converting parameter \"%s\" in path \"%s\" to identifier: %w", elementString, propertyPath, err)
	}
	param, err := container.GetParameter(id)
	if err != nil {
		return "", fmt.Errorf("error getting parameter \"%s\" in path \"%s\": %w", elementString, propertyPath, err)
	}
	return param.GetValue(), nil
}

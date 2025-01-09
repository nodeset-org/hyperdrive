package config

import (
	"errors"
	"fmt"
)

// Deserialize a named property from a map. Assumes the value deserialized by the underlying JSON unmarshaller
// converts the property to the right type.
func deserializeProperty[Type any](data map[string]any, propertyName string, property *Type, optional bool) (bool, error) {
	// Get the property by its name
	value, exists := data[propertyName]
	if !exists {
		if optional {
			return false, nil
		}
		return false, fmt.Errorf("missing property %s", propertyName)
	}

	// Handle nil values
	if optional && value == nil {
		return false, nil
	}

	// Convert it to the right type
	propertyTyped, ok := value.(Type)
	if !ok {
		return true, fmt.Errorf("invalid property %s [%v]", propertyName, value)
	}

	// Set the property
	*property = propertyTyped
	return true, nil
}

func setNumberProperty[Type NumberParameterType](id string, property *Type, value any) error {
	// Check for direct assignment
	typedValue, ok := value.(Type)
	if ok {
		*property = typedValue
	}

	// Handle conversion
	switch valType := value.(type) {
	case int:
		*property = Type(valType)
	case int8:
		*property = Type(valType)
	case int16:
		*property = Type(valType)
	case int32:
		*property = Type(valType)
	case int64:
		*property = Type(valType)
	case uint:
		*property = Type(valType)
	case uint8:
		*property = Type(valType)
	case uint16:
		*property = Type(valType)
	case uint32:
		*property = Type(valType)
	case uint64:
		*property = Type(valType)
	case float32:
		*property = Type(valType)
	case float64:
		*property = Type(valType)
	default:
		return fmt.Errorf("invalid type for number property [%s]: %T", id, value)
	}
	return nil
}

func parseChoiceParameter[ChoiceType ~string](instance map[string]any, paramID string, param *ChoiceParameterMetadata[ChoiceType]) error {
	paramAny, exists := instance[paramID]
	if !exists {
		return errors.New("missing required parameter: " + paramID)
	}
	paramString, ok := paramAny.(string)
	if !ok {
		return fmt.Errorf("invalid type for parameter [%s]: %T", paramID, paramAny)
	}
	param.Value = ChoiceType(paramString)
	return nil
}

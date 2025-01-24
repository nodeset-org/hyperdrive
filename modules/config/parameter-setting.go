package config

import "fmt"

// A general-purpose instance of a parameter. Useful for working with parameters of dynamic module configurations that you don't know the type of during compile time.
type IParameterSetting interface {
	// Get the metadata for the parameter
	GetMetadata() IParameter

	// Get the value of the parameter
	GetValue() any

	// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
	SetValue(value any) error

	// Get the parameter's value as a string
	String() string
}

// Underlying implementation for parameter settings
type parameterSetting[Type any] struct {
	// The parameter metadata
	metadata IParameter

	// The current value of the parameter
	Value Type
}

// Gets the metadata for the parameter
func (p *parameterSetting[Type]) GetMetadata() IParameter {
	return p.metadata
}

// Gets the value of the parameter
func (p *parameterSetting[Type]) GetValue() any {
	return p.Value
}

// Gets the value of the parameter as a string
func (p *parameterSetting[Type]) String() string {
	return fmt.Sprintf("%v", p.Value)
}

/// =======================
/// === Bool Parameters ===
/// =======================

// Underlying implementation for bool parameter settings
type boolParameterSetting struct {
	parameterSetting[bool]
}

// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
func (i *boolParameterSetting) SetValue(value any) error {
	boolValue, ok := value.(bool)
	if !ok {
		return fmt.Errorf("invalid value type for bool parameter [%s]: %T", i.metadata.GetID(), value)
	}
	i.Value = boolValue
	return nil
}

/// =======================================
/// === Prototype for Number Parameters ===
/// =======================================

// Underlying implementation for number parameter settings
type numberParameterSetting[Type NumberParameterType] struct {
	parameterSetting[Type]
}

// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
func (i *numberParameterSetting[Type]) SetValue(value any) error {
	return setNumberProperty(i.metadata.GetID().String(), &i.Value, value)
}

/// =========================
/// === String Parameters ===
/// =========================

// Underlying implementation for string parameter settings
type stringParameterSetting struct {
	parameterSetting[string]
}

// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
func (i *stringParameterSetting) SetValue(value any) error {
	stringValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value type for string parameter [%s]: %T", i.metadata.GetID(), value)
	}
	i.Value = stringValue
	return nil
}

/// =========================
/// === Choice Parameters ===
/// =========================

// Underlying implementation for choice parameter settings
type choiceParameterSetting[ChoiceType ~string] struct {
	parameterSetting[ChoiceType]
}

// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
func (i *choiceParameterSetting[ChoiceType]) SetValue(value any) error {
	// Check for direct assignment
	typedValue, ok := value.(ChoiceType)
	if ok {
		i.Value = typedValue
	}

	// Handle string conversion
	stringValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value type for choice parameter %s: %T", i.metadata.GetID(), value)
	}
	i.Value = ChoiceType(stringValue)
	return nil
}

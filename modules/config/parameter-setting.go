package config

import (
	"fmt"
	"regexp"
)

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

	// Validate the parameter's value. Any issues will be returned as a list of errors.
	Validate() []error
}

// Underlying implementation for parameter settings
type parameterSetting[InfoType IParameter, ValueType any] struct {
	// The metadata for the parameter
	info InfoType

	// The current value of the parameter
	value ValueType
}

// Get the metadata for the parameter
func (p *parameterSetting[InfoType, ValueType]) GetMetadata() IParameter {
	return p.info
}

// Gets the value of the parameter
func (p *parameterSetting[InfoType, ValueType]) GetValue() any {
	return p.value
}

// Gets the value of the parameter as a string
func (p *parameterSetting[InfoType, Type]) String() string {
	return fmt.Sprintf("%v", p.value)
}

/// =======================
/// === Bool Parameters ===
/// =======================

// Underlying implementation for bool parameter settings
type boolParameterSetting struct {
	parameterSetting[IBoolParameter, bool]
}

// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
func (i *boolParameterSetting) SetValue(value any) error {
	boolValue, ok := value.(bool)
	if !ok {
		return fmt.Errorf("invalid value type for bool parameter [%s]: %T", i.info.GetID(), value)
	}
	i.value = boolValue
	return nil
}

// Validate the parameter's value. Any issues will be returned as a list of errors.
func (i *boolParameterSetting) Validate() []error {
	// No validation needed for bool parameters
	return []error{}
}

/// =======================================
/// === Prototype for Number Parameters ===
/// =======================================

// Underlying implementation for number parameter settings
type numberParameterSetting[Type NumberParameterType] struct {
	parameterSetting[INumberParameter[Type], Type]
}

// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
func (i *numberParameterSetting[Type]) SetValue(value any) error {
	return setNumberProperty(i.info.GetID().String(), &i.value, value)
}

// Validate the parameter's value. Any issues will be returned as a list of errors.
func (i *numberParameterSetting[Type]) Validate() []error {
	val := i.value
	min := i.info.GetMinValue()
	max := i.info.GetMaxValue()
	errors := []error{}

	if val < min {
		errors = append(errors, fmt.Errorf("value [%v] is less than the minimum value of %v", val, min))
	}
	if max > 0 && val > max {
		errors = append(errors, fmt.Errorf("value [%v] is greater than the maximum value of %v", val, max))
	}
	return errors
}

/// =========================
/// === String Parameters ===
/// =========================

// Underlying implementation for string parameter settings
type stringParameterSetting struct {
	parameterSetting[IStringParameter, string]
}

// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
func (i *stringParameterSetting) SetValue(value any) error {
	stringValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value type for string parameter [%s]: %T", i.info.GetID(), value)
	}
	i.value = stringValue
	return nil
}

// Validate the parameter's value. Any issues will be returned as a list of errors.
func (i *stringParameterSetting) Validate() []error {
	// Make sure the length is within the allowed range
	val := i.value
	max := i.info.GetMaxLength()
	regexPattern := i.info.GetRegex()
	errors := []error{}

	if max > 0 && uint64(len(val)) > max {
		errors = append(errors, fmt.Errorf("string length [%d] exceeds maximum length of %d", len(val), max))
	}

	// Build the regex
	// TODO: this needs to be in the module loader, not in the settings validator
	if regexPattern != "" {
		regex, err := regexp.Compile(regexPattern)
		if err != nil {
			errors = append(errors, fmt.Errorf("invalid regex pattern [%s]: %w", regexPattern, err))
		}
		if !regex.MatchString(val) {
			errors = append(errors, fmt.Errorf("string [%s] does not match regex pattern [%s]", val, regexPattern))
		}
	}

	return errors
}

/// =========================
/// === Choice Parameters ===
/// =========================

// Underlying implementation for choice parameter settings
type choiceParameterSetting[ChoiceType ~string] struct {
	parameterSetting[IChoiceParameter, ChoiceType]
}

// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
func (i *choiceParameterSetting[ChoiceType]) SetValue(value any) error {
	// Check for direct assignment
	typedValue, ok := value.(ChoiceType)
	if ok {
		i.value = typedValue
		return nil
	}

	// Handle string conversion
	stringValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value type for choice parameter %s: %T", i.info.GetID(), value)
	}
	i.value = ChoiceType(stringValue)
	return nil
}

// Validate the parameter's value. Any issues will be returned as a list of errors.
func (i *choiceParameterSetting[ChoiceType]) Validate() []error {
	// Make sure the value is in the list of choices
	val := i.value
	options := i.info.GetOptions()
	errors := []error{}

	optionFound := false
	for _, option := range options {
		if option.GetValue() == string(val) {
			optionFound = true
			break
		}
	}
	if !optionFound {
		errors = append(errors, fmt.Errorf("value [%s] is not a valid choice", val))
	}

	return errors
}

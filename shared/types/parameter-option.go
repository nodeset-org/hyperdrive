package types

import "fmt"

// Common fields across all ParameterOption instances
type ParameterOptionCommon struct {
	// The option's human-readable name, to be used in config displays
	Name string

	// A description signifying what this option means
	Description string
}

// A single option in a choice parameter
type ParameterOption[Type any] struct {
	*ParameterOptionCommon

	// The underlying value for this option
	Value Type
}

// An interface for typed ParameterOption structs, to get common fields from them
type IParameterOption interface {
	// Get the parameter option's common fields
	Common() *ParameterOptionCommon

	// Get the option's value
	GetValueAsAny() any

	// Ge the option's value as a string
	String() string
}

// Get the parameter option's common fields
func (p *ParameterOption[_]) Common() *ParameterOptionCommon {
	return p.ParameterOptionCommon
}

// Get the parameter's value
func (p *ParameterOption[_]) GetValueAsAny() any {
	return p.Value
}

// Get the parameter option's value as a string
func (p *ParameterOption[_]) String() string {
	return fmt.Sprint(p.Value)
}

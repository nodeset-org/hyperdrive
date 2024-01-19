package types

import (
	"fmt"
	"regexp"
	"strconv"
)

// A setting that has changed
type ChangedSetting struct {
	Name               string
	OldValue           string
	NewValue           string
	AffectedContainers map[ContainerID]bool
}

// =========================
// === Parameter Options ===
// =========================

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
	GetValueAsString() string
}

// Get the parameter option's common fields
func (p *ParameterOption[_]) Common() *ParameterOptionCommon {
	return p.ParameterOptionCommon
}

// Get the parameter option's value as a string
func (p *ParameterOption[_]) GetValueAsString() string {
	return fmt.Sprint(p.Value)
}

// ==================
// === Parameters ===
// ==================

// Common fields across all Parameter instances
type ParameterCommon struct {
	// The parameter's ID, used for serialization and deserialization
	ID string

	// The parameter's human-readable name
	Name string

	// A description of this parameter / setting
	Description string

	// The max length of the parameter, in characters, if it's free-form input
	MaxLength int

	// An optional regex used to validate free-form input for the parameter
	Regex string

	// True if this is an advanced parameter and should be hidden unless advanced configuration mode is enabled
	Advanced bool

	// The list of Docker containers affected by changing this parameter
	// (these containers will require a restart for the change to take effect)
	AffectsContainers []ContainerID

	// A list of Docker container environment variables that should be set to this parameter's value
	EnvironmentVariables []string

	// Whether or not the parameter is allowed to be blank
	CanBeBlank bool

	// True to reset the parameter's value to the default option after Hyperdrive is updated
	OverwriteOnUpgrade bool

	// Descriptions of the parameter that change depending on the selected network
	DescriptionsByNetwork map[Network]string
}

// A parameter that can be configured by the user
type Parameter[Type any] struct {
	*ParameterCommon
	Default map[Network]Type
	Value   Type
	Options []*ParameterOption[Type]
}

// An interface for typed Parameter structs, to get common fields from them
type IParameter interface {
	// Get the parameter's common fields
	GetCommon() *ParameterCommon

	// Get the common fields from each ParameterOption (returns nil if this isn't a choice parameter)
	GetOptions() []*ParameterOptionCommon

	// Set the parameter to the default value
	SetToDefault(network Network) error

	// Get the parameter's value as a string
	GetValueAsString() string

	// Deserializes a string into this parameter's value
	Deserialize(serializedParam string, network Network) error
}

// Get the parameter's common fields
func (p *Parameter[_]) GetCommon() *ParameterCommon {
	return p.ParameterCommon
}

// Get the common fields from each ParameterOption (returns nil if this isn't a choice parameter)
func (p *Parameter[_]) GetOptions() []*ParameterOptionCommon {
	if len(p.Options) == 0 {
		return nil
	}
	opts := make([]*ParameterOptionCommon, len(p.Options))
	for i, param := range p.Options {
		opts[i] = param.ParameterOptionCommon
	}
	return opts
}

// Set the value to the default for the provided config's network
func (p *Parameter[_]) SetToDefault(network Network) error {
	defaultSetting, err := p.GetDefault(network)
	if err != nil {
		return err
	}
	p.Value = defaultSetting
	return nil
}

// Get the default value for the provided network
func (p *Parameter[Type]) GetDefault(network Network) (Type, error) {
	var empty Type
	defaultSetting, exists := p.Default[network]
	if !exists {
		defaultSetting, exists = p.Default[Network_All]
		if !exists {
			return empty, fmt.Errorf("parameter [%s] doesn't have a default for network %s or all networks", p.ID, network)
		}
	}

	return defaultSetting, nil
}

// Get the parameter's value as a string
func (p *Parameter[_]) GetValueAsString() string {
	return fmt.Sprint(p.Value)
}

// Deserializes a string into this parameter's value
func (p *Parameter[Type]) Deserialize(serializedParam string, network Network) error {
	if len(p.Options) > 0 {
		for _, option := range p.Options {
			optionVal := option.GetValueAsString()
			if optionVal == serializedParam {
				p.Value = option.Value
				return nil
			}
		}
		return p.SetToDefault(network)
	}

	var err error
	switch value := any(&p.Value).(type) {
	case *int64:
		*value, err = strconv.ParseInt(serializedParam, 0, 0)
	case *uint64:
		*value, err = strconv.ParseUint(serializedParam, 0, 0)
	case *uint16:
		var result uint64
		result, err = strconv.ParseUint(serializedParam, 0, 16)
		*value = uint16(result)
	case *bool:
		*value, err = strconv.ParseBool(serializedParam)
	case *float64:
		*value, err = strconv.ParseFloat(serializedParam, 64)
	case *string:
		if p.CanBeBlank && serializedParam == "" {
			*value = ""
			return nil
		}
		if p.MaxLength > 0 && len(serializedParam) > p.MaxLength {
			return fmt.Errorf("cannot deserialize parameter [%s]: value [%s] is longer than the max length of [%d]", p.ID, serializedParam, p.MaxLength)
		}
		if p.Regex != "" {
			regex := regexp.MustCompile(p.Regex)
			if !regex.MatchString(serializedParam) {
				return fmt.Errorf("cannot deserialize parameter [%s]: value [%s] did not match the expected format", p.ID, serializedParam)
			}
		}
		if !p.CanBeBlank && serializedParam == "" {
			return p.SetToDefault(network)
		}
		*value = serializedParam
	}

	if err != nil {
		return fmt.Errorf("cannot deserialize parameter [%s]: %w", p.ID, err)
	}

	return nil
}

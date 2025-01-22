package config

import (
	"fmt"
)

// Parameter types
type ParameterType string

const (
	ParameterType_Bool   ParameterType = "bool"
	ParameterType_Int    ParameterType = "int"
	ParameterType_Uint   ParameterType = "uint"
	ParameterType_Float  ParameterType = "float"
	ParameterType_String ParameterType = "string"
	ParameterType_Choice ParameterType = "choice"
)

// Common interface for all Parameter metadata structs
type IParameter interface {
	// Gets the unique ID for the parameter
	GetID() Identifier

	// Gets the human-readable name for the parameter
	GetName() string

	// Gets the description of the parameter
	GetDescription() DynamicProperty[string]

	// Gets the type of the parameter - this can be used for casting to a specific parameter type
	GetType() ParameterType

	// Gets the default value for the parameter as an any
	GetDefault() any

	// Gets whether the parameter is considered "advanced"
	GetAdvanced() bool

	// Gets whether the parameter is disabled
	GetDisabled() DynamicProperty[bool]

	// Gets whether the parameter is hidden
	GetHidden() DynamicProperty[bool]

	// Gets whether the parameter should be overwritten with the default on an upgrade
	GetOverwriteOnUpgrade() bool

	// Gets the list of containers affected if this parameter is changed
	GetAffectedContainers() []string

	// Serializes the parameter metadata to a map
	Serialize() map[string]any

	// Deserializes the parameter metadata from a map
	Deserialize(data map[string]any) error

	// Creates a parameter instance that's linked to this metadata
	CreateInstance() IParameterInstance
}

// A general-purpose instance of a parameter. Useful for working with parameters of dynamic module configurations that you don't know the type of during compile time.
type IParameterInstance interface {
	// Get the metadata for the parameter
	GetMetadata() IParameter

	// Get the value of the parameter
	GetValue() any

	// Set the value of the parameter. If the provided value is the wrong type for the parameter, an error will be returned.
	SetValue(value any) error

	// Get the parameter's value as a string
	String() string
}

// ===========================
/// === Parameter Metadata ===
// ===========================

const (
	// Field names
	IDKey                 string = "id"
	NameKey               string = "name"
	DescriptionKey        string = "description"
	TypeKey               string = "type"
	DefaultKey            string = "default"
	ValueKey              string = "value"
	AdvancedKey           string = "advanced"
	DisabledKey           string = "disabled"
	HiddenKey             string = "hidden"
	OverwriteOnUpgradeKey string = "overwriteOnUpgrade"
	AffectedContainersKey string = "affectedContainers"
)

// Parameter metadata implementation according to the spec
type Parameter[Type any] struct {
	// Unique ID for referencing the parameter behind-the-scenes
	ID Identifier `json:"id" yaml:"id"`

	// Human-readable name for the parameter
	Name string `json:"name" yaml:"name"`

	// Description of the parameter
	Description DynamicProperty[string] `json:"description" yaml:"description"`

	// Default value for the parameter
	Default Type `json:"default" yaml:"default"`

	// Flag for hiding the parameter behind the "advanced mode" toggle
	Advanced bool `json:"advanced,omitempty" yaml:"advanced,omitempty"`

	// Flag for disabling the parameter in the UI, graying it out
	Disabled DynamicProperty[bool] `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// Dynamic flag for hiding the parameter from the UI
	Hidden DynamicProperty[bool] `json:"hidden,omitempty" yaml:"hidden,omitempty"`

	// Flag for overwriting the value with the default on an upgrade
	OverwriteOnUpgrade bool `json:"overwriteOnUpgrade" yaml:"overwriteOnUpgrade"`

	// List of containers affected if this parameter is changed
	AffectedContainers []string `json:"affectedContainers" yaml:"affectedContainers"`
}

func (p Parameter[Type]) GetID() Identifier {
	return p.ID
}

func (p Parameter[Type]) GetName() string {
	return p.Name
}

func (p Parameter[Type]) GetDescription() DynamicProperty[string] {
	return p.Description
}

func (p Parameter[Type]) GetDefault() any {
	return p.Default
}

func (p Parameter[Type]) GetAdvanced() bool {
	return p.Advanced
}

func (p Parameter[Type]) GetDisabled() DynamicProperty[bool] {
	return p.Disabled
}

func (p Parameter[Type]) GetHidden() DynamicProperty[bool] {
	return p.Hidden
}

func (p Parameter[Type]) GetOverwriteOnUpgrade() bool {
	return p.OverwriteOnUpgrade
}

func (p Parameter[Type]) GetAffectedContainers() []string {
	return p.AffectedContainers
}

// Serializes parameter metadata to a map
func (p *Parameter[Type]) serializeImpl() map[string]any {
	props := map[string]any{
		IDKey:                 p.ID,
		NameKey:               p.Name,
		DescriptionKey:        p.Description,
		DefaultKey:            p.Default,
		AdvancedKey:           p.Advanced,
		DisabledKey:           p.Disabled,
		HiddenKey:             p.Hidden,
		OverwriteOnUpgradeKey: p.OverwriteOnUpgrade,
		AffectedContainersKey: p.AffectedContainers,
	}
	return props
}

// DeserializeImpl the parameter metadata from a map
func (p *Parameter[Type]) deserializeImpl(data map[string]any) error {
	// Get the ID
	err := deserializeIdentifier(data, IDKey, &p.ID, false)
	if err != nil {
		return err
	}

	// Get the name
	_, err = deserializeProperty(data, NameKey, &p.Name, false)
	if err != nil {
		return err
	}

	// Get the description
	_, err = deserializeDynamicProperty(data, DescriptionKey, &p.Description, false)
	if err != nil {
		return err
	}

	// Get the advanced flag
	_, err = deserializeProperty(data, AdvancedKey, &p.Advanced, true)
	if err != nil {
		return err
	}

	// Get the disabled flag
	_, err = deserializeDynamicProperty(data, DisabledKey, &p.Disabled, true)
	if err != nil {
		return err
	}

	// Get the hidden flag
	_, err = deserializeDynamicProperty(data, HiddenKey, &p.Hidden, true)
	if err != nil {
		return err
	}

	// Get the overwriteOnUpgrade flag
	_, err = deserializeProperty(data, OverwriteOnUpgradeKey, &p.OverwriteOnUpgrade, false)
	if err != nil {
		return err
	}

	// Get the affectedContainers list
	var affectedContainers []any
	found, err := deserializeProperty(data, AffectedContainersKey, &affectedContainers, true)
	if err != nil {
		return err
	}
	if !found {
		affectedContainers = []any{}
	}
	for _, container := range affectedContainers {
		containerName, ok := container.(string)
		if !ok {
			return fmt.Errorf("invalid affected container name: %T", container)
		}
		p.AffectedContainers = append(p.AffectedContainers, containerName)
	}
	return nil
}

// Underlying implementation for parameter instances
type parameterInstance[Type any] struct {
	// The parameter metadata
	metadata IParameter

	// The current value of the parameter
	Value Type
}

// Gets the metadata for the parameter
func (p *parameterInstance[Type]) GetMetadata() IParameter {
	return p.metadata
}

// Gets the value of the parameter
func (p *parameterInstance[Type]) GetValue() any {
	return p.Value
}

// Gets the value of the parameter as a string
func (p *parameterInstance[Type]) String() string {
	return fmt.Sprintf("%v", p.Value)
}

/// =======================
/// === Bool Parameters ===
/// =======================

// A boolean parameter's metadata
type BoolParameter struct {
	Parameter[bool]
}

func (p BoolParameter) GetType() ParameterType {
	return ParameterType_Bool
}

func (p BoolParameter) Serialize() map[string]any {
	return p.serializeImpl()
}

func (p *BoolParameter) Deserialize(data map[string]any) error {
	err := p.Parameter.deserializeImpl(data)
	if err != nil {
		return err
	}

	// Set the default value
	_, err = deserializeProperty(data, DefaultKey, &p.Default, false)
	if err != nil {
		return err
	}
	return nil
}

func (p *BoolParameter) CreateInstance() IParameterInstance {
	return &boolParameterInstance{
		parameterInstance: parameterInstance[bool]{
			metadata: p,
		},
		Value: p.Default,
	}
}

// Underlying implementation for bool parameter instances
type boolParameterInstance struct {
	parameterInstance[bool]
	Value bool
}

func (i *boolParameterInstance) SetValue(value any) error {
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

const (
	// Field names
	MinValueKey string = "minValue"
	MaxValueKey string = "maxValue"
)

type NumberParameterType interface {
	int64 | uint64 | float64
}

// An integer parameter's metadata
type NumberParameter[Type NumberParameterType] struct {
	Parameter[Type]

	// Minimum value for the parameter
	MinValue Type `json:"minValue,omitempty" yaml:"minValue,omitempty"`

	// Maximum value for the parameter
	MaxValue Type `json:"maxValue,omitempty" yaml:"maxValue,omitempty"`
}

// Underlying implementation for number parameter instances
type numberParameterInstance[Type NumberParameterType] struct {
	parameterInstance[Type]
	Value Type
}

func (p NumberParameter[Type]) Serialize() map[string]any {
	props := p.serializeImpl()
	props[MinValueKey] = p.MinValue
	props[MaxValueKey] = p.MaxValue
	return props
}

func (p *NumberParameter[Type]) Deserialize(data map[string]any) error {
	err := p.Parameter.deserializeImpl(data)
	if err != nil {
		return err
	}

	// Set the min value
	err = p.setProperty(data, MinValueKey, &p.MinValue)
	if err != nil {
		return err
	}

	// Set the max value
	err = p.setProperty(data, MaxValueKey, &p.MaxValue)
	if err != nil {
		return err
	}

	// Set the default value
	err = p.setProperty(data, DefaultKey, &p.Default)
	if err != nil {
		return err
	}
	return nil
}

func (i *numberParameterInstance[Type]) SetValue(value any) error {
	return setNumberProperty(i.metadata.GetID().String(), &i.Value, value)
}

func (p *NumberParameter[Type]) setProperty(data map[string]any, key string, property *Type) error {
	var valueAny any
	_, err := deserializeProperty(data, key, &valueAny, true)
	if err != nil {
		return err
	}
	return setNumberProperty(p.ID.String(), property, valueAny)
}

/// ======================
/// === Int Parameters ===
/// ======================

// An integer parameter's metadata
type IntParameter struct {
	NumberParameter[int64]
}

func (p IntParameter) GetType() ParameterType {
	return ParameterType_Int
}

func (p *IntParameter) CreateInstance() IParameterInstance {
	return &numberParameterInstance[int64]{
		parameterInstance: parameterInstance[int64]{
			metadata: p,
		},
		Value: p.Default,
	}
}

/// =======================
/// === Uint Parameters ===
/// =======================

// An unsigned integer parameter's metadata
type UintParameter struct {
	NumberParameter[uint64]
}

func (p UintParameter) GetType() ParameterType {
	return ParameterType_Uint
}

func (p *UintParameter) CreateInstance() IParameterInstance {
	return &numberParameterInstance[uint64]{
		parameterInstance: parameterInstance[uint64]{
			metadata: p,
		},
		Value: p.Default,
	}
}

/// ========================
/// === Float Parameters ===
/// ========================

// A float parameter's metadata
type FloatParameter struct {
	NumberParameter[float64]
}

func (p FloatParameter) GetType() ParameterType {
	return ParameterType_Float
}

func (p *FloatParameter) CreateInstance() IParameterInstance {
	return &numberParameterInstance[float64]{
		parameterInstance: parameterInstance[float64]{
			metadata: p,
		},
		Value: p.Default,
	}
}

/// =========================
/// === String Parameters ===
/// =========================

const (
	// Field names
	MaxLengthKey string = "maxLength"
	RegexKey     string = "regex"
)

// A string parameter's metadata
type StringParameter struct {
	Parameter[string]

	// The max length of the string
	MaxLength uint64 `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`

	// The pattern for the regular expression the string must match
	Regex string `json:"regex,omitempty" yaml:"regex,omitempty"`
}

func (p *StringParameter) CreateInstance() IParameterInstance {
	return &stringParameterInstance{
		parameterInstance: parameterInstance[string]{
			metadata: p,
		},
		Value: p.Default,
	}
}

type stringParameterInstance struct {
	parameterInstance[string]
	Value string
}

func (p StringParameter) GetType() ParameterType {
	return ParameterType_String
}

func (p StringParameter) Serialize() map[string]any {
	props := p.serializeImpl()
	props[MaxLengthKey] = p.MaxLength
	props[RegexKey] = p.Regex
	return props
}

func (p *StringParameter) Deserialize(data map[string]any) error {
	err := p.Parameter.deserializeImpl(data)
	if err != nil {
		return err
	}

	// Set the max length
	var maxLength any
	_, err = deserializeProperty(data, MaxLengthKey, &maxLength, true)
	if err != nil {
		return err
	}
	err = setNumberProperty(p.ID.String()+"."+MaxLengthKey, &p.MaxLength, maxLength)

	// Set the regex pattern
	_, err = deserializeProperty(data, RegexKey, &p.Regex, true)
	if err != nil {
		return err
	}

	// Set the default value
	_, err = deserializeProperty(data, DefaultKey, &p.Default, false)
	if err != nil {
		return err
	}
	return nil
}

func (i *stringParameterInstance) SetValue(value any) error {
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

const (
	// Field names
	OptionsKey string = "options"
)

type IChoiceParameter interface {
	IParameter
	GetOptions() []IParameterOption
}

type ChoiceParameter[ChoiceType ~string] struct {
	Parameter[ChoiceType]

	// The choices available for the parameter
	Options []ParameterOption[ChoiceType] `json:"options" yaml:"options"`
}

func (p *ChoiceParameter[ChoiceType]) CreateInstance() IParameterInstance {
	return &choiceParameterInstance[ChoiceType]{
		parameterInstance: parameterInstance[ChoiceType]{
			metadata: p,
		},
		Value: p.Default,
	}
}

type choiceParameterInstance[ChoiceType ~string] struct {
	parameterInstance[ChoiceType]
	Value ChoiceType
}

// Gets the type of the parameter
func (p ChoiceParameter[ChoiceType]) GetType() ParameterType {
	return ParameterType_Choice
}

// Gets the options for the choice parameter
func (p ChoiceParameter[ChoiceType]) GetOptions() []IParameterOption {
	options := make([]IParameterOption, len(p.Options))
	for i, option := range p.Options {
		options[i] = option
	}
	return options
}

// Serializes the choice parameter to a map
func (p ChoiceParameter[ChoiceType]) Serialize() map[string]any {
	props := p.serializeImpl()
	options := make([]map[string]any, len(p.Options))
	for i, option := range p.Options {
		options[i] = option.Serialize()
	}
	props[OptionsKey] = options
	return props
}

// Unmarshal the choice parameter from a map
func (p *ChoiceParameter[ChoiceType]) Deserialize(data map[string]any) error {
	err := p.Parameter.deserializeImpl(data)
	if err != nil {
		return err
	}

	// Set the default value
	var defaultVal any
	_, err = deserializeProperty(data, DefaultKey, &defaultVal, false)
	if err != nil {
		return err
	}
	err = p.setProperty(&p.Default, defaultVal)
	if err != nil {
		return err
	}

	// Get the options
	var options []any
	_, err = deserializeProperty(data, OptionsKey, &options, false)
	if err != nil {
		return err
	}

	// Unmarshal the options
	for _, optionData := range options {
		optionDataMap, ok := optionData.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid option data: %T", optionData)
		}
		option := ParameterOption[ChoiceType]{}
		err = option.Deserialize(optionDataMap)
		if err != nil {
			return err
		}
		p.Options = append(p.Options, option)
	}
	return nil
}

func (i *choiceParameterInstance[ChoiceType]) SetValue(value any) error {
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

func (p *ChoiceParameter[ChoiceType]) setProperty(property *ChoiceType, value any) error {
	paramString, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for choice parameter [%s]: %T", p.ID, value)
	}
	*property = ChoiceType(paramString)
	return nil
}

/// =========================
/// === Parameter Options ===
/// =========================

// Interface for a single option for a choice parameter
type IParameterOption interface {
	// The option's name
	GetName() string

	// The description for the option
	GetDescription() DynamicProperty[string]

	// The value for the option
	GetValue() string

	// Flag for disabling the option in the UI, graying it out
	GetDisabled() DynamicProperty[bool]

	// Flag for hiding the option from the UI
	GetHidden() DynamicProperty[bool]
}

// A single option for a choice parameter
type ParameterOption[ChoiceType ~string] struct {
	// The option's name
	Name string `json:"name" yaml:"name"`

	// The description for the option
	Description DynamicProperty[string] `json:"description" yaml:"description"`

	// The value for the option
	Value ChoiceType `json:"value" yaml:"value"`

	// Flag for disabling the option in the UI, graying it out
	Disabled DynamicProperty[bool] `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// Flag for hiding the option from the UI
	Hidden DynamicProperty[bool] `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

// Gets the name of the option
func (o ParameterOption[ChoiceType]) GetName() string {
	return o.Name
}

// Gets the description of the option
func (o ParameterOption[ChoiceType]) GetDescription() DynamicProperty[string] {
	return o.Description
}

// Gets the value of the option
func (o ParameterOption[ChoiceType]) GetValue() string {
	return string(o.Value)
}

// Gets the disabled flag of the option
func (o ParameterOption[ChoiceType]) GetDisabled() DynamicProperty[bool] {
	return o.Disabled
}

// Gets the hidden flag of the option
func (o ParameterOption[ChoiceType]) GetHidden() DynamicProperty[bool] {
	return o.Hidden
}

// Serializes the option to a map
func (o ParameterOption[ChoiceType]) Serialize() map[string]any {
	props := map[string]any{
		NameKey:        o.Name,
		DescriptionKey: o.Description,
		ValueKey:       o.Value,
		DisabledKey:    o.Disabled,
		HiddenKey:      o.Hidden,
	}
	return props
}

// Unmarshal the option from a map
func (o *ParameterOption[ChoiceType]) Deserialize(data map[string]any) error {
	// Get the name
	_, err := deserializeProperty(data, NameKey, &o.Name, false)
	if err != nil {
		return err
	}

	// Get the description
	_, err = deserializeDynamicProperty(data, DescriptionKey, &o.Description, false)
	if err != nil {
		return err
	}

	// Get the value
	_, err = deserializeProperty(data, ValueKey, &o.Value, false)
	if err != nil {
		return err
	}

	// Get the disabled flag
	_, err = deserializeDynamicProperty(data, DisabledKey, &o.Disabled, true)
	if err != nil {
		return err
	}

	// Get the hidden flag
	_, err = deserializeDynamicProperty(data, HiddenKey, &o.Hidden, true)
	if err != nil {
		return err
	}
	return nil
}

/// =====================
/// === Serialization ===
/// =====================

func serializeParameterToMap(p IParameter) map[string]any {
	props := p.Serialize()
	props[TypeKey] = p.GetType()
	return props
}

// Deserializes parameter metadata from a map
func deserializeMapToParameter(serializedParam map[string]any) (IParameter, error) {
	// Get the type
	var paramType string
	_, err := deserializeProperty(serializedParam, TypeKey, &paramType, false)
	if err != nil {
		return nil, err
	}

	// Create the parameter based on the type
	var param IParameter
	switch ParameterType(paramType) {
	case ParameterType_Bool:
		param = &BoolParameter{}
	case ParameterType_Int:
		param = &IntParameter{}
	case ParameterType_Uint:
		param = &UintParameter{}
	case ParameterType_Float:
		param = &FloatParameter{}
	case ParameterType_String:
		param = &StringParameter{}
	case ParameterType_Choice:
		param = &ChoiceParameter[string]{}
	default:
		return nil, fmt.Errorf("invalid parameter type: %s", paramType)
	}

	// Deserialize the parameter
	err = param.Deserialize(serializedParam)
	if err != nil {
		return nil, err
	}
	return param, nil
}

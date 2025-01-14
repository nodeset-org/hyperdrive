package config

const (
	// Field names
	ParametersKey string = "parameters"
	SectionsKey   string = "sections"
)

// Represents the header for a section metadata object
type ISectionHeader interface {
	// Unique ID for referencing the section behind-the-scenes
	GetID() Identifier

	// Name of the section
	GetName() string

	// The description for the section
	GetDescription() DynamicProperty[string]

	// Flag for disabling the section in the UI, graying it out
	GetDisabled() DynamicProperty[bool]

	// Flag for hiding the section from the UI
	GetHidden() DynamicProperty[bool]
}

// Represents a full section metadata object
type ISection interface {
	ISectionHeader

	IMetadataContainer
}

// SectionHeader represents the header of a section in a configuration metadata
type SectionHeader struct {
	// Unique ID for referencing the section behind-the-scenes
	ID Identifier `json:"id" yaml:"id"`

	// Name of the section
	Name string `json:"name" yaml:"name"`

	// The description for the section
	Description DynamicProperty[string] `json:"description" yaml:"description"`

	// Flag for disabling the section in the UI, graying it out
	Disabled DynamicProperty[bool] `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// Flag for hiding the section from the UI
	Hidden DynamicProperty[bool] `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

// Get the unique ID for the section
func (s SectionHeader) GetID() Identifier {
	return s.ID
}

// Get the name of the section
func (s SectionHeader) GetName() string {
	return s.Name
}

// Get the description for the section
func (s SectionHeader) GetDescription() DynamicProperty[string] {
	return s.Description
}

// Get the disabled flag for the section
func (s SectionHeader) GetDisabled() DynamicProperty[bool] {
	return s.Disabled
}

// Get the hidden flag for the section
func (s SectionHeader) GetHidden() DynamicProperty[bool] {
	return s.Hidden
}

/// ====================
/// === Full Section ===
/// ====================

// Section represents a full section in a configuration metadata
type section struct {
	ISectionHeader

	IMetadataContainer
}

/// =====================
/// === Serialization ===
/// =====================

// Serialize the section header to a map
func serializeSectionHeaderToMap(s ISectionHeader) map[string]any {
	props := map[string]any{
		IDKey:          s.GetID(),
		NameKey:        s.GetName(),
		DescriptionKey: s.GetDescription(),
		DisabledKey:    s.GetDisabled(),
		HiddenKey:      s.GetHidden(),
	}
	return props
}

// Serialize the section to a map
func serializeSectionToMap(s ISection) map[string]any {
	// Serialize the header
	props := serializeSectionHeaderToMap(s)
	serializeContainerToMap(s, props)
	return props
}

// Deserialize a section from a map
func deserializeSectionFromMap(data map[string]any) (ISection, error) {
	header := &SectionHeader{}

	// Get the ID
	err := deserializeIdentifier(data, IDKey, &header.ID, false)
	if err != nil {
		return nil, err
	}

	// Get the name
	_, err = deserializeProperty(data, NameKey, &header.Name, false)
	if err != nil {
		return nil, err
	}

	// Get the description
	_, err = deserializeDynamicProperty(data, DescriptionKey, &header.Description, false)
	if err != nil {
		return nil, err
	}

	// Get the disabled flag
	_, err = deserializeDynamicProperty(data, DisabledKey, &header.Disabled, true)
	if err != nil {
		return nil, err
	}

	// Get the hidden flag
	_, err = deserializeDynamicProperty(data, HiddenKey, &header.Hidden, true)
	if err != nil {
		return nil, err
	}

	// Deserialize the parameters and sections
	container, err := deserializeContainerFromMap(data)
	if err != nil {
		return nil, err
	}

	// Create the section
	section := &section{
		ISectionHeader:     header,
		IMetadataContainer: container,
	}
	return section, nil
}

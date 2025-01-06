package config

const (
	// Field names
	ParametersKey string = "parameters"
	SectionsKey   string = "sections"
)

// Represents the header for a section metadata object
type ISectionMetadataHeader interface {
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
type ISectionMetadata interface {
	ISectionMetadataHeader

	IMetadataContainer
}

// SectionMetadataHeader represents the header of a section in a configuration metadata
type SectionMetadataHeader struct {
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
func (s SectionMetadataHeader) GetID() Identifier {
	return s.ID
}

// Get the name of the section
func (s SectionMetadataHeader) GetName() string {
	return s.Name
}

// Get the description for the section
func (s SectionMetadataHeader) GetDescription() DynamicProperty[string] {
	return s.Description
}

// Get the disabled flag for the section
func (s SectionMetadataHeader) GetDisabled() DynamicProperty[bool] {
	return s.Disabled
}

// Get the hidden flag for the section
func (s SectionMetadataHeader) GetHidden() DynamicProperty[bool] {
	return s.Hidden
}

/// ====================
/// === Full Section ===
/// ====================

// SectionMetadata represents a full section in a configuration metadata
type sectionMetadata struct {
	ISectionMetadataHeader

	IMetadataContainer
}

/// =====================
/// === Serialization ===
/// =====================

// Serialize the section header to a map
func serializeSectionMetadataHeaderToMap(s ISectionMetadataHeader) map[string]any {
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
func serializeSectionMetadataToMap(s ISectionMetadata) map[string]any {
	// Serialize the header
	props := serializeSectionMetadataHeaderToMap(s)
	serializeContainerMetadataToMap(s, props)
	return props
}

// Deserialize a section from a map
func deserializeSectionMetadataFromMap(data map[string]any) (ISectionMetadata, error) {
	header := &SectionMetadataHeader{}

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
	container, err := deserializeContainerMetadataFromMap(data)
	if err != nil {
		return nil, err
	}

	// Create the section
	section := &sectionMetadata{
		ISectionMetadataHeader: header,
		IMetadataContainer:     container,
	}
	return section, nil
}

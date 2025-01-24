package config

// A general-purpose section that can be used to create settings of unknown configurations. Use this if you need to explore
// configurations for other Hyperdrive modules dynamically when you don't know their type at compile time.
type SettingsSection struct {
	metadata ISection

	// The parameters in this section
	parameters map[Identifier]IParameterSetting

	// The subsections under this section
	sections map[Identifier]*SettingsSection
}

// Get the metadata for this section
func (s SettingsSection) GetMetadata() ISection {
	return s.metadata
}

// Get the parameter setting with the given ID
func (s SettingsSection) GetParameter(id Identifier) (IParameterSetting, error) {
	param, exists := s.parameters[id]
	if !exists {
		return nil, NewErrorNotFound(id, EntryType_Parameter)
	}
	return param, nil
}

// Get the subsection with the given ID
func (s SettingsSection) GetSection(id Identifier) (*SettingsSection, error) {
	section, exists := s.sections[id]
	if !exists {
		return nil, NewErrorNotFound(id, EntryType_Section)
	}
	return section, nil
}

// Create a new section for configuration settings
func CreateSettingsSection(metadata ISection) *SettingsSection {
	section := &SettingsSection{
		metadata:   metadata,
		parameters: map[Identifier]IParameterSetting{},
		sections:   map[Identifier]*SettingsSection{},
	}

	// Create the parameter instances
	for _, parameter := range metadata.GetParameters() {
		section.parameters[parameter.GetID()] = parameter.CreateSetting()
	}

	// Create the subsection instances
	for _, subsection := range metadata.GetSections() {
		section.sections[subsection.GetID()] = CreateSettingsSection(subsection)
	}

	return section
}

// Serialize the section instance to a map, suitable for marshalling
func (s SettingsSection) Serialize() map[string]any {
	return serializeContainerInstance(s)
}

// Deserialize the section instance from a map
func (s *SettingsSection) Deserialize(instance map[string]any) error {
	return deserializeContainerInstance(s, instance)
}

// Internal method to get the parameters in this configuration instance
func (m SettingsSection) getParameters() map[Identifier]IParameterSetting {
	return m.parameters
}

// Internal method to get the sections in this configuration instance
func (m SettingsSection) getSections() map[Identifier]*SettingsSection {
	return m.sections
}

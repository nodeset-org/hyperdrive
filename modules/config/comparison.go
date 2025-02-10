package config

// Represents the old and new values of a parameter that has had its settings changed
type ParameterDifference struct {
	// The parameter that has changed
	Parameter IParameter

	// Any issues with the parameter
	Error error

	// The old value of the parameter
	OldValue IParameterSetting

	// The new value of the parameter
	NewValue IParameterSetting
}

// Represents the differences between two container instances
type ContainerDifference struct {
	// The parameters that have changed between old and new
	ParameterDifferences []*ParameterDifference

	// The subsections with any changes
	SectionDifferences []*SectionDifference
}

// Represents the differences between two module configuration settings
type ModuleDifference struct {
	*ContainerDifference
}

// Represents the differences between two section instances
type SectionDifference struct {
	*ContainerDifference

	// The section that has changed
	Section ISection

	// Any issues with the section itself
	Error error
}

// Compare the settings between two instances of a module configuration, returning the differences
func CompareSettings(cfg IModuleConfiguration, oldSettings IInstanceContainer, newSettings IInstanceContainer) *ModuleDifference {
	diff := &ModuleDifference{
		ContainerDifference: &ContainerDifference{
			ParameterDifferences: []*ParameterDifference{},
			SectionDifferences:   []*SectionDifference{},
		},
	}
	CompareContainerSettings(cfg, oldSettings, newSettings, diff.ContainerDifference)
	return diff
}

// Compare the settings between two instances of a container, returning the differences
func CompareContainerSettings(info IMetadataContainer, oldSettings IInstanceContainer, newSettings IInstanceContainer, diff *ContainerDifference) {
	// Compare the parameters
	for _, param := range info.GetParameters() {
		paramDiff := &ParameterDifference{
			Parameter: param,
		}
		oldParam, err := oldSettings.GetParameter(param.GetID())
		if err != nil {
			paramDiff.Error = err
		}
		newParam, err := newSettings.GetParameter(param.GetID())
		if err != nil {
			paramDiff.Error = err
		}
		if oldParam.GetValue() != newParam.GetValue() {
			paramDiff.OldValue = oldParam
			paramDiff.NewValue = newParam
			diff.ParameterDifferences = append(diff.ParameterDifferences, paramDiff)
		}
	}

	// Compare the sections
	for _, section := range info.GetSections() {
		sectionDiff := &SectionDifference{
			Section: section,
			ContainerDifference: &ContainerDifference{
				ParameterDifferences: []*ParameterDifference{},
				SectionDifferences:   []*SectionDifference{},
			},
		}
		oldSection, err := oldSettings.GetSection(section.GetID())
		if err != nil {
			sectionDiff.Error = err
		}
		newSection, err := newSettings.GetSection(section.GetID())
		if err != nil {
			sectionDiff.Error = err
		}
		CompareContainerSettings(section, oldSection, newSection, sectionDiff.ContainerDifference)
		if len(sectionDiff.ParameterDifferences) > 0 || len(sectionDiff.SectionDifferences) > 0 {
			diff.SectionDifferences = append(diff.SectionDifferences, sectionDiff)
		}
	}
}

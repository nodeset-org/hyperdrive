package config

// The result of validating a container
type ContainerResult struct {
	// Any issues with the top-level parameters
	ParameterResults []*ParameterValidationResult

	// Any issues within the nested sections
	SectionResults []*SectionValidationResult
}

// The result of validating a parameter - any issues will be returned as a list of errors
type ParameterValidationResult struct {
	// The parameter being validated
	Parameter IParameter

	// Any issues with the parameter
	Errors []error
}

// The result of validating module settings
type ModuleValidationResult struct {
	*ContainerResult
}

// The result of validating module settings
type SectionValidationResult struct {
	*ContainerResult

	// The section being validated
	Section ISection

	// Any issues with the section itself
	Error error
}

// Add a parameter validation result to the container
func (r *SectionValidationResult) AddParameterResult(result *ParameterValidationResult) {
	r.ParameterResults = append(r.ParameterResults, result)
}

// Add a section validation result to the container
func (r *SectionValidationResult) AddSectionResult(result *SectionValidationResult) {
	r.SectionResults = append(r.SectionResults, result)
}

// Validate the configuration settings for a container
func Validate(info IModuleConfiguration, settings IInstanceContainer) *ModuleValidationResult {
	infoResult := &ModuleValidationResult{
		ContainerResult: &ContainerResult{
			ParameterResults: []*ParameterValidationResult{},
			SectionResults:   []*SectionValidationResult{},
		},
	}
	ValidateContainer(info, settings, infoResult.ContainerResult)
	return infoResult
}

// Validate the configuration settings for a container
func ValidateContainer(info IMetadataContainer, settings IInstanceContainer, result *ContainerResult) {
	// Validate the parameters
	for _, param := range info.GetParameters() {
		paramResult := &ParameterValidationResult{
			Parameter: param,
		}
		setting, err := settings.GetParameter(param.GetID())
		if err != nil {
			paramResult.Errors = []error{err}
		} else {
			result := setting.Validate()
			paramResult.Errors = result
		}
		result.ParameterResults = append(result.ParameterResults, paramResult)
	}

	// Validate the sections
	for _, section := range info.GetSections() {
		setting, err := settings.GetSection(section.GetID())
		sectionResult := &SectionValidationResult{
			Section: section,
			ContainerResult: &ContainerResult{
				ParameterResults: []*ParameterValidationResult{},
				SectionResults:   []*SectionValidationResult{},
			},
		}
		if err != nil {
			sectionResult.Error = err
		} else {
			ValidateContainer(section, setting, sectionResult.ContainerResult)
		}
		result.SectionResults = append(result.SectionResults, sectionResult)
	}
}

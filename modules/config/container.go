package config

import (
	"fmt"
	"reflect"
)

var (
	parameterType     = reflect.TypeOf((*IParameter)(nil)).Elem()
	sectionHeaderType = reflect.TypeOf((*ISectionHeader)(nil)).Elem()
	sectionType       = reflect.TypeOf((*ISection)(nil)).Elem()
)

// Interface for deserialized configuration metadata and section metadata that can contain parameters or sections themselves
type IMetadataContainer interface {
	// Get the list of parameters listed in this container
	GetParameters() []IParameter

	// Get the list of sections listed in this container
	GetSections() []ISection
}

// Implementation of the IMetadataContainer interface
type metadataContainer struct {
	// List of parameters in the container
	parameters []IParameter

	// List of sections in the container
	sections []ISection
}

// Get the list of parameters in the container - only populated when deserializing a configuration from JSON
func (c metadataContainer) GetParameters() []IParameter {
	return c.parameters
}

// Get the list of sections in the container - only populated when deserializing a configuration from JSON
func (c metadataContainer) GetSections() []ISection {
	return c.sections
}

// Interface for a container that can contain instances of parameters and sections
type IInstanceContainer interface {
	// Get the parameter instance with the given ID
	GetParameter(id Identifier) (IParameterInstance, error)

	// Get the section instance with the given ID
	GetSection(id Identifier) (*SectionInstance, error)
}

/*
// Get the parameters and sections within a container, either via the interface directly or indirectly via reflection
func getContainerArtifacts(container any) ([]IParameter, []ISection) {
	params := []IParameter{}
	sections := []ISection{}

	// Check if the data explicitly provides its parameters and sections
	dataAsContainer, isContainer := container.(IMetadataContainer)
	if isContainer {
		params = dataAsContainer.GetParameters()
		sections = dataAsContainer.GetSections()
	} else { // Get the parameters and sections via reflection
		// Handle pointers
		dataType := reflect.TypeOf(container)
		dataVal := reflect.ValueOf(container)
		if dataType.Kind() == reflect.Pointer {
			dataType = dataType.Elem()
			dataVal = dataVal.Elem()
		}

		// Ignore if the provided object is not a struct
		if dataType.Kind() != reflect.Struct {
			return nil, nil
		}

		// Parse the fields
		for i := 0; i < dataType.NumField(); i++ {
			fieldVal := dataVal.Field(i)
			fieldType := fieldVal.Type()
			fieldPtr := fieldVal.Addr()
			fieldPtrType := fieldPtr.Type()

			// Handle pointers to the field
			if fieldPtrType.Implements(parameterType) ||
				fieldPtrType.Implements(sectionType) ||
				fieldPtrType.Implements(sectionHeaderType) {
				fieldType = fieldPtrType
				fieldVal = fieldPtr
			}

			if fieldType.Implements(parameterType) {
				// Handle the parameter
				param := fieldVal.Interface().(IParameter)
				params = append(params, param)
			} else if fieldType.Implements(sectionType) {
				// Handle the section
				section := fieldVal.Interface().(ISection)
				sections = append(sections, section)
			} else if fieldType.Implements(sectionHeaderType) {
				childParams, childSections := getContainerArtifacts(fieldVal.Interface())
				if len(childParams) == 0 && len(childSections) == 0 {
					continue
				}
				section := section{
					ISectionHeader: fieldVal.Interface().(ISectionHeader),
					IMetadataContainer: metadataContainer{
						parameters: childParams,
						sections:   childSections,
					},
				}
				sections = append(sections, section)
			}
		}
	}

	return params, sections
}
*/

// Serialize the container into an existing map
func serializeContainerToMap(container IMetadataContainer, existingData map[string]any) {
	//params, sections := getContainerArtifacts(container)
	params := container.GetParameters()
	sections := container.GetSections()

	// Handle the parameters
	paramMaps := []map[string]any{}
	for _, param := range params {
		paramMap := serializeParameterToMap(param)
		paramMaps = append(paramMaps, paramMap)
	}
	existingData[ParametersKey] = paramMaps

	// Handle the sections
	sectionMaps := []map[string]any{}
	for _, section := range sections {
		sectionMap := serializeSectionToMap(section)
		sectionMaps = append(sectionMaps, sectionMap)
	}
	existingData[SectionsKey] = sectionMaps
}

// Deserialize the container from a map
func deserializeContainerFromMap(data map[string]any) (IMetadataContainer, error) {
	container := &metadataContainer{}

	// Handle the parameters
	var parameters []any
	_, err := deserializeProperty(data, ParametersKey, &parameters, false)
	if err != nil {
		return nil, err
	}
	for _, parameterData := range parameters {
		paramMap, ok := parameterData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid parameter data: %T", parameterData)
		}
		parameter, err := deserializeMapToParameter(paramMap)
		if err != nil {
			return nil, err
		}
		container.parameters = append(container.parameters, parameter)
	}

	// Handle subsections
	var subsections []any
	_, err = deserializeProperty(data, SectionsKey, &subsections, false)
	if err != nil {
		return nil, err
	}
	for _, subsectionData := range subsections {
		subsectionMap, ok := subsectionData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid subsection data: %T", subsectionData)
		}
		subsection, err := deserializeSectionFromMap(subsectionMap)
		if err != nil {
			return nil, err
		}
		container.sections = append(container.sections, subsection)
	}

	return container, nil
}

// Interface for deserialized configuration metadata and section metadata that can contain parameters or sections themselves
type iContainerInstance interface {
	// Get the list of parameters listed in this container
	getParameters() map[Identifier]IParameterInstance

	// Get the list of sections listed in this container
	getSections() map[Identifier]*SectionInstance
}

func serializeContainerInstance(container iContainerInstance) map[string]any {
	instanceMap := map[string]any{}
	for paramID, parameter := range container.getParameters() {
		instanceMap[paramID.String()] = parameter.GetValue()
	}
	for sectionID, section := range container.getSections() {
		instanceMap[sectionID.String()] = section.Serialize()
	}
	return instanceMap
}

// Deserialize the section instance from a map
func deserializeContainerInstance(container iContainerInstance, instance map[string]any) error {
	for _, parameter := range container.getParameters() {
		paramID := parameter.GetMetadata().GetID()
		paramData, exists := instance[paramID.String()]
		if !exists {
			return NewErrorNotFound(paramID, EntryType_Parameter)
		}
		err := parameter.SetValue(paramData)
		if err != nil {
			return fmt.Errorf("error setting parameter [%s]: %w", paramID, err)
		}
	}
	for _, section := range container.getSections() {
		sectionID := section.GetMetadata().GetID()
		sectionData, exists := instance[sectionID.String()]
		if !exists {
			return NewErrorNotFound(sectionID, EntryType_Section)
		}
		sectionDataAsMap, ok := sectionData.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid type for section [%s]: %T", sectionID, sectionData)
		}
		err := section.Deserialize(sectionDataAsMap)
		if err != nil {
			return fmt.Errorf("error processing section [%s]: %w", sectionID, err)
		}
	}
	return nil
}

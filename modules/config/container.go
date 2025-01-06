package config

import (
	"fmt"
	"reflect"
)

var (
	parameterType     = reflect.TypeOf((*IParameterMetadata)(nil)).Elem()
	sectionHeaderType = reflect.TypeOf((*ISectionMetadataHeader)(nil)).Elem()
	sectionType       = reflect.TypeOf((*ISectionMetadata)(nil)).Elem()
)

// Interface for deserialized configuration metadata and section metadata that can contain parameters or sections themselves
type IMetadataContainer interface {
	// Get the list of parameters listed in this container
	GetParameters() []IParameterMetadata

	// Get the list of sections listed in this container
	GetSections() []ISectionMetadata
}

// Implementation of the IMetadataContainer interface
type metadataContainer struct {
	// List of parameters in the container
	parameters []IParameterMetadata

	// List of sections in the container
	sections []ISectionMetadata
}

// Get the list of parameters in the container - only populated when deserializing a configuration from JSON
func (c metadataContainer) GetParameters() []IParameterMetadata {
	return c.parameters
}

// Get the list of sections in the container - only populated when deserializing a configuration from JSON
func (c metadataContainer) GetSections() []ISectionMetadata {
	return c.sections
}

// Get the parameters and sections within a container, either via the interface directly or indirectly via reflection
func getContainerArtifacts(container any) ([]IParameterMetadata, []ISectionMetadata) {
	params := []IParameterMetadata{}
	sections := []ISectionMetadata{}

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
				param := fieldVal.Interface().(IParameterMetadata)
				params = append(params, param)
			} else if fieldType.Implements(sectionType) {
				// Handle the section
				section := fieldVal.Interface().(ISectionMetadata)
				sections = append(sections, section)
			} else if fieldType.Implements(sectionHeaderType) {
				childParams, childSections := getContainerArtifacts(fieldVal.Interface())
				if len(childParams) == 0 && len(childSections) == 0 {
					continue
				}
				section := sectionMetadata{
					ISectionMetadataHeader: fieldVal.Interface().(ISectionMetadataHeader),
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

// Serialize the container into an existing map
func serializeContainerMetadataToMap(container any, existingData map[string]any) {
	params, sections := getContainerArtifacts(container)

	// Handle the parameters
	paramMaps := []map[string]any{}
	for _, param := range params {
		paramMap := serializeParameterMetadataToMap(param)
		paramMaps = append(paramMaps, paramMap)
	}
	existingData[ParametersKey] = paramMaps

	// Handle the sections
	sectionMaps := []map[string]any{}
	for _, section := range sections {
		sectionMap := serializeSectionMetadataToMap(section)
		sectionMaps = append(sectionMaps, sectionMap)
	}
	existingData[SectionsKey] = sectionMaps
}

// Serialize the container into an instance
func serializeContainerMetadataToInstance(container any, existingData map[string]any) {
	params, sections := getContainerArtifacts(container)

	// Handle the parameters
	for _, param := range params {
		existingData[param.GetID().String()] = param.GetValueAsAny()
	}

	// Handle the sections
	for _, section := range sections {
		sectionData := map[string]any{}
		serializeContainerMetadataToInstance(section, sectionData)
		existingData[section.GetID().String()] = sectionData
	}
}

// Deserialize the container from a map
func deserializeContainerMetadataFromMap(data map[string]any) (IMetadataContainer, error) {
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
		parameter, err := deserializeMapToParameterMetadata(paramMap)
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
		subsection, err := deserializeSectionMetadataFromMap(subsectionMap)
		if err != nil {
			return nil, err
		}
		container.sections = append(container.sections, subsection)
	}

	return container, nil
}

// Deserialize a container instance into a container metadata to assign the current values
func deserializeContainerInstanceToMetadata(instance map[string]any, container any) error {
	params, sections := getContainerArtifacts(container)

	// Handle the parameters
	for _, param := range params {
		paramID := param.GetID().String()
		paramData, exists := instance[paramID]
		if !exists {
			return fmt.Errorf("missing parameter [%s]", paramID)
		}
		err := param.SetValue(paramData)
		if err != nil {
			return fmt.Errorf("error setting parameter [%s]: %w", paramID, err)
		}
	}

	// Handle the sections
	for _, section := range sections {
		sectionID := section.GetID().String()
		sectionData, exists := instance[sectionID]
		if !exists {
			return fmt.Errorf("missing section [%s]", sectionID)
		}
		sectionDataAsMap, ok := sectionData.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid type for section [%s]: %T", sectionID, sectionData)
		}

		err := deserializeContainerInstanceToMetadata(sectionDataAsMap, section)
		if err != nil {
			return fmt.Errorf("error processing section [%s]: %w", sectionID, err)
		}
	}

	return nil
}

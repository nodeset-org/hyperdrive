package config

import "fmt"

const (
	TemplateKey string = "template"
)

// A property that can be customized at runtime by template logic
type DynamicProperty[Type any] struct {
	// Default value
	Default Type `json:"default" yaml:"default"`

	// Template for customizing the value
	Template string `json:"template,omitempty" yaml:"template,omitempty"`
}

// Unmarshal the dynamic property from a map
func (p *DynamicProperty[Type]) Unmarshal(data map[string]any) error {
	// Get the default value
	_, err := deserializeProperty(data, DefaultKey, &p.Default, false)
	if err != nil {
		return err
	}

	// Get the template
	_, err = deserializeProperty(data, TemplateKey, &p.Template, true)
	if err != nil {
		return err
	}
	return nil
}

// Deserialize a dynamic property from a map
func deserializeDynamicProperty[Type any](data map[string]any, propertyName string, property *DynamicProperty[Type], optional bool) (bool, error) {
	// Get the property as a map
	var value map[string]any
	exists, err := deserializeProperty(data, propertyName, &value, optional)
	if err != nil {
		return false, err
	}
	if !exists && optional {
		return false, nil
	}

	// Unmarshal the property
	err = property.Unmarshal(value)
	if err != nil {
		return true, fmt.Errorf("invalid %s \"%s\"", propertyName, value)
	}
	return true, nil
}

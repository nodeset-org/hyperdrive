package metadata

import (
	"fmt"
	"strings"
)

// Struct to pass into the template engine containing all necessary data and methods for populating a template.
type MetadataDataSource struct {
	CustomFields map[string]string
}

// Create a new TemplateDataSource struct.
func NewMetadataDataSource(
	customFields map[string]string,
) (*MetadataDataSource, error) {
	return &MetadataDataSource{
		CustomFields: customFields,
	}, nil

}

func (t *MetadataDataSource) GetValue(fqpn string) (string, error) {
	if value, exists := t.CustomFields[fqpn]; exists {
		return value, nil
	}
	return "", fmt.Errorf("key not found: %s", fqpn)
}

func (t *MetadataDataSource) GetValueArray(fqpn string, delimiter string) ([]string, error) {
	if value, exists := t.CustomFields[fqpn]; exists {
		return strings.Split(value, delimiter), nil
	}
	return nil, fmt.Errorf("key not found: %s", fqpn)
}

func (t *MetadataDataSource) UseDefault() string {
	// TODO
	return "defaultValue"
}

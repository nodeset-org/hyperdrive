package module

import (
	"fmt"
	"strings"
)

// Struct to pass into the template engine containing all necessary data and methods for populating a template.
type TemplateDataSource struct {
	ModuleConfigDir     string
	ModuleSecretFile    string
	ModuleLogDir        string
	ModuleJwtKeyFile    string
	HyperdriveDaemonUrl string

	CustomFields map[string]string
}

// Create a new TemplateDataSource struct.
func NewTemplateDataSource(
	moduleConfigDir string,
	moduleSecretFile string,
	moduleLogDir string,
	moduleJwtKeyFile string,
	hyperdriveDaemonUrl string,
	customFields map[string]string,
) (*TemplateDataSource, error) {
	return &TemplateDataSource{
		ModuleConfigDir:     moduleConfigDir,
		ModuleSecretFile:    moduleSecretFile,
		ModuleLogDir:        moduleLogDir,
		ModuleJwtKeyFile:    moduleJwtKeyFile,
		HyperdriveDaemonUrl: hyperdriveDaemonUrl,
		CustomFields:        customFields,
	}, nil

}

func (t *TemplateDataSource) GetValue(fqpn string) (string, error) {
	if value, exists := t.CustomFields[fqpn]; exists {
		return value, nil
	}
	return "", fmt.Errorf("key not found: %s", fqpn)
}

func (t *TemplateDataSource) GetValueArray(fqpn string, delimiter string) ([]string, error) {
	if value, exists := t.CustomFields[fqpn]; exists {
		return strings.Split(value, delimiter), nil
	}
	return nil, fmt.Errorf("key not found: %s", fqpn)
}

func (t *TemplateDataSource) UseDefault() string {
	// TODO
	return ""
}

package service

import (
	"fmt"
	"strings"
)

// Struct to pass into the template engine containing all necessary data and methods for populating a template.
type ServiceDataSource struct {
	GetProjectName       string
	ModuleConfigDir      string
	ModuleDataDir        string
	HyperdriveDaemonUrl  string
	HyperdriveJwtKeyFile string

	CustomFields map[string]string
}

// Create a new ServiceDataSource struct.
func NewServiceDataSource(
	getProjectName string,
	moduleConfigDir string,
	moduleDataDir string,
	hyperdriveDaemonUrl string,
	hyperdriveJwtKeyFile string,
	customFields map[string]string,
) (*ServiceDataSource, error) {
	return &ServiceDataSource{
		GetProjectName:       getProjectName,
		ModuleConfigDir:      moduleConfigDir,
		ModuleDataDir:        moduleDataDir,
		HyperdriveDaemonUrl:  hyperdriveDaemonUrl,
		HyperdriveJwtKeyFile: hyperdriveJwtKeyFile,

		CustomFields: customFields,
	}, nil

}

func (t *ServiceDataSource) GetValue(fqpn string) (string, error) {
	if value, exists := t.CustomFields[fqpn]; exists {
		return value, nil
	}
	return "", fmt.Errorf("key not found: %s", fqpn)
}

func (t *ServiceDataSource) GetValueArray(fqpn string, delimiter string) ([]string, error) {
	if value, exists := t.CustomFields[fqpn]; exists {
		return strings.Split(value, delimiter), nil
	}
	return nil, fmt.Errorf("key not found: %s", fqpn)
}

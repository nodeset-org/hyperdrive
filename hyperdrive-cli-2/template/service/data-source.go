package service

// Struct to pass into the template engine containing all necessary data and methods for populating a template.
type ServiceDataSource struct {
	ProjectName          func() string
	ModuleConfigDir      func() string
	ModuleDataDir        func() string
	HyperdriveDaemonUrl  func() string
	HyperdriveJwtKeyFile func() string

	CustomFields map[string]string
}

// Create a new ServiceDataSource struct.
func NewServiceDataSource(
	getProjectName func() string,
	getModuleConfigDir func() string,
	getModuleDataDir func() string,
	getHyperdriveDaemonUrl func() string,
	getHyperdriveJwtKeyFile func() string,
	customFields map[string]string,
) (*ServiceDataSource, error) {
	return &ServiceDataSource{
		ProjectName:          getProjectName,
		ModuleConfigDir:      getModuleConfigDir,
		ModuleDataDir:        getModuleDataDir,
		HyperdriveDaemonUrl:  getHyperdriveDaemonUrl,
		HyperdriveJwtKeyFile: getHyperdriveJwtKeyFile,

		CustomFields: customFields,
	}, nil

}

func (t *ServiceDataSource) GetProjectName() string {
	return t.ProjectName()
}

func (t *ServiceDataSource) GetModuleConfigDir() string {
	return t.ModuleConfigDir()
}

func (t *ServiceDataSource) GetModuleDataDir() string {
	return t.ModuleDataDir()
}

func (t *ServiceDataSource) GetHyperdriveDaemonUrl() string {
	return t.HyperdriveDaemonUrl()
}

func (t *ServiceDataSource) GetHyperdriveJwtKeyFile() string {
	return t.HyperdriveJwtKeyFile()
}

// func (t *ServiceDataSource) GetValue(fqpn string) (string, error) {
// 	if value, exists := t.CustomFields[fqpn]; exists {
// 		return value, nil
// 	}
// 	return "", fmt.Errorf("key not found: %s", fqpn)
// }

// func (t *ServiceDataSource) GetValueArray(fqpn string, delimiter string) ([]string, error) {
// 	if value, exists := t.CustomFields[fqpn]; exists {
// 		return strings.Split(value, delimiter), nil
// 	}
// 	return nil, fmt.Errorf("key not found: %s", fqpn)
// }

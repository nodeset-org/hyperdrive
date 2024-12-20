package adapter

// Struct to pass into the template engine containing all necessary data and methods for populating a template.
type AdapterDataSource struct {
	ProjectName      func() string
	ModuleConfigDir  func() string
	ModuleSecretFile func() string
	ModuleLogDir     func() string
	ModuleJwtKeyFile func() string
}

// Create a new AdapterDataSource struct.
func NewAdapterDataSource(
	getProjectName func() string,
	getModuleConfigDir func() string,
	getModuleSecretFile func() string,
	getModuleLogDir func() string,
	getModuleJwtKeyFile func() string,
	customFields map[string]string,
) (*AdapterDataSource, error) {
	return &AdapterDataSource{
		ProjectName:      getProjectName,
		ModuleConfigDir:  getModuleConfigDir,
		ModuleSecretFile: getModuleSecretFile,
		ModuleLogDir:     getModuleLogDir,
		ModuleJwtKeyFile: getModuleJwtKeyFile,
	}, nil

}

func (a *AdapterDataSource) GetProjectName() string {
	return a.ProjectName()
}

func (a *AdapterDataSource) GetModuleConfigDir() string {
	return a.ModuleConfigDir()
}

func (a *AdapterDataSource) GetModuleSecretFile() string {
	return a.ModuleSecretFile()
}

func (a *AdapterDataSource) GetModuleLogDir() string {
	return a.ModuleLogDir()
}

func (a *AdapterDataSource) GetModuleJwtKeyFile() string {
	return a.ModuleJwtKeyFile()
}

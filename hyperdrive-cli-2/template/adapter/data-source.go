package adapter

// Struct to pass into the template engine containing all necessary data and methods for populating a template.
type AdapterDataSource struct {
	GetProjectName   string
	ModuleConfigDir  string
	ModuleSecretFile string
	ModuleLogDir     string
	ModuleJwtKeyFile string
}

// Create a new AdapterDataSource struct.
func NewAdapterDataSource(
	getProjectName string,
	moduleConfigDir string,
	moduleSecretFile string,
	moduleLogDir string,
	moduleJwtKeyFile string,
	customFields map[string]string,
) (*AdapterDataSource, error) {
	return &AdapterDataSource{
		GetProjectName:   getProjectName,
		ModuleConfigDir:  moduleConfigDir,
		ModuleSecretFile: moduleSecretFile,
		ModuleLogDir:     moduleLogDir,
		ModuleJwtKeyFile: moduleJwtKeyFile,
	}, nil

}

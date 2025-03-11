package management

import (
	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

// Simple status of a Docker container
type ContainerStatus string

const (
	// The container is running
	ContainerStatus_Running ContainerStatus = "running"

	// The container exists, but isn't currently running
	ContainerStatus_Stopped ContainerStatus = "stopped"

	// The container doesn't exist yet
	ContainerStatus_Missing ContainerStatus = "missing"
)

// Binding for a module that's been installed on the system
type ModuleInstallation struct {
	// The full path of the directory the module is installed in
	InstallationPath string

	// The path to the module's descriptor file
	DescriptorPath string

	// The module descriptor
	Descriptor *modules.ModuleDescriptor

	// An error that occurred while loading the module's descriptor, if it couldn't be loaded
	DescriptorLoadError error

	// The path to the module's global adapter compose file
	GlobalAdapterRuntimeFilePath string

	// An error that occured while checking if the module's global adapter container file has been instantiated, if it hasn't
	GlobalAdapterRuntimeFileError error

	// The name of the global adapter container
	GlobalAdapterContainerName string

	// The status of the global adapter container
	GlobalAdapterContainerStatus ContainerStatus

	// The module's configuration if it was loaded properly
	Configuration modconfig.IModuleConfiguration

	// An error that occurred while loading the module's configuration, if it couldn't be loaded
	ConfigurationLoadError error
}

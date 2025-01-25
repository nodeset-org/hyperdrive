package config

import "github.com/nodeset-org/hyperdrive/modules"

// The configuration for a module, along with some module metadata
type ModuleInfo struct {
	// The module's descriptor
	Descriptor modules.ModuleDescriptor

	// The configuration metadata for the module
	Configuration IModuleConfiguration
}

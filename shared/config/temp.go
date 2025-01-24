package config

import "github.com/nodeset-org/hyperdrive/modules"

// This file contains TEMPORARY placeholder functions that exist solely to facilitate development and will eventually be replaced by real implementations prior to release.

func GetModuleAdapterContainerName(descriptor modules.ModuleDescriptor, projectName string) string {
	return projectName + "_" + string(descriptor.Shortcut) + "_adapter"
}

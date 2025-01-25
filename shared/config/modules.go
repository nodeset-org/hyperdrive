package config

import (
	"errors"
)

var (
	// The module descriptor was missing
	ErrNoDescriptor error = errors.New("descriptor file is missing")

	// The module config was not loaded because the module adapter container is missing
	ErrNoAdapterContainer error = errors.New("adapter container is missing")

	// The module config was not loaded because the module adapter is not running
	ErrAdapterContainerOffline error = errors.New("adapter container is offline")
)

// An error that occurs when loading module information
type ModuleInfoLoadError struct {
	// The error thrown by the adapter container while getting the module config
	internalError error
}

// Create a new module info load error
func NewModuleInfoLoadError(err error) ModuleInfoLoadError {
	return ModuleInfoLoadError{internalError: err}
}

// Get the error message for a module config load error
func (e ModuleInfoLoadError) Error() string {
	return "error loading module info: " + e.internalError.Error()
}

// Result of loading a module's info
type ModuleInfoLoadResult struct {
	// The name of the module
	Name string

	// An error that occurred while loading the module's info
	LoadError error
}

// Result of processing a module's config
type ModuleConfigProcessResult struct {
	// An error that occurred at the system level while trying to process the module config
	ProcessError error

	// A list of errors or issues with the module's config that need to be addressed prior to saving
	Issues []string

	// A list of ports that the module will expose on the host machine
	Ports map[string]uint16
}

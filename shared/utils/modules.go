package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/shared"
)

const (
	AdapterContainerName string = "adapter"
)

// Get the compose project name for a module
func GetModuleComposeProjectName(hdProjectName string, descriptor *modules.ModuleDescriptor) string {
	return fmt.Sprintf("%s-%s", hdProjectName, descriptor.Shortcut)
}

// Get the container name for a service
func GetServiceContainerName(hdProjectName string, descriptor *modules.ModuleDescriptor, serviceName string) string {
	modProject := GetModuleComposeProjectName(hdProjectName, descriptor)
	return fmt.Sprintf("%s_%s", modProject, serviceName)
}

// Get the name of the global adapter container for a module
func GetGlobalAdapterContainerName(descriptor *modules.ModuleDescriptor) string {
	return GetModuleComposeProjectName(shared.GlobalAdapterProjectName, descriptor) + "_" + AdapterContainerName
}

// Get the name of the project adapter container for a module
func GetProjectAdapterContainerName(projectName string, descriptor *modules.ModuleDescriptor) string {
	return GetModuleComposeProjectName(projectName, descriptor) + "_" + AdapterContainerName
}

// Get the descriptors for all installed modules
func GetInstalledDescriptors(modulePath string) ([]*modules.ModuleDescriptor, error) {
	// Enumerate the installed modules
	entries, err := os.ReadDir(modulePath)
	if err != nil {
		return nil, fmt.Errorf("error reading module directory: %w", err)
	}

	// Find the modules
	descriptors := []*modules.ModuleDescriptor{}
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		moduleDir := filepath.Join(modulePath, entry.Name())

		// Check if the descriptor exists - this is the key for modules
		var descriptor modules.ModuleDescriptor
		descriptorPath := filepath.Join(moduleDir, modules.DescriptorFilename)
		bytes, err := os.ReadFile(descriptorPath)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("error reading descriptor file \"%s\": %w", descriptorPath, err)
		}

		// Load the descriptor
		err = json.Unmarshal(bytes, &descriptor)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling descriptor: %w", err)
		}
		descriptors = append(descriptors, &descriptor)
	}
	return descriptors, nil
}

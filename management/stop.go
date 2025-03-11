package management

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	hdconfig "github.com/nodeset-org/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/utils"
)

// Stop the Hyperdrive daemon and all services, including module services and project adapters
func (m *HyperdriveManager) StopService(settings *hdconfig.HyperdriveSettings) error {
	// Get the list of installed descriptors
	descriptors, err := utils.GetInstalledDescriptors(m.GetModuleManager().GetModuleSystemDir())
	if err != nil {
		return fmt.Errorf("error getting installed module descriptors: %w", err)
	}

	// Stop each module
	err = stopModules(settings, descriptors)
	if err != nil {
		return fmt.Errorf("error stopping modules: %w", err)
	}

	// Stop the project adapters for each module, along with the Hyperdrive daemon
	composeFiles := []string{} // TODO
	err = stopHdServices(m.Context.UserDirPath, settings.ProjectName, descriptors, composeFiles)
	if err != nil {
		return fmt.Errorf("error stopping project adapters: %w", err)
	}

	// Stop the global adapters
	err = stopGlobalAdapters()
	if err != nil {
		return fmt.Errorf("error stopping global adapters: %w", err)
	}

	return nil
}

// Stop the project adapters for each module, along with the Hyperdrive daemon
func stopHdServices(
	userDir string,
	projectName string,
	descriptors []*modules.ModuleDescriptor,
	composeFiles []string,
) error {
	// Add the project adapters to the compose files
	// TODO: add override files for these too
	for _, descriptor := range descriptors {
		moduleDir := filepath.Join(userDir, shared.ModulesDir, string(descriptor.Name))
		moduleRuntimeDir := filepath.Join(moduleDir, shared.RuntimeDir)
		adapterRuntimePath := filepath.Join(moduleRuntimeDir, "adapter.yml")
		composeFiles = append(composeFiles, adapterRuntimePath)
	}
	args := []string{
		"compose",
		"-p",
		projectName,
	}
	for _, composeFile := range composeFiles {
		args = append(args, "-f", composeFile)
	}
	args = append(args, "stop")
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error stopping project adapters: %w", err)
	}
	return nil
}

// For each project, have the project adapter stop all services.
func stopModules(
	hdSettings *hdconfig.HyperdriveSettings,
	descriptors []*modules.ModuleDescriptor,
) error {
	for _, descriptor := range descriptors {
		modProject := utils.GetModuleComposeProjectName(hdSettings.ProjectName, descriptor)
		err := StopProject(modProject, nil)
		if err != nil {
			return fmt.Errorf("error stopping module [%s]: %w", descriptor.Name, err)
		}
	}
	return nil
}

// For each project, have the project adapter stop the provided services.
// If a module doesn't have any services to stop, it will be skipped.
func stopModuleServices(
	projectName string,
	services map[string][]string,
	descriptors []*modules.ModuleDescriptor,
) error {
	for _, descriptor := range descriptors {
		// Get the list of services to stop
		fqmn := descriptor.GetFullyQualifiedModuleName()
		modServices, exists := services[fqmn]
		if !exists {
			continue
		}

		// Stop the services
		modProject := utils.GetModuleComposeProjectName(projectName, descriptor)
		err := StopProject(modProject, modServices)
		if err != nil {
			return fmt.Errorf("error stopping module [%s]: %w", descriptor.Name, err)
		}
	}
	return nil
}

// Stop the global adapters for each module
func stopGlobalAdapters() error {
	return StopProject(shared.GlobalAdapterProjectName, nil)
}

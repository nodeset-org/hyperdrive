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

// Delete the Hyperdrive daemon and all services, including module services and project adapters
func (m *HyperdriveManager) DownService(settings *hdconfig.HyperdriveSettings, includeVolumes bool) error {
	// Get the list of installed descriptors
	descriptors, err := utils.GetInstalledDescriptors(m.GetModuleManager().GetModuleSystemDir())
	if err != nil {
		return fmt.Errorf("error getting installed module descriptors: %w", err)
	}

	// Delete each module
	err = downModules(settings, descriptors, includeVolumes)
	if err != nil {
		return fmt.Errorf("error deleting modules: %w", err)
	}

	// Delete the project adapters for each module, along with the Hyperdrive daemon
	composeFiles := []string{} // TODO
	err = downHdServices(m.Context.UserDirPath, settings.ProjectName, descriptors, composeFiles, includeVolumes)
	if err != nil {
		return fmt.Errorf("error deleting project adapters: %w", err)
	}

	// Delete the global adapters
	err = downGlobalAdapters(includeVolumes)
	if err != nil {
		return fmt.Errorf("error deleting global adapters: %w", err)
	}

	return nil
}

// Delete the project adapters for each module, along with the Hyperdrive daemon
func downHdServices(
	userDir string,
	projectName string,
	descriptors []*modules.ModuleDescriptor,
	composeFiles []string,
	includeVolumes bool,
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
	args = append(args, "down")
	if includeVolumes {
		args = append(args, "--volumes")
	}
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error deleting project adapters: %w", err)
	}
	return nil
}

// For each project, have the project adapter stop and delete all services.
func downModules(
	hdSettings *hdconfig.HyperdriveSettings,
	descriptors []*modules.ModuleDescriptor,
	includeVolumes bool,
) error {
	for _, descriptor := range descriptors {
		modProject := utils.GetModuleComposeProjectName(hdSettings.ProjectName, descriptor)
		err := DownProject(modProject, includeVolumes)
		if err != nil {
			return fmt.Errorf("error deleting module \"%s\": %w", descriptor.Name, err)
		}
	}
	return nil
}

// Delete the global adapters
func downGlobalAdapters(includeVolumes bool) error {
	return DownProject(shared.GlobalAdapterProjectName, includeVolumes)
}

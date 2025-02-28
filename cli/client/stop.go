package client

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	hdconfig "github.com/nodeset-org/hyperdrive/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/utils"
)

// Stop the Hyperdrive daemon and all services
func (c *HyperdriveClient) StopService(settings *hdconfig.HyperdriveSettings) error {
	// Build the module settings map
	hdCfg := c.GetHyperdriveConfiguration()
	modInfos, moduleSettingsMap, err := createModuleSettingsArtifacts(hdCfg, settings)
	if err != nil {
		return fmt.Errorf("error creating module settings: %w", err)
	}

	modMgr := c.GetModuleManager()

	// Start all of the base services and project module adapters
	composeFiles, err := deployTemplates(c.Context.SystemDirPath, c.Context.UserDirPath, settings)
	if err != nil {
		return fmt.Errorf("error deploying templates: %w", err)
	}
	err = deployModules(modMgr, c.Context.ModulesDir(), settings, moduleSettingsMap, modInfos)
	if err != nil {
		return fmt.Errorf("error deploying modules: %w", err)
	}
	err = startComposeFiles(c.Context.UserDirPath, settings.ProjectName, modInfos, composeFiles)
	if err != nil {
		return fmt.Errorf("error starting project adapters: %w", err)
	}

	// Stop each module
	err = stopModules(modMgr, settings, modInfos)
	if err != nil {
		return fmt.Errorf("error stopping modules: %w", err)
	}

	// Stop the project adapters for each module, along with the Hyperdrive daemon
	err = stopComposeFiles(c.Context.UserDirPath, settings.ProjectName, modInfos, composeFiles)
	if err != nil {
		return fmt.Errorf("error stopping project adapters: %w", err)
	}
	return nil
}

// Stop the project adapters for each module, along with the Hyperdrive daemon
func stopComposeFiles(
	userDir string,
	projectName string,
	infos map[string]*modconfig.ModuleInfo,
	composeFiles []string,
) error {
	// Add the project adapters to the compose files
	// TODO: add override files for these too
	for _, info := range infos {
		moduleDir := filepath.Join(userDir, shared.ModulesDir, string(info.Descriptor.Name))
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
	modMgr *utils.ModuleManager,
	hdSettings *hdconfig.HyperdriveSettings,
	infos map[string]*modconfig.ModuleInfo,
) error {
	for _, info := range infos {
		// Stop the services
		fqmn := info.Descriptor.GetFullyQualifiedModuleName()
		pac, err := modMgr.GetProjectAdapterClient(hdSettings.ProjectName, fqmn)
		if err != nil {
			return fmt.Errorf("error getting project adapter client for module [%s]: %w", info.Descriptor.Name, err)
		}
		err = pac.Stop(context.Background(), hdSettings.ProjectName+"-"+string(info.Descriptor.Shortcut), nil)
		if err != nil {
			return fmt.Errorf("error stopping module [%s]: %w", info.Descriptor.Name, err)
		}
	}
	return nil
}

// For each project, have the project adapter stop the provided services.
// If a module doesn't have any services to stop, it will be skipped.
func stopModuleServices(
	modMgr *utils.ModuleManager,
	projectName string,
	services map[string][]string,
	infos map[string]*modconfig.ModuleInfo,
) error {
	for _, info := range infos {
		// Get the list of services to stop
		fqmn := info.Descriptor.GetFullyQualifiedModuleName()
		modServices, exists := services[fqmn]
		if !exists {
			continue
		}

		// Stop the services
		pac, err := modMgr.GetProjectAdapterClient(projectName, fqmn)
		if err != nil {
			return fmt.Errorf("error getting project adapter client for module [%s]: %w", info.Descriptor.Name, err)
		}
		err = pac.Stop(context.Background(), projectName+"-"+string(info.Descriptor.Shortcut), modServices)
		if err != nil {
			return fmt.Errorf("error stopping module [%s]: %w", info.Descriptor.Name, err)
		}
	}
	return nil
}

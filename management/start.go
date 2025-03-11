package management

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	hdconfig "github.com/nodeset-org/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive/management/template"
	"github.com/nodeset-org/hyperdrive/modules"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
)

func (m *HyperdriveManager) StartService(currentSettings *hdconfig.HyperdriveSettings, pendingSettings *hdconfig.HyperdriveSettings) error {
	if !m.modulesLoaded {
		return fmt.Errorf("modules have not been loaded yet")
	}
	hdCfg := m.GetHyperdriveConfiguration()
	modMgr := m.GetModuleManager()

	// If there are pending settings, stop the services that need to be restarted first
	if pendingSettings != nil {
		// See if this had a project name change, which requires taking down the old project
		if currentSettings.ProjectName != pendingSettings.ProjectName {
			err := m.DownService(currentSettings, false)
			if err != nil {
				return fmt.Errorf("error taking down old project: %w", err)
			}
		} else {
			// Make sure the project adapters are all running
			modInfos, moduleSettingsMap, err := createModuleSettingsArtifacts(hdCfg, currentSettings)
			if err != nil {
				return fmt.Errorf("error creating module settings: %w", err)
			}
			err = m.StartProjectAdapters(currentSettings, modInfos, moduleSettingsMap)
			if err != nil {
				return fmt.Errorf("error starting project adapters: %w", err)
			}

			// Stop the services that need to be restarted
			services := map[string][]string{}
			for fqmn := range modInfos {
				pendingModSettings, exists := pendingSettings.Modules[fqmn]
				if !exists {
					continue
				}
				services[fqmn] = pendingModSettings.Restart
			}
			descriptors := []*modules.ModuleDescriptor{}
			for _, info := range modInfos {
				descriptors = append(descriptors, info.Descriptor)
			}
			err = stopModuleServices(currentSettings.ProjectName, services, descriptors)
			if err != nil {
				return fmt.Errorf("error stopping modules: %w", err)
			}
		}

		// Use the pending settings for the rest of the process
		currentSettings = pendingSettings
	}

	// Build the pending module settings map
	modInfos, moduleSettingsMap, err := createModuleSettingsArtifacts(hdCfg, currentSettings)
	if err != nil {
		return fmt.Errorf("error creating module settings: %w", err)
	}

	// Start all of the base services and project module adapters
	err = m.StartProjectAdapters(currentSettings, modInfos, moduleSettingsMap)
	if err != nil {
		return fmt.Errorf("error starting project adapters: %w", err)
	}

	// Commit the pending settings now
	if pendingSettings != nil {
		// Save the settings
		for _, mod := range pendingSettings.Modules {
			mod.Restart = nil
		}
		err = m.SavePendingSettings(pendingSettings)
		if err != nil {
			return fmt.Errorf("error updating pending settings: %w", err)
		}
		err = m.CommitPendingSettings(true)
		if err != nil {
			return fmt.Errorf("error committing pending settings: %w", err)
		}
	}

	// Set the settings for each module
	for _, info := range modInfos {
		pac, err := modMgr.GetProjectAdapterClient(currentSettings.ProjectName, info.Descriptor.GetFullyQualifiedModuleName())
		if err != nil {
			return fmt.Errorf("error getting project adapter client for module [%s]: %w", info.Descriptor.Name, err)
		}
		err = pac.SetSettings(context.Background(), currentSettings)
		if err != nil {
			return fmt.Errorf("error saving settings for module [%s]: %w", info.Descriptor.Name, err)
		}
	}

	// Start the service
	err = startModules(modMgr, currentSettings, modInfos)
	if err != nil {
		return fmt.Errorf("error starting modules: %w", err)
	}
	return nil
}

func (m *HyperdriveManager) StartProjectAdapters(settings *hdconfig.HyperdriveSettings, modInfos map[string]*modconfig.ModuleInfo, moduleSettingsMap map[string]*modconfig.ModuleSettings) error {
	// Deploy the base HD templates
	composeFiles, err := deployTemplates(m.Context.SystemDirPath, m.Context.UserDirPath, settings)
	if err != nil {
		return fmt.Errorf("error deploying templates: %w", err)
	}

	// Remove modules that are not enabled / missing from the info map
	enabledModules := map[string]*modconfig.ModuleInfo{}
	for fqmn, info := range modInfos {
		_, exists := moduleSettingsMap[fqmn]
		if exists {
			enabledModules[fqmn] = info
		}
	}

	// Deploy the modules
	err = deployModules(m.modMgr, m.Context.ModulesDir(), settings, moduleSettingsMap, enabledModules)
	if err != nil {
		return fmt.Errorf("error deploying modules: %w", err)
	}

	// Start the project
	err = startComposeFiles(m.Context.UserDirPath, settings.ProjectName, enabledModules, composeFiles)
	if err != nil {
		return fmt.Errorf("error starting project adapters: %w", err)
	}
	return nil
}

// Get the info and dynamic settings for all enabled modules (TEMP)
// Assumes the modules are already loaded (LoadModules)
func createModuleSettingsArtifacts(
	hdCfg *hdconfig.HyperdriveConfig,
	settings *hdconfig.HyperdriveSettings,
) (
	map[string]*modconfig.ModuleInfo,
	map[string]*modconfig.ModuleSettings,
	error,
) {
	modInfos := map[string]*modconfig.ModuleInfo{}
	moduleSettingsMap := map[string]*modconfig.ModuleSettings{}
	for fqmn, modInfo := range hdCfg.Modules {
		settings, exists := settings.Modules[fqmn]
		if !exists {
			continue
		}
		if !settings.Enabled {
			continue
		}
		moduleSettings := modconfig.CreateModuleSettings(modInfo.Configuration)
		err := moduleSettings.CopySettingsFromKnownType(settings.Settings)
		if err != nil {
			return nil, nil, fmt.Errorf("error loading settings for module [%s]: %w", fqmn, err)
		}
		moduleSettingsMap[fqmn] = moduleSettings
		modInfos[fqmn] = modInfo
	}
	return modInfos, moduleSettingsMap, nil
}

// Deploy the main Hyperdrive containers and the user's override files
func deployTemplates(systemDir string, userDir string, settings *hdconfig.HyperdriveSettings) ([]string, error) {
	templateSourceDir := filepath.Join(systemDir, shared.TemplatesDir)
	runtimeDir := filepath.Join(userDir, shared.RuntimeDir)
	overrideSourceDir := filepath.Join(systemDir, shared.OverrideDir)
	overrideTargetDir := filepath.Join(userDir, shared.OverrideDir)
	//extraScrapeJobsDir := filepath.Join(userDir, shared.ExtraScrapeJobsDir)

	// Prep the override folder
	err := copyOverrideFiles(overrideSourceDir, overrideTargetDir)
	if err != nil {
		return nil, fmt.Errorf("error copying override files: %w", err)
	}

	// Remove the obsolete Docker Compose version from the overrides
	err = removeComposeVersion(overrideTargetDir)
	if err != nil {
		return nil, fmt.Errorf("error removing obsolete Docker Compose version from overrides: %w", err)
	}

	// Clear out the runtime folder and remake it
	err = os.RemoveAll(runtimeDir)
	if err != nil {
		return nil, fmt.Errorf("error deleting runtime folder [%s]: %w", runtimeDir, err)
	}
	err = os.Mkdir(runtimeDir, 0775)
	if err != nil {
		return nil, fmt.Errorf("error creating runtime folder [%s]: %w", runtimeDir, err)
	}

	// Make the extra scrape jobs folder
	/*
		err = os.MkdirAll(extraScrapeJobsDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating extra-scrape-jobs folder: %w", err)
		}
	*/

	composePaths := template.ComposePaths{
		RuntimePath:  runtimeDir,
		TemplatePath: templateSourceDir,
		OverridePath: overrideTargetDir,
	}

	// Read and substitute the templates
	deployedContainers := []string{}

	// These containers always run
	toDeploy := []string{
		//string(config.ContainerID_Daemon),
	}

	// Deploy main containers
	for _, containerName := range toDeploy {
		containers, err := composePaths.File(string(containerName)).Write(settings)
		if err != nil {
			return nil, fmt.Errorf("could not create %s container definition: %w", containerName, err)
		}
		deployedContainers = append(deployedContainers, containers...)
	}

	return toDeploy, nil
}

// Deploy the modules for the project, instantiating templates and scaffolding their folder structure
func deployModules(
	modMgr *ModuleManager,
	moduleInstallDir string,
	hdSettings *hdconfig.HyperdriveSettings,
	moduleSettingsMap map[string]*modconfig.ModuleSettings,
	infos map[string]*modconfig.ModuleInfo,
) error {
	for _, module := range infos {
		info := infos[module.Descriptor.GetFullyQualifiedModuleName()]
		err := modMgr.DeployModule(moduleInstallDir, hdSettings, moduleSettingsMap, info)
		if err != nil {
			return fmt.Errorf("error deploying module [%s]: %w", module.Descriptor.GetFullyQualifiedModuleName(), err)
		}
	}
	return nil
}

// Start the project adapters for each module, along with the Hyperdrive daemon
func startComposeFiles(
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
	if len(composeFiles) == 0 {
		// If there are no compose files at all, nothing to do
		return nil
	}

	// Start the project
	args = append(args, "up", "-d", "--remove-orphans", "--quiet-pull")
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error starting project adapters: %w", err)
	}
	return nil
}

// For each project, have the project adapter start the module
func startModules(
	modMgr *ModuleManager,
	hdSettings *hdconfig.HyperdriveSettings,
	infos map[string]*modconfig.ModuleInfo,
) error {
	for _, info := range infos {
		pac, err := modMgr.GetProjectAdapterClient(hdSettings.ProjectName, info.Descriptor.GetFullyQualifiedModuleName())
		if err != nil {
			return fmt.Errorf("error getting project adapter client for module [%s]: %w", info.Descriptor.Name, err)
		}
		err = pac.Start(context.Background(), hdSettings)
		if err != nil {
			return fmt.Errorf("error starting module [%s]: %w", info.Descriptor.Name, err)
		}
	}
	return nil
}

// Make sure the override files have all been copied to the local user dir
func copyOverrideFiles(sourceDir string, targetDir string) error {
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating override folder: %w", err)
	}

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("error enumerating override source folder: %w", err)
	}

	// Copy any override files that don't exist in the local user directory
	for _, file := range files {
		filename := file.Name()
		targetPath := filepath.Join(targetDir, filename)
		if file.IsDir() {
			// Recurse
			srcPath := filepath.Join(sourceDir, file.Name())
			err = copyOverrideFiles(srcPath, targetPath)
			if err != nil {
				return err
			}
		}

		_, err := os.Stat(targetPath)
		if !os.IsNotExist(err) {
			// Ignore files that already exist
			continue
		}

		// Read the source
		srcPath := filepath.Join(sourceDir, filename)
		contents, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("error reading override file [%s]: %w", srcPath, err)
		}

		// Write a copy to the user dir
		err = os.WriteFile(targetPath, contents, 0644)
		if err != nil {
			return fmt.Errorf("error writing local override file [%s]: %w", targetPath, err)
		}
	}
	return nil
}

// Remove the obsolete Docker Compose version from each compose file in the target directory
func removeComposeVersion(targetDir string) error {
	files, err := os.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("error enumerating folder [%s]: %w", targetDir, err)
	}

	// Copy any override files that don't exist in the local user directory
	for _, file := range files {
		filename := file.Name()
		targetPath := filepath.Join(targetDir, filename)
		if file.IsDir() {
			// Recurse
			subdir := filepath.Join(targetDir, file.Name())
			err = removeComposeVersion(subdir)
			if err != nil {
				return err
			}
		}

		// Ignore it if it's not a YAML file
		if filepath.Ext(filename) != ".yml" {
			continue
		}

		// Read the source
		contents, err := os.ReadFile(targetPath)
		if err != nil {
			return fmt.Errorf("error reading file [%s]: %w", targetPath, err)
		}

		// Remove the version field, accounting for both Windows and Unix line endings
		newContents := bytes.ReplaceAll(contents, []byte("\r\nversion: \"3.7\""), []byte("\r\n"))
		newContents = bytes.ReplaceAll(newContents, []byte("\nversion: \"3.7\""), []byte("\n"))

		// Write the updated contents if they differ
		if len(newContents) != len(contents) {
			err = os.WriteFile(targetPath, newContents, 0644)
			if err != nil {
				return fmt.Errorf("error updating file [%s]: %w", targetPath, err)
			}
		}
	}
	return nil
}

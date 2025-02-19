package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	cliutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/urfave/cli/v2"
)

func startService(c *cli.Context) error {
	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Load the modules
	modLoadResults, err := hd.LoadModules()
	if err != nil {
		return fmt.Errorf("error loading modules: %w", err)
	}
	failures := false
	for _, result := range modLoadResults {
		if result.LoadError != nil {
			fmt.Printf("WARNING: Module %s failed to load: %s\n", result.Info.Descriptor.Name, result.LoadError.Error())
			failures = true
		}
	}
	if failures {
		fmt.Println("The above modules will be disabled if you proceed until their problems are resolved.")
		if !(c.Bool(cliutils.YesFlag.Name) || cliutils.Confirm("Are you sure you want to continue?")) {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Load the config, checking to see if it's new (hasn't been installed before)
	hdCfg := hd.GetHyperdriveConfiguration()
	settings, isNew, err := hd.LoadMainSettingsFile()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}
	if isNew {
		fmt.Println("Hyperdrive has not been configured yet. Please run 'hyperdrive service configure' first.")
		return nil
	}

	// Disable modules that failed to load
	for _, result := range modLoadResults {
		if result.LoadError == nil {
			continue
		}
		modInstance, exists := settings.Modules[result.Info.Descriptor.GetFullyQualifiedModuleName()]
		if !exists {
			continue
		}
		modInstance.Enabled = false
	}

	// Build the module settings map
	modInfos, moduleSettingsMap, err := createModuleSettingsArtifacts(hdCfg, settings)
	if err != nil {
		return fmt.Errorf("error creating module settings: %w", err)
	}

	// Start the service
	err = startImpl(hd.GetModuleManager(), hd.Context.UserDirPath, hd.Context.ModulesDir(), settings, moduleSettingsMap, modInfos)
	if err != nil {
		return fmt.Errorf("error starting service: %w", err)
	}
	return nil
}

// Get the info and dynamic settings for all enabled modules (TEMP)
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

// Starts the service
func startImpl(
	modMgr *utils.ModuleManager,
	userDir string,
	moduleInstallDir string,
	hdSettings *hdconfig.HyperdriveSettings,
	moduleSettingsMap map[string]*modconfig.ModuleSettings,
	infos map[string]*modconfig.ModuleInfo,
) error {
	err := deployModules(modMgr, moduleInstallDir, hdSettings, moduleSettingsMap, infos)
	if err != nil {
		return fmt.Errorf("error deploying modules: %w", err)
	}
	err = startProjectAdapters(userDir, hdSettings.ProjectName, infos)
	if err != nil {
		return fmt.Errorf("error starting project adapters: %w", err)
	}
	err = startModules(modMgr, hdSettings, infos)
	if err != nil {
		return fmt.Errorf("error starting modules: %w", err)
	}
	return nil
}

// Deploy the modules for the project, instantiating templates and scaffolding their folder structure
func deployModules(
	modMgr *utils.ModuleManager,
	moduleInstallDir string,
	hdSettings *hdconfig.HyperdriveSettings,
	moduleSettingsMap map[string]*modconfig.ModuleSettings,
	infos map[string]*modconfig.ModuleInfo,
) error {
	for _, info := range infos {
		err := modMgr.DeployModule(moduleInstallDir, hdSettings, moduleSettingsMap, info)
		if err != nil {
			return fmt.Errorf("error deploying module [%s]: %w", info.Descriptor.GetFullyQualifiedModuleName(), err)
		}
	}
	return nil
}

// Start the project adapters for each module
func startProjectAdapters(
	userDir string,
	projectName string,
	infos map[string]*modconfig.ModuleInfo,
) error {
	projectAdapterFiles := []string{}
	for _, info := range infos {
		moduleDir := filepath.Join(userDir, shared.ModulesDir, string(info.Descriptor.Name))
		moduleRuntimeDir := filepath.Join(moduleDir, shared.RuntimeDir)
		adapterRuntimePath := filepath.Join(moduleRuntimeDir, "adapter.yml")
		projectAdapterFiles = append(projectAdapterFiles, adapterRuntimePath)
	}
	args := []string{"compose", "-p", projectName}
	for _, file := range projectAdapterFiles {
		args = append(args, "-f", file)
	}
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
	modMgr *utils.ModuleManager,
	hdSettings *hdconfig.HyperdriveSettings,
	infos map[string]*modconfig.ModuleInfo,
) error {
	for _, info := range infos {
		pac, err := modMgr.GetProjectAdapterClient(hdSettings.ProjectName, info.Descriptor.GetFullyQualifiedModuleName())
		if err != nil {
			return fmt.Errorf("error getting project adapter client for module [%s]: %w", info.Descriptor.Name, err)
		}
		err = pac.Start(context.Background(), hdSettings, hdSettings.ProjectName+"-"+string(info.Descriptor.Shortcut))
		if err != nil {
			return fmt.Errorf("error starting module [%s]: %w", info.Descriptor.Name, err)
		}
	}
	return nil
}

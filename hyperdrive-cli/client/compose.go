package client

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/mitchellh/go-homedir"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client/template"
	"github.com/rocket-pool/node-manager-core/config"
)

// Build a docker compose command
func (c *HyperdriveClient) compose(composeFiles []string, args string) (string, error) {
	// Get the expanded config path
	expandedConfigPath, err := homedir.Expand(c.Context.ConfigPath)
	if err != nil {
		return "", err
	}

	// Load config
	cfg, isNew, err := c.LoadConfig()
	if err != nil {
		return "", err
	}
	if isNew {
		return "", fmt.Errorf("settings file not found. Please run `hyperdrive service config` to set up Hyperdrive before starting it")
	}

	// Check config
	if cfg.Hyperdrive.ClientMode.Value == config.ClientMode_Unknown {
		return "", fmt.Errorf("you haven't selected local or external mode for your clients yet.\nPlease run 'hyperdrive service config' before running this command")
	} else if cfg.Hyperdrive.IsLocalMode() && cfg.Hyperdrive.LocalExecutionClient.ExecutionClient.Value == config.ExecutionClient_Unknown {
		return "", errors.New("no Execution Client selected. Please run 'hyperdrive service config' before running this command")
	}
	if cfg.Hyperdrive.IsLocalMode() && cfg.Hyperdrive.LocalBeaconClient.BeaconNode.Value == config.BeaconNode_Unknown {
		return "", errors.New("no Beacon Node selected. Please run 'hyperdrive service config' before running this command")
	}

	// Make sure the external IP is loaded
	cfg.LoadExternalIP()

	// Deploy the templates and run environment variable substitution on them
	deployedContainers, err := c.deployTemplates(cfg, expandedConfigPath)
	if err != nil {
		return "", fmt.Errorf("error deploying Docker templates: %w", err)
	}

	// Include all of the relevant docker compose definition files
	composeFileFlags := []string{}
	for _, container := range deployedContainers {
		composeFileFlags = append(composeFileFlags, fmt.Sprintf("-f %s", shellescape.Quote(container)))
	}
	for _, container := range composeFiles {
		composeFileFlags = append(composeFileFlags, fmt.Sprintf("-f %s", shellescape.Quote(container)))
	}

	// Return command
	return fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker compose --project-directory %s %s %s", cfg.Hyperdrive.ProjectName.Value, shellescape.Quote(expandedConfigPath), strings.Join(composeFileFlags, " "), args), nil
}

// Deploys all of the appropriate docker compose template files and provisions them based on the provided configuration
func (c *HyperdriveClient) deployTemplates(cfg *GlobalConfig, hyperdriveDir string) ([]string, error) {
	// Prep the override folder
	overrideFolder := filepath.Join(hyperdriveDir, overrideDir)
	err := copyOverrideFiles(c.Context.OverrideSourceDir, overrideFolder)
	if err != nil {
		return []string{}, fmt.Errorf("error copying override files: %w", err)
	}

	// Remove the obsolete Docker Compose version from the overrides
	err = removeComposeVersion(overrideFolder)
	if err != nil {
		return nil, fmt.Errorf("error removing obsolete Docker Compose version from overrides: %w", err)
	}

	// Clear out the runtime folder and remake it
	runtimeFolder := filepath.Join(hyperdriveDir, runtimeDir)
	err = os.RemoveAll(runtimeFolder)
	if err != nil {
		return []string{}, fmt.Errorf("error deleting runtime folder [%s]: %w", runtimeFolder, err)
	}
	err = os.Mkdir(runtimeFolder, 0775)
	if err != nil {
		return []string{}, fmt.Errorf("error creating runtime folder [%s]: %w", runtimeFolder, err)
	}

	// Make the extra scrape jobs folder
	extraScrapeJobsFolder := filepath.Join(hyperdriveDir, extraScrapeJobsDir)
	err = os.MkdirAll(extraScrapeJobsFolder, 0755)
	if err != nil {
		return []string{}, fmt.Errorf("error creating extra-scrape-jobs folder: %w", err)
	}

	composePaths := template.ComposePaths{
		RuntimePath:  runtimeFolder,
		TemplatePath: c.Context.TemplatesDir,
		OverridePath: overrideFolder,
	}

	// Read and substitute the templates
	deployedContainers := []string{}

	// These containers always run
	toDeploy := []config.ContainerID{
		config.ContainerID_Daemon,
	}

	// Check if we are running the Execution Layer locally
	if cfg.Hyperdrive.IsLocalMode() {
		toDeploy = append(toDeploy, config.ContainerID_ExecutionClient)
		toDeploy = append(toDeploy, config.ContainerID_BeaconNode)
	}

	// Check the metrics containers
	if cfg.Hyperdrive.Metrics.EnableMetrics.Value {
		toDeploy = append(toDeploy,
			config.ContainerID_Grafana,
			config.ContainerID_Exporter,
			config.ContainerID_Prometheus,
		)
	}

	// Check if we are running the MEV-Boost container locally
	if cfg.Hyperdrive.MevBoost.Enable.Value && cfg.Hyperdrive.MevBoost.Mode.Value == config.ClientMode_Local {
		toDeploy = append(toDeploy, config.ContainerID_MevBoost)
	}

	// Deploy main containers
	for _, containerName := range toDeploy {
		containers, err := composePaths.File(string(containerName)).Write(cfg)
		if err != nil {
			return []string{}, fmt.Errorf("could not create %s container definition: %w", containerName, err)
		}
		deployedContainers = append(deployedContainers, containers...)
	}

	// Deploy modules
	for _, module := range cfg.GetAllModuleConfigs() {
		if module.IsEnabled() {
			deployedContainers, err = c.composeModule(cfg, module, hyperdriveDir, deployedContainers)
			if err != nil {
				return nil, err
			}
		}
	}
	return deployedContainers, nil
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

// Handle composing for modules
func (c *HyperdriveClient) composeModule(global *GlobalConfig, module hdconfig.IModuleConfig, hyperdriveDir string, deployedContainers []string) ([]string, error) {
	moduleName := module.GetModuleName()
	composePaths := template.ComposePaths{
		RuntimePath:  filepath.Join(hyperdriveDir, runtimeDir, hdconfig.ModulesName, moduleName),
		TemplatePath: filepath.Join(c.Context.TemplatesDir, hdconfig.ModulesName, moduleName),
		OverridePath: filepath.Join(hyperdriveDir, overrideDir, hdconfig.ModulesName, moduleName),
	}

	// These containers always run
	toDeploy := module.GetContainersToDeploy()

	// Make the modules folder
	err := os.MkdirAll(composePaths.RuntimePath, 0775)
	if err != nil {
		return []string{}, fmt.Errorf("error creating modules runtime folder (%s): %w", composePaths.RuntimePath, err)
	}

	// Deploy the container templates
	for _, containerName := range toDeploy {
		containers, err := composePaths.File(string(containerName)).Write(global)
		if err != nil {
			return []string{}, fmt.Errorf("could not create %s container definition: %w", containerName, err)
		}
		deployedContainers = append(deployedContainers, containers...)
	}

	// Check if the module has a Prometheus config
	prometheusConfigFilename := modulePrometheusSd + template.TemplateSuffix
	_, err = os.Stat(filepath.Join(composePaths.TemplatePath, prometheusConfigFilename))
	if os.IsNotExist(err) {
		return deployedContainers, nil
	}

	// Make the modules dir
	modulesDir := filepath.Join(hyperdriveDir, metricsDir, hdconfig.ModulesName)
	err = os.MkdirAll(modulesDir, metricsDirMode)
	if err != nil {
		return []string{}, fmt.Errorf("error creating metrics module directory [%s]: %w", modulesDir, err)
	}

	// Deploy the Prometheus config
	t := template.Template{
		Src: filepath.Join(composePaths.TemplatePath, prometheusConfigFilename),
		Dst: filepath.Join(modulesDir, moduleName+template.ComposeFileSuffix),
	}
	err = t.Write(global)
	if err != nil {
		return []string{}, fmt.Errorf("could not write module [%s] Prometheus config: %w", moduleName, err)
	}

	return deployedContainers, nil
}

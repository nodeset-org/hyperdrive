package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/blang/semver/v4"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client/template"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	debugColor                    color.Attribute = color.FgYellow
	nethermindPruneStarterCommand string          = "DELETE_ME"
	nethermindAdminUrl            string          = "http://127.0.0.1:7434"

	templatesDir       string = "/usr/share/hyperdrive/templates"
	overrideSourceDir  string = "/usr/share/hyperdrive/override"
	overrideDir        string = "override"
	runtimeDir         string = "runtime"
	extraScrapeJobsDir string = "extra-scrape-jobs"
)

// Install Hyperdrive
func (c *HyperdriveClient) InstallService(verbose bool, noDeps bool, version string, path string, useLocalInstaller bool) error {
	// Get installation script flags
	flags := []string{
		"-v", shellescape.Quote(version),
	}
	if path != "" {
		flags = append(flags, fmt.Sprintf("-p %s", shellescape.Quote(path)))
	}
	if noDeps {
		flags = append(flags, "-d")
	}

	var script []byte
	if useLocalInstaller {
		// Make sure it exists
		_, err := os.Stat(InstallerName)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("local install script [%s] does not exist", InstallerName)
		}
		if err != nil {
			return fmt.Errorf("error checking install script [%s]: %w", InstallerName, err)
		}

		// Read it
		script, err = os.ReadFile(InstallerName)
		if err != nil {
			return fmt.Errorf("error reading local install script [%s]: %w", InstallerName, err)
		}

		// Set the "local mode" flag
		flags = append(flags, "-l")
	} else {
		// Download the installation script
		resp, err := http.Get(fmt.Sprintf(InstallerURL, version))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected http status downloading installation script: %d", resp.StatusCode)
		}

		// Sanity check that the script octet length matches content-length
		script, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if fmt.Sprint(len(script)) != resp.Header.Get("content-length") {
			return fmt.Errorf("downloaded script length %d did not match content-length header %s", len(script), resp.Header.Get("content-length"))
		}
	}

	// Get the escalation command
	escalationCmd, err := c.getEscalationCommand()
	if err != nil {
		return fmt.Errorf("error getting escalation command: %w", err)
	}

	// Initialize installation command
	cmd := c.newCommand(fmt.Sprintf("%s sh -s -- %s", escalationCmd, strings.Join(flags, " ")))

	// Pass the script to sh via its stdin fd
	cmd.SetStdin(bytes.NewReader(script))

	// Get command output pipes
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Print progress from stdout
	go (func() {
		scanner := bufio.NewScanner(cmdOut)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	})()

	// Read command & error output from stderr; render in verbose mode
	var errMessage string
	go (func() {
		c := color.New(debugColor)
		scanner := bufio.NewScanner(cmdErr)
		for scanner.Scan() {
			errMessage = scanner.Text()
			if verbose {
				_, _ = c.Println(scanner.Text())
			}
		}
	})()

	// Run command and return error output
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not install Hyperdrive service: %s", errMessage)
	}
	return nil
}

// Start the Hyperdrive service
func (c *HyperdriveClient) StartService(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "up -d --remove-orphans --quiet-pull")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Pause the Hyperdrive service
func (c *HyperdriveClient) PauseService(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "stop")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Stop the Hyperdrive service
func (c *HyperdriveClient) StopService(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "down -v")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Stop Hyperdrive and remove the config folder
func (c *HyperdriveClient) TerminateService(composeFiles []string, configPath string) error {
	// Get the command to run with root privileges
	rootCmd, err := c.getEscalationCommand()
	if err != nil {
		return fmt.Errorf("could not get privilege escalation command: %w", err)
	}

	// Terminate the Docker containers
	cmd, err := c.compose(composeFiles, "down -v")
	if err != nil {
		return fmt.Errorf("error creating Docker artifact removal command: %w", err)
	}
	err = c.printOutput(cmd)
	if err != nil {
		return fmt.Errorf("error removing Docker artifacts: %w", err)
	}

	// Delete the Hyperdrive directory
	path, err := homedir.Expand(configPath)
	if err != nil {
		return fmt.Errorf("error loading Hyperdrive directory: %w", err)
	}
	fmt.Printf("Deleting Hyperdrive directory (%s)...\n", path)
	cmd = fmt.Sprintf("%s rm -rf %s", rootCmd, path)
	_, err = c.readOutput(cmd)
	if err != nil {
		return fmt.Errorf("error deleting Hyperdrive directory: %w", err)
	}

	fmt.Println("Termination complete.")

	return nil
}

// Print the Hyperdrive service status
func (c *HyperdriveClient) PrintServiceStatus(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "ps")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Print the Hyperdrive service logs
func (c *HyperdriveClient) PrintServiceLogs(composeFiles []string, tail string, serviceNames ...string) error {
	sanitizedStrings := make([]string, len(serviceNames))
	for i, serviceName := range serviceNames {
		sanitizedStrings[i] = shellescape.Quote(serviceName)
	}
	cmd, err := c.compose(composeFiles, fmt.Sprintf("logs -f --tail %s %s", shellescape.Quote(tail), strings.Join(sanitizedStrings, " ")))
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Print the Hyperdrive service stats
func (c *HyperdriveClient) PrintServiceStats(composeFiles []string) error {
	// Get service container IDs
	cmd, err := c.compose(composeFiles, "ps -q")
	if err != nil {
		return err
	}
	containers, err := c.readOutput(cmd)
	if err != nil {
		return err
	}
	containerIds := strings.Split(strings.TrimSpace(string(containers)), "\n")

	// Print stats
	return c.printOutput(fmt.Sprintf("docker stats %s", strings.Join(containerIds, " ")))
}

// Print the Hyperdrive service compose config
func (c *HyperdriveClient) PrintServiceCompose(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "config")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Get the Hyperdrive service version
func (c *HyperdriveClient) GetServiceVersion() (string, error) {
	// Get service container version output
	response, err := c.Api.Service.Version()
	if err != nil {
		return "", fmt.Errorf("error requesting Hyperdrive service version: %w", err)
	}
	versionString := response.Data.Version

	// Make sure it's a semantic version
	version, err := semver.Make(versionString)
	if err != nil {
		return "", fmt.Errorf("error parsing Hyperdrive service version number from output '%s': %w", versionString, err)
	}

	// Return the parsed semantic version (extra safety)
	return version.String(), nil
}

// Deletes the data directory, including the node wallet and all validator keys, and restarts the Docker containers
func (c *HyperdriveClient) PurgeData(composeFiles []string) error {
	// Get the command to run with root privileges
	rootCmd, err := c.getEscalationCommand()
	if err != nil {
		return fmt.Errorf("could not get privilege escalation command: %w", err)
	}

	// Get the config
	cfg, _, err := c.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading user settings: %w", err)
	}

	// Shut down the containers
	fmt.Println("Stopping containers...")
	err = c.PauseService(composeFiles)
	if err != nil {
		return fmt.Errorf("error stopping Docker containers: %w", err)
	}

	// Delete the user's data directory
	dataPath, err := homedir.Expand(cfg.Hyperdrive.UserDataPath.Value)
	if err != nil {
		return fmt.Errorf("error loading data path: %w", err)
	}
	fmt.Println("Deleting data...")
	cmd := fmt.Sprintf("%s rm -f %s", rootCmd, dataPath)
	_, err = c.readOutput(cmd)
	if err != nil {
		return fmt.Errorf("error deleting data: %w", err)
	}

	// Start the containers
	fmt.Println("Starting containers...")
	err = c.StartService(composeFiles)
	if err != nil {
		return fmt.Errorf("error starting Docker containers: %w", err)
	}

	fmt.Println("Purge complete.")
	return nil
}

// Runs the prune provisioner
func (c *HyperdriveClient) RunPruneProvisioner(container string, volume string, image string) error {

	// Run the prune provisioner
	cmd := fmt.Sprintf("docker run --rm --name %s -v %s:/ethclient %s", container, volume, image)
	output, err := c.readOutput(cmd)
	if err != nil {
		return err
	}

	outputString := strings.TrimSpace(string(output))
	if outputString != "" {
		return fmt.Errorf("unexpected output running the prune provisioner: %s", outputString)
	}

	return nil

}

// Runs the prune provisioner
func (c *HyperdriveClient) RunNethermindPruneStarter(container string) error {
	cmd := fmt.Sprintf("docker exec %s %s %s", container, nethermindPruneStarterCommand, nethermindAdminUrl)
	err := c.printOutput(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Runs the EC migrator
func (c *HyperdriveClient) RunEcMigrator(container string, volume string, targetDir string, mode string, image string) error {
	cmd := fmt.Sprintf("docker run --rm --name %s -v %s:/ethclient -v %s:/mnt/external -e EC_MIGRATE_MODE='%s' %s", container, volume, targetDir, mode, image)
	err := c.printOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Gets the size of the target directory via the EC migrator for importing, which should have the same permissions as exporting
func (c *HyperdriveClient) GetDirSizeViaEcMigrator(container string, targetDir string, image string) (uint64, error) {
	cmd := fmt.Sprintf("docker run --rm --name %s -v %s:/mnt/external -e OPERATION='size' %s", container, targetDir, image)
	output, err := c.readOutput(cmd)
	if err != nil {
		return 0, fmt.Errorf("error getting source directory size: %w", err)
	}

	trimmedOutput := strings.TrimRight(string(output), "\n")
	dirSize, err := strconv.ParseUint(trimmedOutput, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing directory size output [%s]: %w", trimmedOutput, err)
	}

	return dirSize, nil
}

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
	copyOverrideFiles(overrideSourceDir, overrideFolder)

	// Clear out the runtime folder and remake it
	runtimeFolder := filepath.Join(hyperdriveDir, runtimeDir)
	err := os.RemoveAll(runtimeFolder)
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
		TemplatePath: templatesDir,
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
			copyOverrideFiles(srcPath, targetPath)
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

// Handle composing for modules
func (c *HyperdriveClient) composeModule(global *GlobalConfig, module hdconfig.IModuleConfig, hyperdriveDir string, deployedContainers []string) ([]string, error) {
	moduleName := module.GetModuleName()
	composePaths := template.ComposePaths{
		RuntimePath:  filepath.Join(hyperdriveDir, runtimeDir, hdconfig.ModulesName, moduleName),
		TemplatePath: filepath.Join(templatesDir, hdconfig.ModulesName, moduleName),
		OverridePath: filepath.Join(hyperdriveDir, overrideDir, hdconfig.ModulesName, moduleName),
	}

	// These containers always run
	toDeploy := module.GetContainersToDeploy()

	// Make the modules folder
	err := os.MkdirAll(composePaths.RuntimePath, 0775)
	if err != nil {
		return []string{}, fmt.Errorf("error creating modules runtime folder (%s): %w", composePaths.RuntimePath, err)
	}

	for _, containerName := range toDeploy {
		containers, err := composePaths.File(string(containerName)).Write(global)
		if err != nil {
			return []string{}, fmt.Errorf("could not create %s container definition: %w", containerName, err)
		}
		deployedContainers = append(deployedContainers, containers...)
	}

	return deployedContainers, nil
}

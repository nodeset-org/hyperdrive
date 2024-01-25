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
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	debugColor                    color.Attribute = color.FgYellow
	nethermindPruneStarterCommand string          = "DELETE_ME"
)

// Install Hyperdrive
func (c *Client) InstallService(verbose, noDeps bool, version, path string) error {
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
	script, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if fmt.Sprint(len(script)) != resp.Header.Get("content-length") {
		return fmt.Errorf("downloaded script length %d did not match content-length header %s", len(script), resp.Header.Get("content-length"))
	}

	// Initialize installation command
	cmd := c.newCommand(fmt.Sprintf("sh -s -- %s", strings.Join(flags, " ")))

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
		return fmt.Errorf("Could not install Hyperdrive service: %s", errMessage)
	}
	return nil
}

// Install the update tracker
func (c *Client) InstallUpdateTracker(verbose bool, version string) error {
	// Get installation script flags
	flags := []string{
		"-v", fmt.Sprintf("%s", shellescape.Quote(version)),
	}

	// Download the installer package
	url := fmt.Sprintf(UpdateTrackerURL, version)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading installer package: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading installer package failed with code %d - [%s]", resp.StatusCode, resp.Status)
	}

	// Create a temp file for it
	path := filepath.Join(os.TempDir(), "install-update-tracker.sh")
	tempFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("error creating temporary file for installer script: %w", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing installer package to disk: %w", err)
	}
	tempFile.Close()

	// Initialize installation command
	cmd := c.newCommand(fmt.Sprintf("sh -c \"%s %s\"", path, strings.Join(flags, " ")))
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
		return fmt.Errorf("Could not install Hyperdrive update tracker: %s", errMessage)
	}
	return nil
}

// Start the Hyperdrive service
func (c *Client) StartService(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "up -d --remove-orphans --quiet-pull")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Pause the Hyperdrive service
func (c *Client) PauseService(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "stop")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Stop the Hyperdrive service
func (c *Client) StopService(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "down -v")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Stop Hyperdrive and remove the config folder
func (c *Client) TerminateService(composeFiles []string, configPath string) error {
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

	// Delete the RP directory
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
func (c *Client) PrintServiceStatus(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "ps")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Print the Hyperdrive service logs
func (c *Client) PrintServiceLogs(composeFiles []string, tail string, serviceNames ...string) error {
	sanitizedStrings := make([]string, len(serviceNames))
	for i, serviceName := range serviceNames {
		sanitizedStrings[i] = fmt.Sprintf("%s", shellescape.Quote(serviceName))
	}
	cmd, err := c.compose(composeFiles, fmt.Sprintf("logs -f --tail %s %s", shellescape.Quote(tail), strings.Join(sanitizedStrings, " ")))
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Print the Hyperdrive service stats
func (c *Client) PrintServiceStats(composeFiles []string) error {
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
func (c *Client) PrintServiceCompose(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "config")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Get the Hyperdrive service version
func (c *Client) GetServiceVersion() (string, error) {
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

// Runs the prune provisioner
func (c *Client) RunPruneProvisioner(container string, volume string, image string) error {

	// Run the prune provisioner
	cmd := fmt.Sprintf("docker run --rm --name %s -v %s:/ethclient %s", container, volume, image)
	output, err := c.readOutput(cmd)
	if err != nil {
		return err
	}

	outputString := strings.TrimSpace(string(output))
	if outputString != "" {
		return fmt.Errorf("Unexpected output running the prune provisioner: %s", outputString)
	}

	return nil

}

// Runs the prune provisioner
func (c *Client) RunNethermindPruneStarter(container string) error {
	cmd := fmt.Sprintf("docker exec %s %s %s", container, nethermindPruneStarterCommand, nethermindAdminUrl)
	err := c.printOutput(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Runs the EC migrator
func (c *Client) RunEcMigrator(container string, volume string, targetDir string, mode string, image string) error {
	cmd := fmt.Sprintf("docker run --rm --name %s -v %s:/ethclient -v %s:/mnt/external -e EC_MIGRATE_MODE='%s' %s", container, volume, targetDir, mode, image)
	err := c.printOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Gets the size of the target directory via the EC migrator for importing, which should have the same permissions as exporting
func (c *Client) GetDirSizeViaEcMigrator(container string, targetDir string, image string) (uint64, error) {
	cmd := fmt.Sprintf("docker run --rm --name %s -v %s:/mnt/external -e OPERATION='size' %s", container, targetDir, image)
	output, err := c.readOutput(cmd)
	if err != nil {
		return 0, fmt.Errorf("Error getting source directory size: %w", err)
	}

	trimmedOutput := strings.TrimRight(string(output), "\n")
	dirSize, err := strconv.ParseUint(trimmedOutput, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Error parsing directory size output [%s]: %w", trimmedOutput, err)
	}

	return dirSize, nil
}

// Build a docker compose command
func (c *Client) compose(composeFiles []string, args string) (string, error) {
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
		return "", fmt.Errorf("Settings file not found. Please run `hyperdrive service config` to set up Hyperdrive before starting it.")
	}

	// Check config
	if cfg.ClientMode.Value == types.ClientMode_Unknown {
		return "", fmt.Errorf("You haven't selected local or external mode for your clients yet.\nPlease run 'hyperdrive service config' before running this command.")
	} else if cfg.IsLocalMode() && cfg.LocalExecutionConfig.ExecutionClient.Value == types.ExecutionClient_Unknown {
		return "", errors.New("No Execution Client selected. Please run 'hyperdrive service config' before running this command.")
	}
	if cfg.IsLocalMode() && cfg.LocalBeaconConfig.BeaconNode.Value == types.BeaconNode_Unknown {
		return "", errors.New("No Beacon Node selected. Please run 'hyperdrive service config' before running this command.")
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
	return fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker compose --project-directory %s %s %s", cfg.ProjectName.Value, shellescape.Quote(expandedConfigPath), strings.Join(composeFileFlags, " "), args), nil
}

// Deploys all of the appropriate docker compose template files and provisions them based on the provided configuration
func (c *Client) deployTemplates(cfg *config.HyperdriveConfig, hyperdriveDir string) ([]string, error) {
	// Check for the folders
	runtimeFolder := filepath.Join(hyperdriveDir, runtimeDir)
	templatesFolder := filepath.Join(hyperdriveDir, templatesDir)
	_, err := os.Stat(templatesFolder)
	if os.IsNotExist(err) {
		return []string{}, fmt.Errorf("templates folder [%s] does not exist", templatesFolder)
	}
	overrideFolder := filepath.Join(hyperdriveDir, overrideDir)
	_, err = os.Stat(overrideFolder)
	if os.IsNotExist(err) {
		return []string{}, fmt.Errorf("override folder [%s] does not exist", overrideFolder)
	}

	// Clear out the runtime folder and remake it
	err = os.RemoveAll(runtimeFolder)
	if err != nil {
		return []string{}, fmt.Errorf("error deleting runtime folder [%s]: %w", runtimeFolder, err)
	}
	err = os.Mkdir(runtimeFolder, 0775)
	if err != nil {
		return []string{}, fmt.Errorf("error creating runtime folder [%s]: %w", runtimeFolder, err)
	}

	composePaths := template.ComposePaths{
		RuntimePath:  runtimeFolder,
		TemplatePath: templatesFolder,
		OverridePath: overrideFolder,
	}

	// Read and substitute the templates
	deployedContainers := []string{}

	// These containers always run
	toDeploy := []types.ContainerID{
		types.ContainerID_Daemon,
	}

	// Check if we are running the Execution Layer locally
	if cfg.IsLocalMode() {
		toDeploy = append(toDeploy, types.ContainerID_ExecutionClient)
		toDeploy = append(toDeploy, types.ContainerID_BeaconNode)
	}

	// Check the metrics containers
	if cfg.Metrics.EnableMetrics.Value == true {
		toDeploy = append(toDeploy,
			types.ContainerID_Grafana,
			types.ContainerID_Exporter,
			types.ContainerID_Prometheus,
		)
	}

	for _, containerName := range toDeploy {
		containers, err := composePaths.File(string(containerName)).Write(cfg)
		if err != nil {
			return []string{}, fmt.Errorf("could not create %s container definition: %w", containerName, err)
		}
		deployedContainers = append(deployedContainers, containers...)
	}

	return deployedContainers, nil
}

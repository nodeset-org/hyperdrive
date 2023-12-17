package client

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/a8m/envsubst"
	"github.com/alessio/shellescape"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client/routes"
	"github.com/nodeset-org/hyperdrive/shared/config"
	cfgtypes "github.com/nodeset-org/hyperdrive/shared/types/config"
)

const (
	runtimeDir        string = "runtime"
	templatesDir      string = "templates"
	overrideDir       string = "override"
	templateSuffix    string = ".tmpl"
	composeFileSuffix string = ".yml"
)

// Rocket Pool client
type Client struct {
	HyperdriveDir string
	Cfg           *config.HyperdriveConfig
	Api           *routes.ApiRequester
}

// Create a new Hyperdrive client
func NewClient(installDir string) (*Client, error) {
	mgr := config.NewConfigManager(installDir)
	cfg, isNew, err := mgr.LoadOrCreateConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting Hyperdrive config: %w", err)
	}
	if isNew {
		return nil, fmt.Errorf("Settings file not found. Please run `hyperdrive service config` to set up Hyperdrive before starting it.")
	}

	socketPath := cfg.DaemonSocketPath.Value
	client := &Client{
		HyperdriveDir: installDir,
		Cfg:           cfg,
		Api:           routes.NewApiRequester(socketPath),
	}

	return client, nil
}

// Start the Rocket Pool service
func (c *Client) StartService(composeFiles []string) error {
	// Start all of the containers
	cmdArgs, err := c.compose(composeFiles, "up", "-d", "--remove-orphans", "--quiet-pull")
	if err != nil {
		return err
	}
	cmd := exec.Command("docker", cmdArgs...)
	return c.printOutput(cmd)
}

// Pause the Rocket Pool service
func (c *Client) PauseService(composeFiles []string) error {
	cmdArgs, err := c.compose(composeFiles, "stop")
	if err != nil {
		return err
	}
	cmd := exec.Command("docker", cmdArgs...)
	return c.printOutput(cmd)
}

// Stop the Rocket Pool service
func (c *Client) StopService(composeFiles []string) error {
	cmdArgs, err := c.compose(composeFiles, "down", "-v")
	if err != nil {
		return err
	}
	cmd := exec.Command("docker", cmdArgs...)
	return c.printOutput(cmd)
}

// Build a docker compose command
func (c *Client) compose(composeFiles []string, args ...string) ([]string, error) {
	cmdArgs := []string{
		"compose",
		"--project-directory",
		shellescape.Quote(c.HyperdriveDir),
	}

	// Set up environment variables and deploy the template config files
	settings := c.Cfg.GenerateEnvironmentVariables()

	// Deploy the templates and run environment variable substitution on them
	deployedContainers, err := c.deployTemplates(settings)
	if err != nil {
		return nil, fmt.Errorf("error deploying Docker templates: %w", err)
	}

	// Include all of the relevant docker compose definition files
	for _, container := range deployedContainers {
		cmdArgs = append(cmdArgs, "-f", shellescape.Quote(container))
	}
	for _, container := range composeFiles {
		cmdArgs = append(cmdArgs, "-f", shellescape.Quote(container))
	}

	cmdArgs = append(cmdArgs, args...)
	return cmdArgs, nil
}

// Deploys all of the appropriate docker compose template files and provisions them based on the provided configuration
func (c *Client) deployTemplates(settings map[string]string) ([]string, error) {

	// Check for the folders
	runtimeFolder := filepath.Join(c.HyperdriveDir, runtimeDir)
	templatesFolder := filepath.Join(c.HyperdriveDir, templatesDir)
	_, err := os.Stat(templatesFolder)
	if os.IsNotExist(err) {
		return []string{}, fmt.Errorf("templates folder [%s] does not exist", templatesFolder)
	}
	overrideFolder := filepath.Join(c.HyperdriveDir, overrideDir)
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

	// Set the environment variables for substitution
	oldValues := map[string]string{}
	for varName, varValue := range settings {
		oldValues[varName] = os.Getenv(varName)
		os.Setenv(varName, varValue)
	}
	defer func() {
		// Unset the env vars
		for name, value := range oldValues {
			os.Setenv(name, value)
		}
	}()

	// Read and substitute the templates
	deployedContainers := []string{}

	// Node
	contents, err := envsubst.ReadFile(filepath.Join(templatesFolder, string(cfgtypes.ContainerID_Daemon)+templateSuffix))
	if err != nil {
		return []string{}, fmt.Errorf("error reading and substituting node container template: %w", err)
	}
	nodeComposePath := filepath.Join(runtimeFolder, string(cfgtypes.ContainerID_Daemon)+composeFileSuffix)
	err = os.WriteFile(nodeComposePath, contents, 0664)
	if err != nil {
		return []string{}, fmt.Errorf("could not write node container file to %s: %w", nodeComposePath, err)
	}
	deployedContainers = append(deployedContainers, nodeComposePath)
	deployedContainers = append(deployedContainers, filepath.Join(overrideFolder, string(cfgtypes.ContainerID_Daemon)+composeFileSuffix))

	return deployedContainers, nil
}

func (c *Client) printOutput(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Wait for the command to exit
	return cmd.Wait()

	/*
		output, err := cmd.Output()
		if err != nil {
			exitErr, isExitErr := err.(*exec.ExitError)
			if isExitErr {
				return fmt.Errorf("exit code %d, message: %s", exitErr.ExitCode(), string(exitErr.Stderr))
			}
			return err
		}
		fmt.Println(output)
		return nil
	*/
}

package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/blang/semver/v4"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

const (
	debugColor                    color.Attribute = color.FgYellow
	nethermindPruneStarterCommand string          = "DELETE_ME"
	nethermindAdminUrl            string          = "http://127.0.0.1:7434"

	templatesDir       string = "/usr/share/hyperdrive/templates"
	overrideSourceDir  string = "/usr/share/hyperdrive/override"
	overrideDir        string = "override"
	runtimeDir         string = "runtime"
	metricsDir         string = "metrics"
	extraScrapeJobsDir string = "extra-scrape-jobs"
	modulePrometheusSd string = "prometheus-sd"
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

// Pause the Hyperdrive service, shutting it down without removing the Docker artifacts
func (c *HyperdriveClient) StopService(composeFiles []string) error {
	cmd, err := c.compose(composeFiles, "stop")
	if err != nil {
		return err
	}
	return c.printOutput(cmd)
}

// Stop the Hyperdrive service, shutting it down and removing the Docker artifacts
func (c *HyperdriveClient) DownService(composeFiles []string, includeVolumes bool) error {
	args := "down"
	if includeVolumes {
		args += " -v"
	}
	cmd, err := c.compose(composeFiles, args)
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

// Print the Hyperdrive daemon logs
func (c *HyperdriveClient) PrintDaemonLogs(composeFiles []string, tail string, logPaths ...string) error {
	cmd := fmt.Sprintf("tail -f %s %s", tail, strings.Join(logPaths, " "))
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

// Deletes the data directory, including the node wallet and all validator keys, and restarts the Docker containers if requested
func (c *HyperdriveClient) PurgeData(composeFiles []string, restart bool) error {
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
	err = c.StopService(composeFiles)
	if err != nil {
		return fmt.Errorf("error stopping Docker containers: %w", err)
	}

	// Delete the user's data directory
	dataPath, err := homedir.Expand(cfg.Hyperdrive.UserDataPath.Value)
	if err != nil {
		return fmt.Errorf("error loading data path: %w", err)
	}
	fmt.Println("Deleting data...")
	cmd := fmt.Sprintf("%s rm -rf %s", rootCmd, dataPath)
	_, err = c.readOutput(cmd)
	if err != nil {
		return fmt.Errorf("error deleting data: %w", err)
	}

	if restart {
		// Start the containers
		fmt.Println("Starting containers...")
		err = c.StartService(composeFiles)
		if err != nil {
			return fmt.Errorf("error starting Docker containers: %w", err)
		}
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

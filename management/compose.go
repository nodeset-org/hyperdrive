package management

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
)

// OutputStreams is an interface for getting standard output and error streams to attach to exec processes.
type OutputStreams interface {
	// The standard output stream to use for the process
	Stdout() io.Writer

	// The standard error stream to use for the process
	Stderr() io.Writer
}

// Docker Compose's format for project file info
type projectFileDetails struct {
	Name        string `json:"Name"`
	Status      string `json:"Status"`
	ConfigFiles string `json:"ConfigFiles"` // Comma-separated list of full file paths
}

// Get the list of full file paths for each Docker Compose file in a project.
// If the project doesn't exist, an empty list is returned.
func GetFilesForProject(project string) ([]string, error) {
	cmd := exec.Command("docker", "compose", "ls", "--format", "json", "--all")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error listing Docker Compose projects: %w", err)
	}
	var response []projectFileDetails
	err = json.Unmarshal(out, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling Docker Compose project list: %w", err)
	}
	for _, projectDetails := range response {
		if projectDetails.Name != project {
			continue
		}
		return strings.Split(projectDetails.ConfigFiles, ","), nil
	}
	return []string{}, nil
}

// Start a Docker Compose project
func StartProject(project string, files []string) error {
	args := []string{
		"compose",
		"-p",
		project,
	}
	for _, file := range files {
		args = append(args, "-f", file)
	}
	args = append(args, "up", "--detach", "--remove-orphans", "--quiet-pull")
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error starting project [%s]: %w", project, err)
	}
	return nil
}

// Stop a Docker Compose project.
// If services is empty, all services are stopped.
// Otherwise, only the specified services within the project are stopped.
func StopProject(project string, services []string) error {
	args := []string{
		"compose",
		"-p",
		project,
	}

	// Get the files corresponding to each service
	if len(services) > 0 {
		files, err := GetFilesForProject(project)
		if err != nil {
			return fmt.Errorf("error stopping project [%s]: %w", project, err)
		}
		foundServices := make(map[string]bool)
		for _, service := range services {
			foundServices[service] = false
			for _, file := range files {
				config, err := ParseComposeFile(project, file)
				if err != nil {
					return fmt.Errorf("error stopping project [%s]: %w", project, err)
				}
				if config.Name == service {
					foundServices[service] = true
					args = append(args, "-f", file)
					break
				}
			}
		}
		for service, found := range foundServices {
			if !found {
				return fmt.Errorf("service [%s] not found in project [%s]", service, project)
			}
		}
	}

	args = append(args, "stop")
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error stopping project [%s]: %w", project, err)
	}
	return nil
}

// Delete a Docker Compose project
func DownProject(project string, includeVolumes bool) error {
	args := []string{
		"compose",
		"-p",
		project,
		"down",
	}
	if includeVolumes {
		args = append(args, "--volumes")
	}
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error deleting project [%s]: %w", project, err)
	}
	return nil
}

// Parse a Docker Compose file and return the service configuration
func ParseComposeFile(projectName string, file string) (*types.ServiceConfig, error) {
	// Load the Docker Compose file
	details, err := loader.LoadConfigFiles(context.Background(), []string{file}, filepath.Dir(file), func(o *loader.Options) {
		o.SetProjectName(projectName, true)
	})
	if err != nil {
		return nil, fmt.Errorf("error loading Docker Compose file [%s]: %w", file, err)
	}
	project, err := loader.LoadWithContext(context.Background(), *details, func(o *loader.Options) {
		o.SetProjectName(projectName, true)
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing Docker Compose file [%s]: %w", file, err)
	}
	if len(project.Services) > 1 {
		return nil, fmt.Errorf("multiple services found in Docker Compose file [%s]", file)
	}
	for _, service := range project.Services {
		// Return the first service found since it should be the only one
		return &service, nil
	}
	return nil, fmt.Errorf("no services found in Docker Compose file [%s]", file)
}

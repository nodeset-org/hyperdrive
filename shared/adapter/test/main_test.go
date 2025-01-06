package adapter_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/nodeset-org/hyperdrive/shared/adapter"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
)

const (
	exampleTag           string = "nodeset/hyperdrive-example-adapter:v0.1.0"
	exampleContainerName string = "hde-adapter_testing"
	logDir               string = "/tmp/hde-adapter-test/log"
	cfgDir               string = "/tmp/hde-adapter-test/cfg"
	keyPath              string = "/tmp/hde-adapter-test/key"
	testKey              string = "test-key"
)

var (
	// Adapter client
	ac *adapter.AdapterClient

	// Docker client
	docker *client.Client
)

func TestMain(m *testing.M) {
	// Create a Docker client
	var err error
	docker, err = client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		fail(fmt.Errorf("error creating Docker client: %w", err))
	}

	// Check if the adapter container is already created
	adapterID := getAdapterContainerID()
	if adapterID != "" {
		fail(fmt.Errorf("adapter container already exists - please remove it before running tests"))
	}

	// Initialize everything and get the adapter container info
	initializeArtifacts()
	adapterID = getAdapterContainerID()
	if adapterID == "" {
		fail(fmt.Errorf("adapter container not found"))
	}
	containerJson, err := docker.ContainerInspect(context.Background(), adapterID)
	if err != nil {
		fail(fmt.Errorf("error inspecting adapter container: %w", err))
	}

	// Create the adapter client
	ac, err = adapter.NewAdapterClient(exampleContainerName, containerJson.Config.Entrypoint, string(testKey))
	if err != nil {
		fail(fmt.Errorf("error creating adapter client: %w", err))
	}

	// Run the tests and clean up after
	code := m.Run()
	cleanup()
	os.Exit(code)
}

// Get the ID of the adapter container
func getAdapterContainerID() string {
	containerList, err := docker.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		fail(fmt.Errorf("error inspecting adapter container: %w", err))
	}
	for _, container := range containerList {
		for _, name := range container.Names {
			if name == "/"+exampleContainerName {
				return container.ID
			}
		}
	}
	return ""
}

// Create the Docker container and initialize the adapter client
func initializeArtifacts() {
	// Make the dirs
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fail(fmt.Errorf("error creating log dir: %w", err))
	}
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		fail(fmt.Errorf("error creating config dir: %w", err))
	}

	// Create the key file, or get the key if it already exists
	err := os.WriteFile(keyPath, []byte(testKey), 0644)
	if err != nil {
		fail(fmt.Errorf("error creating key file: %w", err))
	}

	// Create the container via the Docker CLI since it does stuff like pulling / tagging the image
	runCmd := fmt.Sprintf("docker run --rm -d -v %s:/hd/logs -v %s:/hd/config -v %s:/hd/secret --name %s %s", logDir, cfgDir, keyPath, exampleContainerName, exampleTag)
	_, err = command.ReadOutput(runCmd)
	if err != nil {
		fail(fmt.Errorf("error running adapter container: %w", err))
	}
}

// Clean up the test environment
func cleanup() {
	// Stop the adapter container
	if docker != nil {
		timeout := 0
		_ = docker.ContainerStop(context.Background(), exampleContainerName, container.StopOptions{Timeout: &timeout})
		_ = docker.ContainerRemove(context.Background(), exampleContainerName, container.RemoveOptions{Force: true})
	}

	// Remove the temp key file
	_ = os.Remove(keyPath)
}

// Clean up and exit with a failure
func fail(err error) {
	cleanup()
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(1)
}

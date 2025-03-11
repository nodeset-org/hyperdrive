package adapter_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	internal_test "github.com/nodeset-org/hyperdrive/internal/test"
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared/adapter"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
)

var (
	// Adapter client for global mode
	gac *adapter.AdapterClient

	// Adapter client for project mode
	pac *adapter.AdapterClient

	// Docker client
	docker *client.Client

	// Info for the example module
	modInfo *config.ModuleInfo
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

	// Check if the adapter containers are already created
	globalAdapterID := getContainerID(internal_test.GlobalAdapterContainerName)
	if globalAdapterID != "" {
		fail(fmt.Errorf("global adapter container already exists - please remove it before running tests"))
	}
	projectAdapterID := getContainerID(internal_test.ProjectAdapterContainerName)
	if projectAdapterID != "" {
		fail(fmt.Errorf("project adapter container already exists - please remove it before running tests"))
	}

	// Initialize everything and get the adapter container info
	initializeArtifacts()
	globalAdapterID = getContainerID(internal_test.GlobalAdapterContainerName)
	if globalAdapterID == "" {
		fail(fmt.Errorf("global adapter container not found"))
	}
	projectAdapterID = getContainerID(internal_test.ProjectAdapterContainerName)
	if projectAdapterID == "" {
		fail(fmt.Errorf("project adapter container not found"))
	}

	// Create the adapter clients
	gac, err = adapter.NewAdapterClient(internal_test.GlobalAdapterContainerName, string(internal_test.TestKey))
	if err != nil {
		fail(fmt.Errorf("error creating global adapter client: %w", err))
	}
	pac, err = adapter.NewAdapterClient(internal_test.ProjectAdapterContainerName, string(internal_test.TestKey))
	if err != nil {
		fail(fmt.Errorf("error creating project adapter client: %w", err))
	}

	// Get the module config info
	modCfgMeta, err := gac.GetConfigMetadata(context.Background())
	if err != nil {
		fail(fmt.Errorf("error getting module config metadata: %w", err))
	}
	modInfo = &config.ModuleInfo{
		Descriptor:    &internal_test.ExampleDescriptor,
		Configuration: modCfgMeta,
	}

	// Run the tests and clean up after
	code := m.Run()
	cleanup()
	os.Exit(code)
}

// Get the ID of the container with the given name
func getContainerID(name string) string {
	containerList, err := docker.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		fail(fmt.Errorf("error inspecting adapter container: %w", err))
	}
	for _, container := range containerList {
		for _, containerName := range container.Names {
			if containerName == "/"+name {
				return container.ID
			}
		}
	}
	return ""
}

// Create the Docker container and initialize the adapter client
func initializeArtifacts() {
	// Make the dirs
	if err := os.MkdirAll(internal_test.LogDir, 0755); err != nil {
		fail(fmt.Errorf("error creating log dir: %w", err))
	}
	if err := os.MkdirAll(internal_test.CfgDir, 0755); err != nil {
		fail(fmt.Errorf("error creating config dir: %w", err))
	}
	if err := os.MkdirAll(internal_test.RuntimeDir, 0755); err != nil {
		fail(fmt.Errorf("error creating runtime dir: %w", err))
	}
	if err := os.MkdirAll(internal_test.UserDir, 0755); err != nil {
		fail(fmt.Errorf("error creating user dir: %w", err))
	}
	if err := os.MkdirAll(internal_test.DataDir, 0755); err != nil {
		fail(fmt.Errorf("error creating data dir: %w", err))
	}
	if err := os.MkdirAll(filepath.Dir(internal_test.AdapterKeyPath), 0755); err != nil {
		fail(fmt.Errorf("error creating secrets dir: %w", err))
	}

	// Create the key file, or get the key if it already exists
	err := os.WriteFile(internal_test.AdapterKeyPath, []byte(internal_test.TestKey), 0644)
	if err != nil {
		fail(fmt.Errorf("error creating key file: %w", err))
	}

	// Create the docker network
	composeProjectName := utils.GetModuleComposeProjectName(internal_test.ProjectName, &internal_test.ExampleDescriptor)
	netCreateResponse, err := docker.NetworkCreate(context.Background(), composeProjectName+"_net", network.CreateOptions{
		Driver: "bridge",
		Scope:  "local",
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{}, // Use an autogenerated IPAM config
		},
		Internal:   true,
		Attachable: false,
		Ingress:    false,
		ConfigOnly: false,
	})
	if err != nil {
		fail(fmt.Errorf("error creating network: %w", err))
	}
	if netCreateResponse.Warning != "" {
		fmt.Printf("warning creating network: %s\n", netCreateResponse.Warning)
	}

	// Create the adapters via the Docker CLI since it does stuff like pulling / tagging the image
	runCmd := fmt.Sprintf(
		"docker run --rm -d -e HD_ADAPTER_MODE=global --name %s %s",
		internal_test.GlobalAdapterContainerName,
		internal_test.AdapterTag,
	)
	_, err = command.ReadOutput(runCmd)
	if err != nil {
		fail(fmt.Errorf("error running global adapter container: %w", err))
	}
	runCmd = fmt.Sprintf(
		"docker run --rm -d -e HD_ADAPTER_MODE=project -e HD_PROJECT_NAME=%s -e HD_CONFIG_DIR=%s -e HD_LOG_DIR=%s -e HD_KEY_FILE=%s -e HD_COMPOSE_DIR=%s -e HD_COMPOSE_PROJECT=%s -v %s:%s -v %s:%s -v %s:%s -v %s:%s -v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker:ro -v /usr/libexec/docker:/usr/libexec/docker:ro --network %s_net --name %s %s",
		internal_test.ProjectName,
		internal_test.CfgDir,
		internal_test.LogDir,
		internal_test.AdapterKeyPath,
		internal_test.RuntimeDir,
		composeProjectName,
		internal_test.CfgDir,
		internal_test.CfgDir,
		internal_test.LogDir,
		internal_test.LogDir,
		internal_test.RuntimeDir,
		internal_test.RuntimeDir,
		internal_test.AdapterKeyPath,
		internal_test.AdapterKeyPath,
		composeProjectName,
		internal_test.ProjectAdapterContainerName,
		internal_test.AdapterTag,
	)
	_, err = command.ReadOutput(runCmd)
	if err != nil {
		fail(fmt.Errorf("error running project adapter container: %w", err))
	}
}

// Delete all of the config files from disk
func deleteConfigs() error {
	info, err := os.ReadDir(internal_test.CfgDir)
	if err != nil {
		return fmt.Errorf("error enumerating config directory: %w", err)
	}
	for _, entry := range info {
		err = os.Remove(filepath.Join(internal_test.CfgDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("error removing config file [%s]: %w", entry.Name(), err)
		}
	}
	return nil
}

// Clean up the test environment
func cleanup() {
	// Stop the adapter container
	if docker != nil {
		composeProjectName := utils.GetModuleComposeProjectName(internal_test.ProjectName, &internal_test.ExampleDescriptor)
		timeout := 0
		_ = docker.ContainerStop(context.Background(), internal_test.GlobalAdapterContainerName, container.StopOptions{Timeout: &timeout})
		_ = docker.ContainerStop(context.Background(), internal_test.ProjectAdapterContainerName, container.StopOptions{Timeout: &timeout})
		_ = docker.ContainerRemove(context.Background(), internal_test.GlobalAdapterContainerName, container.RemoveOptions{Force: true})
		_ = docker.ContainerRemove(context.Background(), internal_test.ProjectAdapterContainerName, container.RemoveOptions{Force: true})
		_ = docker.NetworkRemove(context.Background(), composeProjectName+"_net")
	}

	// Remove the temp key file
	_ = os.Remove(internal_test.AdapterKeyPath)
}

// Clean up and exit with a failure
func fail(err error) {
	cleanup()
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(1)
}

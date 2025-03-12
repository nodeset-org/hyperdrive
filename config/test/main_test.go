package config_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/config"
	internal_test "github.com/nodeset-org/hyperdrive/internal/test"
	"github.com/nodeset-org/hyperdrive/management"
	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/adapter"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
)

var (
	// Adapter client
	ac *adapter.AdapterClient

	// Docker client
	docker *client.Client

	cfgMgr      *management.ConfigurationManager
	cfgInstance *config.HyperdriveSettings
	modMgr      *management.ModuleManager
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

	// Create the adapter client
	ac, err = adapter.NewAdapterClient(internal_test.GlobalAdapterContainerName, string(internal_test.TestKey))
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
			if name == "/"+internal_test.GlobalAdapterContainerName {
				return container.ID
			}
		}
	}
	return ""
}

// Create the Docker container and initialize the adapter client
func initializeArtifacts() {
	// Serialize the descriptor
	bytes, err := json.Marshal(internal_test.ExampleDescriptor)
	if err != nil {
		fail(fmt.Errorf("error serializing descriptor: %w", err))
	}

	// Make the dirs
	modulePath := filepath.Join(internal_test.SystemDir, shared.ModulesDir, string(internal_test.ExampleDescriptor.Name))
	descriptorPath := filepath.Join(modulePath, modules.DescriptorFilename)
	if err := os.MkdirAll(internal_test.LogDir, 0755); err != nil {
		fail(fmt.Errorf("error creating log dir: %w", err))
	}
	if err := os.MkdirAll(internal_test.CfgDir, 0755); err != nil {
		fail(fmt.Errorf("error creating config dir: %w", err))
	}
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		fail(fmt.Errorf("error creating module dir: %w", err))
	}
	if err := os.MkdirAll(filepath.Dir(internal_test.AdapterKeyPath), 0755); err != nil {
		fail(fmt.Errorf("error creating secrets dir: %w", err))
	}
	if err := os.MkdirAll(internal_test.UserDataPath, 0755); err != nil {
		fail(fmt.Errorf("error creating data dir: %w", err))
	}
	if err := os.WriteFile(descriptorPath, bytes, 0644); err != nil {
		fail(fmt.Errorf("error writing descriptor file: %w", err))
	}

	// Create the descriptor file
	err = os.WriteFile(descriptorPath, bytes, 0644)

	// Create the container via the Docker CLI since it does stuff like pulling / tagging the image
	runCmd := fmt.Sprintf(
		"docker run --rm -d -e HD_PROJECT_NAME=%s -v %s:/hd/logs -v %s:/hd/config -v %s:/hd/secret --name %s %s", internal_test.ProjectName, internal_test.LogDir, internal_test.CfgDir, internal_test.AdapterKeyPath, internal_test.GlobalAdapterContainerName, internal_test.AdapterTag)
	_, err = command.ReadOutput(runCmd)
	if err != nil {
		fail(fmt.Errorf("error running adapter container: %w", err))
	}

	// Set up the test config
	cfgMgr = management.NewConfigurationManager(internal_test.UserDir, internal_test.SystemDir)
	cfgInstance, err := config.CreateDefaultHyperdriveSettingsFromConfiguration(cfgMgr.HyperdriveConfiguration)
	if err != nil {
		fail(fmt.Errorf("error converting instance to known config type: %w", err))
	}
	cfgInstance.ProjectName = internal_test.ProjectName
	cfgInstance.UserDataPath = internal_test.UserDataPath

	// Set up the mod manager
	modMgr, err = management.NewModuleManager(shared.GetModulesDirectoryPath(internal_test.SystemDir), "", internal_test.CfgDir)
	if err != nil {
		fail(fmt.Errorf("error creating module manager: %w", err))
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
		timeout := 0
		_ = docker.ContainerStop(context.Background(), internal_test.GlobalAdapterContainerName, container.StopOptions{Timeout: &timeout})
		_ = docker.ContainerRemove(context.Background(), internal_test.GlobalAdapterContainerName, container.RemoveOptions{Force: true})
	}

	// Remove the temp files
	_ = os.Remove(internal_test.AdapterKeyPath)
	_ = os.RemoveAll(internal_test.UserDataPath)
	_ = os.RemoveAll(internal_test.SystemDir)
	_ = os.RemoveAll(internal_test.LogDir)
	_ = os.RemoveAll(internal_test.CfgDir)
	_ = os.RemoveAll(internal_test.UserDir)
}

// Clean up and exit with a failure
func fail(err error) {
	cleanup()
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(1)
}

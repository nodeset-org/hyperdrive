package adapter_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/nodeset-org/hyperdrive/modules/config"
	adapter "github.com/nodeset-org/hyperdrive/shared/adapter/test"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	version, err := ac.GetVersion(context.Background())
	if err != nil {
		t.Errorf("error getting version: %v", err)
	}
	t.Logf("Adapter version: %s", version)
	require.Equal(t, "0.1.0", version)
}

func TestGetLogFile(t *testing.T) {
	// Get the adapter log
	response, err := ac.GetLogFile(context.Background(), "adapter")
	require.NoError(t, err)
	require.Equal(t, "adapter.log", response.Path)
	t.Logf("Adapter log file path: %s", response.Path)

	// Get the service log
	response, err = ac.GetLogFile(context.Background(), "example")
	require.NoError(t, err)
	require.Equal(t, "service.log", response.Path)
	t.Logf("Service log file path: %s", response.Path)
}

func TestGetConfigMetadata(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	cfg, err := ac.GetConfigMetadata(context.Background())
	require.NoError(t, err)

	// Do a spot check of the config - full test is done elsewhere
	paramMap := make(map[string]config.IParameterMetadata)
	for _, param := range cfg.GetParameters() {
		paramMap[param.GetID().String()] = param
	}
	require.Len(t, paramMap, 6)
	exampleFloatName := "exampleFloat"
	require.Contains(t, paramMap, exampleFloatName)
	exampleFloat := paramMap[exampleFloatName]
	require.Equal(t, exampleFloat.GetName(), "Example Float")
	require.NotEmpty(t, exampleFloat.GetDescription())
	require.Equal(t, 50.0, exampleFloat.GetValueAsAny().(float64))
	require.Equal(t, 50.0, exampleFloat.GetDefaultAsAny().(float64))
	castedFloat := exampleFloat.(*config.FloatParameterMetadata)
	require.Equal(t, 0.0, castedFloat.MinValue)
	require.Equal(t, 100.0, castedFloat.MaxValue)

	// Make a map of sections
	sections := cfg.GetSections()
	sectionMap := make(map[string]config.ISectionMetadata)
	for _, section := range cfg.GetSections() {
		sectionMap[section.GetID().String()] = section
	}
	require.Len(t, sections, 2)
	require.Contains(t, sectionMap, "subConfig")
	require.Contains(t, sectionMap, "server")
	t.Logf("Config metadata: %v", cfg)
}

func TestProcessConfig(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	cfg, err := ac.GetConfigMetadata(context.Background())
	require.NoError(t, err)
	updateConfigSettings(t, cfg)

	// Process the config
	response, err := ac.ProcessConfig(context.Background(), cfg)
	require.NoError(t, err)
	require.Empty(t, response.Errors)
	require.Len(t, response.Ports, 1)
	require.Equal(t, uint16(8085), response.Ports["server/portMode"])
	t.Log("Config processed successfully")
}

func TestSetConfig(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	cfg, err := ac.GetConfigMetadata(context.Background())
	require.NoError(t, err)
	updateConfigSettings(t, cfg)

	// Set the config
	err = ac.SetConfig(context.Background(), cfg)
	require.NoError(t, err)

	// Get the config and check the values
	cfg, err = ac.GetConfigMetadata(context.Background())
	checkConfigSettings(t, cfg)
	t.Log("Config set successfully")
}

func TestGetConfigInstance(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	cfg, err := ac.GetConfigMetadata(context.Background())
	require.NoError(t, err)
	updateConfigSettings(t, cfg)

	// Set the config
	err = ac.SetConfig(context.Background(), cfg)
	require.NoError(t, err)

	// Get the config and check the values
	inst, err := ac.GetConfigInstance(context.Background())
	require.Equal(t, true, inst["exampleBool"])
	require.Equal(t, "three", inst["exampleChoice"])
	require.Equal(t, float64(8085), inst["server"].(map[string]any)["port"])
	t.Log("Config instance was correct")
}

func TestGetContaners(t *testing.T) {
	containers, err := ac.GetContainers(context.Background())
	if err != nil {
		t.Errorf("error getting containers: %v", err)
	}
	require.Len(t, containers, 1)
	require.Equal(t, "example", containers[0])
	t.Logf("Containers: %v", containers)
}

func TestRunCommand(t *testing.T) {
	defer func() {
		// Stop the service container
		if docker != nil {
			timeout := 0
			_ = docker.ContainerStop(context.Background(), serviceContainerName, container.StopOptions{Timeout: &timeout})
			_ = docker.ContainerRemove(context.Background(), serviceContainerName, container.RemoveOptions{Force: true})
		}
	}()

	// Make a logger
	logger := adapter.CreateLogger(t)

	// Set the config
	err := deleteConfigs()
	require.NoError(t, err)
	cfg, err := ac.GetConfigMetadata(context.Background())
	require.NoError(t, err)
	updateConfigSettings(t, cfg)
	err = ac.SetConfig(context.Background(), cfg)
	require.NoError(t, err)

	// Make sure the service is running
	runCmd := fmt.Sprintf("docker run --rm -d -v %s:/hd/logs -v %s:/hd/config -v %s:/hd/secret --network %s_net --name %s %s -i 0.0.0.0 -p 8085", logDir, cfgDir, keyPath, projectName, serviceContainerName, serviceTag)
	serviceRunOut, err := command.ReadOutput(runCmd)
	require.NoError(t, err)
	t.Logf("Service container started: %s", serviceRunOut)

	// Run the get-param command
	cmd := "config get-param exampleFloat"
	stdout, stderr, err := ac.RunNoninteractive(context.Background(), logger, cmd)
	require.Empty(t, stderr)
	require.NoError(t, err)

	// Check the output
	out := strings.TrimSpace(stdout)
	paramVal, err := strconv.ParseFloat(out, 64)
	require.NoError(t, err)
	require.Equal(t, 75.0, paramVal)
	t.Logf("Command ran successfully and returned %s", out)
}

func updateConfigSettings(t *testing.T, cfg config.IConfigurationMetadata) {
	// Set some values
	paramMap := make(map[string]config.IParameterMetadata)
	for _, param := range cfg.GetParameters() {
		paramMap[param.GetID().String()] = param
	}
	err := paramMap["exampleBool"].SetValue(true)
	require.NoError(t, err)
	err = paramMap["exampleChoice"].SetValue("three")
	require.NoError(t, err)
	err = paramMap["exampleFloat"].SetValue(75.0)
	require.NoError(t, err)

	// Set a subconfig value
	sectionMap := make(map[string]config.ISectionMetadata)
	for _, section := range cfg.GetSections() {
		sectionMap[section.GetID().String()] = section
	}
	subParamMap := make(map[string]config.IParameterMetadata)
	for _, param := range sectionMap["server"].GetParameters() {
		subParamMap[param.GetID().String()] = param
	}
	err = subParamMap["port"].SetValue(8085)
	require.NoError(t, err)
	err = subParamMap["portMode"].SetValue("open")
	require.NoError(t, err)
}

func checkConfigSettings(t *testing.T, cfg config.IConfigurationMetadata) {
	// Check some values
	paramMap := make(map[string]config.IParameterMetadata)
	for _, param := range cfg.GetParameters() {
		paramMap[param.GetID().String()] = param
	}
	require.True(t, paramMap["exampleBool"].GetValueAsAny().(bool))
	require.Equal(t, "three", paramMap["exampleChoice"].GetValueAsAny().(string))

	// Check a subconfig value
	sectionMap := make(map[string]config.ISectionMetadata)
	for _, section := range cfg.GetSections() {
		sectionMap[section.GetID().String()] = section
	}
	subParamMap := make(map[string]config.IParameterMetadata)
	for _, param := range sectionMap["server"].GetParameters() {
		subParamMap[param.GetID().String()] = param
	}
	require.Equal(t, uint64(8085), subParamMap["port"].GetValueAsAny().(uint64))
	require.Equal(t, "open", subParamMap["portMode"].GetValueAsAny().(string))
}

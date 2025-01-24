package adapter_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	internal_test "github.com/nodeset-org/hyperdrive/internal/test"
	"github.com/nodeset-org/hyperdrive/modules/config"
	adapter "github.com/nodeset-org/hyperdrive/shared/adapter/test"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
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
	paramMap := make(map[string]config.IParameter)
	for _, param := range cfg.GetParameters() {
		paramMap[param.GetID().String()] = param
	}
	require.Len(t, paramMap, 6)
	exampleFloatName := "exampleFloat"
	require.Contains(t, paramMap, exampleFloatName)
	exampleFloat := paramMap[exampleFloatName]
	require.Equal(t, exampleFloat.GetName(), "Example Float")
	require.NotEmpty(t, exampleFloat.GetDescription())
	require.Equal(t, 50.0, exampleFloat.GetDefault().(float64))
	castedFloat := exampleFloat.(*config.FloatParameter)
	require.Equal(t, 0.0, castedFloat.MinValue)
	require.Equal(t, 100.0, castedFloat.MaxValue)

	// Make a map of sections
	sections := cfg.GetSections()
	sectionMap := make(map[string]config.ISection)
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
	hdInstance := createHyperdriveConfigInstance(t)
	modInstance := hdInstance.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	updateConfigSettings(t, modInstance.Settings.GetSettings())

	// Process the config
	response, err := ac.ProcessConfig(context.Background(), hdInstance.SerializeToMap())
	require.NoError(t, err)
	require.Empty(t, response.Errors)
	require.Len(t, response.Ports, 1)
	require.Equal(t, uint16(8085), response.Ports["server/port"])
	t.Log("Config processed successfully")
}

func TestSetConfig(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	hdInstance := createHyperdriveConfigInstance(t)
	modInstance := hdInstance.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	updateConfigSettings(t, modInstance.Settings.GetSettings())

	// Set the config
	err = ac.SetConfig(context.Background(), hdInstance.SerializeToMap())
	require.NoError(t, err)
	t.Log("Config set successfully")
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
			_ = docker.ContainerStop(context.Background(), internal_test.ServiceContainerName, container.StopOptions{Timeout: &timeout})
			_ = docker.ContainerRemove(context.Background(), internal_test.ServiceContainerName, container.RemoveOptions{Force: true})
		}
	}()

	// Make a logger
	logger := adapter.CreateLogger(t)

	// Set the config
	err := deleteConfigs()
	require.NoError(t, err)
	hdInstance := createHyperdriveConfigInstance(t)
	modInstance := hdInstance.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	updateConfigSettings(t, modInstance.Settings.GetSettings())
	err = ac.SetConfig(context.Background(), hdInstance.SerializeToMap())
	require.NoError(t, err)

	// Make sure the service is running
	runCmd := fmt.Sprintf("docker run --rm -d -v %s:/hd/logs -v %s:/hd/config -v %s:/hd/secret --network %s_net --name %s %s -i 0.0.0.0 -p 8085", internal_test.LogDir, internal_test.CfgDir, internal_test.KeyPath, internal_test.ProjectName, internal_test.ServiceContainerName, internal_test.ServiceTag)
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

// Create a full Hyperdrive config instance for the test
func createHyperdriveConfigInstance(t *testing.T) *hdconfig.HyperdriveConfigInstance {
	cfgMgr := hdconfig.NewConfigurationManager(internal_test.UserDir, internal_test.SystemDir)
	inst := config.CreateModuleConfigurationInstance(cfgMgr.HyperdriveConfiguration)
	cfgInstance := hdconfig.NewHyperdriveConfigInstance()
	err := inst.ConvertToKnownType(cfgInstance)
	if err != nil {
		fail(fmt.Errorf("error converting instance to known config type: %w", err))
	}
	cfgInstance.ProjectName = internal_test.ProjectName
	cfgInstance.UserDataPath = internal_test.UserDataPath

	modCfgMeta, err := ac.GetConfigMetadata(context.Background())
	require.NoError(t, err)
	modInstance := &config.ModuleInstance{
		Enabled: true,
	}
	modInstance.Settings.CreateSettingsFromMetadata(modCfgMeta)
	cfgInstance.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()] = modInstance
	return cfgInstance
}

func updateConfigSettings(t *testing.T, cfg *config.ModuleConfigurationInstance) {
	// Set some values
	param, err := cfg.GetParameter("exampleBool")
	require.NoError(t, err)
	err = param.SetValue(true)
	require.NoError(t, err)

	param, err = cfg.GetParameter("exampleChoice")
	require.NoError(t, err)
	err = param.SetValue("three")
	require.NoError(t, err)

	param, err = cfg.GetParameter("exampleFloat")
	require.NoError(t, err)
	err = param.SetValue(75.0)
	require.NoError(t, err)

	// Set a subconfig value
	serverCfg, err := cfg.GetSection("server")
	require.NoError(t, err)

	subPort, err := serverCfg.GetParameter("port")
	require.NoError(t, err)
	err = subPort.SetValue(8085)
	require.NoError(t, err)

	subPortMode, err := serverCfg.GetParameter("portMode")
	require.NoError(t, err)
	err = subPortMode.SetValue("open")
	require.NoError(t, err)
}

func checkConfigSettings(t *testing.T, cfg *config.ModuleConfigurationInstance) {
	// Check some values
	exampleBool, err := cfg.GetParameter("exampleBool")
	require.NoError(t, err)
	require.True(t, exampleBool.GetValue().(bool))

	exampleChoice, err := cfg.GetParameter("exampleChoice")
	require.NoError(t, err)
	require.Equal(t, "three", exampleChoice.GetValue().(string))

	// Check a subconfig value
	serverCfg, err := cfg.GetSection("server")
	require.NoError(t, err)

	subPort, err := serverCfg.GetParameter("port")
	require.NoError(t, err)
	require.Equal(t, uint64(8085), subPort.GetValue().(uint64))

	subPortMode, err := serverCfg.GetParameter("portMode")
	require.NoError(t, err)
	require.Equal(t, "open", subPortMode.GetValue().(string))
}

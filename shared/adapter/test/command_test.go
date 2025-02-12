package adapter_test

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/docker/docker/api/types/container"
	internal_test "github.com/nodeset-org/hyperdrive/internal/test"
	"github.com/nodeset-org/hyperdrive/modules/config"
	adapter "github.com/nodeset-org/hyperdrive/shared/adapter/test"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	version, err := gac.GetVersion(context.Background())
	if err != nil {
		t.Errorf("error getting version: %v", err)
	}
	t.Logf("Adapter version: %s", version)
	require.Equal(t, "0.2.0", version)
}

func TestGetConfigMetadata(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	cfg, err := gac.GetConfigMetadata(context.Background())
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

func TestUpgradeInstance(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	hdSettings := createHyperdriveConfigInstance()
	modInstance := hdSettings.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	modSettings := config.CreateModuleSettings(modInfo.Configuration)
	updateConfigSettings(t, modSettings)
	modInstance.SetSettings(modSettings)

	// Manually downgrade the old config to v0.1.0
	legacyModInstance := &config.ModuleInstance{
		Enabled: modInstance.Enabled,
		Version: "0.1.0",
	}
	legacyModInstance.SetSettings(modSettings)
	delete(legacyModInstance.Settings, "exampleUint")

	// Process the config
	require.Equal(t, "0.1.0", legacyModInstance.Version)
	require.NotContains(t, legacyModInstance.Settings, "exampleUint")
	upgradedInstance, err := gac.UpgradeInstance(context.Background(), legacyModInstance)
	require.NoError(t, err)
	require.Equal(t, modInstance.Version, upgradedInstance.Version)
	require.Contains(t, upgradedInstance.Settings, "exampleUint")
	require.Equal(t, upgradedInstance.Settings["exampleUint"], float64(42))
	t.Log("Config upgraded successfully")
}

func TestProcessSettings(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	oldHdSettings := createHyperdriveConfigInstance()
	oldModInstance := oldHdSettings.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	oldModSettings := config.CreateModuleSettings(modInfo.Configuration)
	oldModInstance.SetSettings(oldModSettings)

	hdSettings := createHyperdriveConfigInstance()
	modInstance := hdSettings.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	modSettings := config.CreateModuleSettings(modInfo.Configuration)
	updateConfigSettings(t, modSettings)
	modInstance.SetSettings(modSettings)

	// Process the config
	response, err := gac.ProcessSettings(context.Background(), oldHdSettings, hdSettings)
	require.NoError(t, err)
	require.Empty(t, response.Errors)
	require.Len(t, response.Ports, 1)
	require.Equal(t, uint16(8085), response.Ports["server/port"])
	require.Len(t, response.ServicesToRestart, 1)
	require.Equal(t, "example", response.ServicesToRestart[0])
	t.Log("Config processed successfully")
}

func TestSetSettings(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	hdSettings := createHyperdriveConfigInstance()
	modInstance := hdSettings.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	modSettings := config.CreateModuleSettings(modInfo.Configuration)
	updateConfigSettings(t, modSettings)
	modInstance.SetSettings(modSettings)

	// Set the config
	err = pac.SetSettings(context.Background(), hdSettings)
	require.NoError(t, err)
	t.Log("Config set successfully")
}

func TestStartStopRun(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)

	// Remove the service container
	if docker != nil {
		timeout := 0
		_ = docker.ContainerStop(context.Background(), internal_test.ServiceContainerName, container.StopOptions{Timeout: &timeout})
		_ = docker.ContainerRemove(context.Background(), internal_test.ServiceContainerName, container.RemoveOptions{Force: true})
	}

	// Set the config
	hdSettings := createHyperdriveConfigInstance()
	modInstance := hdSettings.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	modSettings := config.CreateModuleSettings(modInfo.Configuration)
	updateConfigSettings(t, modSettings)
	modInstance.SetSettings(modSettings)
	err = pac.SetSettings(context.Background(), hdSettings)
	require.NoError(t, err)

	// Make sure the service is not running
	found := false
	containers, err := docker.ContainerList(context.Background(), container.ListOptions{All: true})
	require.NoError(t, err)
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+internal_test.ServiceContainerName {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	require.False(t, found)

	// Create the service compose file
	serviceFilePath := filepath.Join(internal_test.RuntimeDir, "example.yml")
	tmpl := template.New("service")
	_, err = tmpl.Parse(internal_test.ServiceComposeTemplate)
	require.NoError(t, err)
	runtimeFile, err := os.OpenFile(serviceFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	require.NoError(t, err)
	defer runtimeFile.Close()
	tmplSrc := &internal_test.InternalTestTemplateSource{}
	err = tmpl.Execute(runtimeFile, tmplSrc)
	require.NoError(t, err)
	runtimeFile.Close()

	// Start the services
	composeProjectName := internal_test.ProjectName + "-" + string(internal_test.ExampleDescriptor.Shortcut)
	err = pac.Start(context.Background(), hdSettings, composeProjectName)
	require.NoError(t, err)
	t.Log("Services started successfully")

	// Make sure the service is running now
	found = false
	containers, err = docker.ContainerList(context.Background(), container.ListOptions{All: false})
	require.NoError(t, err)
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+internal_test.ServiceContainerName {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	require.True(t, found)

	// Make a logger
	logger := adapter.CreateLogger(t)

	// Run the get-param command
	cmd := "config get-param exampleFloat"
	stdout, stderr, err := pac.RunNoninteractive(context.Background(), logger, cmd)
	require.Empty(t, stderr)
	require.NoError(t, err)

	// Check the output
	out := strings.TrimSpace(stdout)
	paramVal, err := strconv.ParseFloat(out, 64)
	require.NoError(t, err)
	require.Equal(t, 75.0, paramVal)
	t.Logf("Command ran successfully and returned %s", out)

	// Stop the services
	err = pac.Stop(context.Background(), composeProjectName)
	require.NoError(t, err)
	t.Log("Services stopped successfully")

	// Make sure the service is not running
	found = false
	containers, err = docker.ContainerList(context.Background(), container.ListOptions{All: true})
	require.NoError(t, err)
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+internal_test.ServiceContainerName {
				found = true
				require.Equal(t, "exited", container.State)
				break
			}
		}
		if found {
			break
		}
	}
	require.True(t, found)

	// Clean up
	err = os.Remove(serviceFilePath)
	require.NoError(t, err)
}

// Create a full Hyperdrive config instance for the test
func createHyperdriveConfigInstance() *hdconfig.HyperdriveSettings {
	hdSettings := hdconfig.NewHyperdriveSettings()
	hdSettings.ProjectName = internal_test.ProjectName
	hdSettings.UserDataPath = internal_test.UserDataPath

	// Create a module instance manually
	modInstance := &config.ModuleInstance{
		Enabled:  true,
		Version:  internal_test.ExampleDescriptor.Version.String(),
		Settings: map[string]any{},
	}
	hdSettings.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()] = modInstance
	return hdSettings
}

func updateConfigSettings(t *testing.T, cfg *config.ModuleSettings) {
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

	param, err = cfg.GetParameter("exampleUint")
	require.NoError(t, err)
	err = param.SetValue(6)
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

func checkConfigSettings(t *testing.T, cfg *config.ModuleSettings) {
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

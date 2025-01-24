package config_test

import (
	"path/filepath"
	"testing"

	internal_test "github.com/nodeset-org/hyperdrive/internal/test"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/stretchr/testify/require"
)

var (
	isCfgLoaded bool = false
)

func TestLoadModuleConfigs(t *testing.T) {
	if isCfgLoaded {
		return
	}
	err := deleteConfigs()
	require.NoError(t, err)
	results, err := cfgMgr.LoadModuleInfo(cfgInstance.ProjectName)
	require.NoError(t, err)

	require.Len(t, results, 1)
	require.NoError(t, results[0].LoadError)

	modCfgs := cfgMgr.HyperdriveConfiguration.Modules
	require.Len(t, modCfgs, 1)
	modCfg := modCfgs[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	require.Equal(t, internal_test.ExampleDescriptor.Name, modCfg.Descriptor.Name)

	// Create a default instance of the module config
	require.Len(t, cfgInstance.Modules, 0)
	modInstance := &modconfig.ModuleInstance{
		Enabled: false,
	}
	modInstance.Settings.CreateSettingsFromMetadata(modCfg.Configuration)
	cfgInstance.Modules[modCfg.Descriptor.GetFullyQualifiedModuleName()] = modInstance
	isCfgLoaded = true
	t.Log("Module config loaded successfully")
}

func TestSerialization(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	TestLoadModuleConfigs(t)
	mod := cfgInstance.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	mod.Enabled = true

	// Do some simple tweaks
	cfgInstance.ClientTimeout = 10
	modCfg := mod.Settings.GetSettings()
	floatParam, err := modCfg.GetParameter("exampleFloat")
	require.NoError(t, err)
	err = floatParam.SetValue(80.0)
	serverSection, err := modCfg.GetSection("server")
	require.NoError(t, err)
	portParam, err := serverSection.GetParameter("port")
	require.NoError(t, err)
	err = portParam.SetValue(8085)
	require.NoError(t, err)
	t.Log("Configs modified")

	// Process the configs to make sure they're good
	processResults, err := cfgMgr.ProcessModuleConfigurations(cfgInstance)
	require.NoError(t, err)
	for _, result := range processResults {
		require.NoError(t, result.ProcessError)
		require.Empty(t, result.Issues)
	}
	t.Log("Module config processed successfully")

	// Save everything
	cfgPath := filepath.Join(internal_test.UserDir, config.ConfigFilename)
	err = cfgInstance.SaveToFile(cfgPath)
	require.NoError(t, err)
	t.Log("Main config saved to file")

	// Load the config back in
	newCfg, err := cfgMgr.HyperdriveConfiguration.CreateInstanceFromFile(cfgPath, internal_test.SystemDir)
	require.NoError(t, err)
	require.Equal(t, cfgInstance.ProjectName, newCfg.ProjectName)
	require.Equal(t, cfgInstance.ClientTimeout, newCfg.ClientTimeout)
	t.Log("Main config loaded from file")

	// Load the module configs back in
	newModCfgs := newCfg.Modules
	require.Len(t, newModCfgs, 1)
	newModCfg := newModCfgs[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	require.Equal(t, mod.Enabled, newModCfg.Enabled)
	newSettings := newModCfg.Settings.GetSettings()
	newFloat, err := newSettings.GetParameter("exampleFloat")
	require.NoError(t, err)
	require.Equal(t, floatParam.GetValue(), newFloat.GetValue())
	newServerSection, err := newSettings.GetSection("server")
	require.NoError(t, err)
	newPort, err := newServerSection.GetParameter("port")
	require.NoError(t, err)
	require.Equal(t, portParam.GetValue(), newPort.GetValue())
	t.Log("Module configs loaded from file")
}

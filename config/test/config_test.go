package config_test

import (
	"path/filepath"
	"testing"

	"github.com/nodeset-org/hyperdrive/config"
	internal_test "github.com/nodeset-org/hyperdrive/internal/test"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
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
	results, err := modMgr.LoadModuleInfo(false)
	require.NoError(t, err)
	require.Len(t, results, 1)

	for _, result := range results {
		require.NoError(t, result.LoadError)
		cfgMgr.HyperdriveConfiguration.Modules[result.Info.Descriptor.GetFullyQualifiedModuleName()] = result.Info.ModuleInfo
	}

	modCfgs := cfgMgr.HyperdriveConfiguration.Modules
	require.Len(t, modCfgs, 1)
	modCfg := modCfgs[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	require.Equal(t, internal_test.ExampleDescriptor.Name, modCfg.Descriptor.Name)

	// Create a default instance of the module config
	modSettings := modconfig.CreateModuleSettings(modCfg.Configuration)
	require.Len(t, cfgInstance.Modules, 0)
	modInstance := &modconfig.ModuleInstance{
		Enabled:  false,
		Version:  internal_test.ExampleDescriptor.Version.String(),
		Settings: modSettings.SerializeToMap(),
	}
	require.NoError(t, err)
	modInstance.SetSettings(modSettings)
	cfgInstance.Modules[modCfg.Descriptor.GetFullyQualifiedModuleName()] = modInstance
	isCfgLoaded = true
	t.Log("Module config loaded successfully")
}

func TestSerialization(t *testing.T) {
	err := deleteConfigs()
	require.NoError(t, err)
	TestLoadModuleConfigs(t)
	modCfg := cfgMgr.HyperdriveConfiguration.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	modInstance := cfgInstance.Modules[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	modInstance.Enabled = true
	modSettings, err := modInstance.CreateSettingsFromMetadata(modCfg.Configuration)
	require.NoError(t, err)

	// Create a copy of the config without any modified settings
	oldInstance := cfgInstance.CreateCopy()

	// Do some simple tweaks
	cfgInstance.ClientTimeout = 10
	floatParam, err := modSettings.GetParameter("exampleFloat")
	require.NoError(t, err)
	err = floatParam.SetValue(80.0)
	serverSection, err := modSettings.GetSection("server")
	require.NoError(t, err)
	portParam, err := serverSection.GetParameter("port")
	require.NoError(t, err)
	err = portParam.SetValue(8085)
	require.NoError(t, err)
	portMode, err := serverSection.GetParameter("portMode")
	require.NoError(t, err)
	err = portMode.SetValue("open")
	require.NoError(t, err)
	modInstance.SetSettings(modSettings)
	t.Log("Configs modified")

	// Process the configs to make sure they're good
	processResults, err := cfgMgr.ProcessModuleSettings(modMgr, oldInstance, cfgInstance)
	require.NoError(t, err)
	for _, result := range processResults {
		require.NoError(t, result.ProcessError)
		require.Empty(t, result.Issues)
		require.Len(t, result.Ports, 1)
		require.Equal(t, uint16(8085), result.Ports["server/port"])
		require.Len(t, result.ServicesToRestart, 1)
		require.Equal(t, "example", result.ServicesToRestart[0])
	}
	t.Log("Module config processed successfully")

	// Save everything
	cfgPath := filepath.Join(internal_test.UserDir, config.SettingsFilename)
	err = cfgInstance.SaveToFile(cfgPath)
	require.NoError(t, err)
	t.Log("Main config saved to file")

	// Load the config back in
	newCfg, err := cfgMgr.HyperdriveConfiguration.LoadSettingsFromFile(cfgPath)
	require.NoError(t, err)
	require.Equal(t, cfgInstance.ProjectName, newCfg.ProjectName)
	require.Equal(t, cfgInstance.ClientTimeout, newCfg.ClientTimeout)
	t.Log("Main config loaded from file")

	// Load the module configs back in
	newModCfgs := newCfg.Modules
	require.Len(t, newModCfgs, 1)
	newModInstance := newModCfgs[internal_test.ExampleDescriptor.GetFullyQualifiedModuleName()]
	require.Equal(t, modInstance.Enabled, newModInstance.Enabled)
	newModSettings, err := newModInstance.CreateSettingsFromMetadata(modCfg.Configuration)
	require.NoError(t, err)

	// Make sure the settings were loaded properly
	newFloat, err := newModSettings.GetParameter("exampleFloat")
	require.NoError(t, err)
	require.Equal(t, floatParam.GetValue(), newFloat.GetValue())
	newServerSection, err := newModSettings.GetSection("server")
	require.NoError(t, err)
	newPort, err := newServerSection.GetParameter("port")
	require.NoError(t, err)
	require.Equal(t, portParam.GetValue(), newPort.GetValue())
	t.Log("Module configs loaded from file")
}

package config_test

import (
	"path/filepath"
	"testing"

	internal_test "github.com/nodeset-org/hyperdrive/internal/test"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/stretchr/testify/require"
)

func TestLoadModuleConfigs(t *testing.T) {
	results, err := cfgMgr.LoadModuleInfo(cfgInstance.ProjectName)
	require.NoError(t, err)

	require.Len(t, results, 1)
	require.NoError(t, results[0].LoadError)

	modCfgs := cfgMgr.HyperdriveConfiguration.ModuleInfo
	require.Len(t, modCfgs, 1)
	modCfg := modCfgs[exampleDescriptor.GetFullyQualifiedModuleName()]
	require.Equal(t, exampleDescriptor.Name, modCfg.Descriptor.Name)
	t.Log("Module config loaded successfully")
}

func TestSerialization(t *testing.T) {
	TestLoadModuleConfigs(t)
	var mod *config.HyperdriveModuleInstanceInfo
	for _, m := range cfgInstance.Modules {
		mod = m
		break
	}
	mod.Enabled = true

	// Do some simple tweaks
	cfgInstance.ClientTimeout = 10
	modCfg := mod.Configuration
	param, err := modCfg.GetParameter("exampleFloat")
	require.NoError(t, err)
	err = param.SetValue(80.0)
	require.NoError(t, err)
	t.Log("Configs modified")

	// Process the configs to make sure they're good
	modCfgs := []*config.HyperdriveModuleInstanceInfo{
		mod,
	}
	processResults, err := cfgMgr.ProcessModuleConfigurations(modCfgs)
	require.NoError(t, err)
	for _, result := range processResults {
		require.NoError(t, result.ProcessError)
		require.Empty(t, result.Issues)
	}
	t.Log("Module config processed successfully")

	// Save everything
	cfgPath := filepath.Join(internal_test.UserDir, config.ConfigFilename)
	err = cfgMgr.SaveInstanceToFile(cfgPath, cfgInstance)
	require.NoError(t, err)
	t.Log("Main config saved to file")

	// Load the config back in
	newCfg, err := cfgMgr.LoadInstanceFromFile(cfgPath, internal_test.SystemDir)
	require.NoError(t, err)
	require.Equal(t, cfgInstance.ProjectName, newCfg.ProjectName)
	require.Equal(t, cfgInstance.ClientTimeout, newCfg.ClientTimeout)
	t.Log("Main config loaded from file")

	// Load the module configs back in
	newModCfgs := newCfg.Modules
	require.Len(t, newModCfgs, 1)
	newModCfg, exists := newModCfgs[exampleDescriptor.GetFullyQualifiedModuleName()]
	require.True(t, exists)
	require.Equal(t, mod.Enabled, newModCfg.Enabled)
	newParam, err := newModCfg.Configuration.GetParameter("exampleFloat")
	require.NoError(t, err)
	require.Equal(t, param.GetValue(), newParam.GetValue())
	t.Log("Module configs loaded from file")
}

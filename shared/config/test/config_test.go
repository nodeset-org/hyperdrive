package config_test

import (
	"path/filepath"
	"testing"

	internal_test "github.com/nodeset-org/hyperdrive/internal/test"
	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/stretchr/testify/require"
)

func TestLoadModuleConfigs(t *testing.T) {
	err := hdCfg.LoadModuleConfigs()
	require.NoError(t, err)

	modCfgs := hdCfg.GetModuleConfigs()
	require.Len(t, modCfgs, 1)
	modCfg := modCfgs[0]
	require.NoError(t, modCfg.LoadError)
	require.Equal(t, modules.DescriptorString("example-module"), modCfg.Descriptor.Name)
	t.Log("Module config loaded successfully")
}

func TestSerialization(t *testing.T) {
	TestLoadModuleConfigs(t)
	hdCfg.GetModuleConfigs()[0].Enabled = true

	// Do some simple tweaks
	hdCfg.ClientTimeout.Value = 10
	modCfg := hdCfg.GetModuleConfigs()[0]
	err := modCfg.Config.GetParameters()[3].SetValue(80.0)
	require.NoError(t, err)
	t.Log("Configs modified")

	// Process the configs to make sure they're good
	processResults, err := hdCfg.ProcessModuleConfigs()
	require.NoError(t, err)
	for _, result := range processResults {
		require.NoError(t, result.ProcessError)
		require.Empty(t, result.Issues)
	}
	t.Log("Module config processed successfully")

	// Save everything
	cfgPath := filepath.Join(internal_test.UserDir, "user-settings.yml")
	err = hdCfg.SaveToFile(cfgPath)
	require.NoError(t, err)
	t.Log("Main config saved to file")

	saveErrors, err := hdCfg.SaveModuleConfigs()
	require.NoError(t, err)
	for _, saveErr := range saveErrors {
		require.NoError(t, saveErr)
	}
	t.Log("Module configs saved to disk")

	// Load the config back in
	newCfg, err := config.LoadFromFile(cfgPath, internal_test.SystemDir)
	require.NoError(t, err)
	require.Equal(t, hdCfg.ProjectName.Value, newCfg.ProjectName.Value)
	require.Equal(t, hdCfg.ClientTimeout.Value, newCfg.ClientTimeout.Value)
	t.Log("Main config loaded from file")

	// Load the module configs back in
	err = newCfg.LoadModuleConfigs()
	require.NoError(t, err)
	newModCfgs := newCfg.GetModuleConfigs()
	require.Len(t, newModCfgs, 1)
	newModCfg := newModCfgs[0]
	require.NoError(t, newModCfg.LoadError)
	require.Equal(t, modCfg.Descriptor.Name, newModCfg.Descriptor.Name)
	require.Equal(t, modCfg.Enabled, newModCfg.Enabled)
	require.Equal(t, modCfg.Config.GetParameters()[3].GetValueAsAny(), newModCfg.Config.GetParameters()[3].GetValueAsAny())
	t.Log("Module configs loaded from file")
}

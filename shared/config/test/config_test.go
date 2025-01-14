package config_test

import (
	"testing"

	"github.com/nodeset-org/hyperdrive/modules"
	"github.com/stretchr/testify/require"
)

func TestLoadModuleConfigs(t *testing.T) {
	err := hdCfg.LoadModuleConfigs()
	require.NoError(t, err)

	require.Len(t, hdCfg.ModuleConfigs, 1)
	modCfg := hdCfg.ModuleConfigs[0]
	require.NoError(t, modCfg.ConfigLoadError)
	require.Equal(t, modules.DescriptorString("example-module"), modCfg.Descriptor.Name)
	t.Log("Module config loaded successfully")
}

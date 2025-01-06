package adapter_test

import (
	"context"
	"testing"

	"github.com/nodeset-org/hyperdrive/modules/config"
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

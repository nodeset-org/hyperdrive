package adapter

import (
	"context"
	"fmt"

	"github.com/nodeset-org/hyperdrive/modules/config"
)

const (
	GetConfigMetadataCommandString string = HyperdriveModuleCommand + " get-config-metadata"
)

// Get the module config metadata from the adapter
func (c *AdapterClient) GetConfigMetadata(ctx context.Context) (config.IModuleConfiguration, error) {
	configMap := map[string]any{}

	// Get the config from the adapter
	err := runCommand[struct{}](c, ctx, GetConfigMetadataCommandString, nil, &configMap)
	if err != nil {
		return nil, fmt.Errorf("error getting configuration metadata: %w", err)
	}

	// Unmarshal the config from the map
	response, err := config.UnmarshalConfigurationFromMap(configMap)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling configuration metadata: %w", err)
	}
	return response, nil
}

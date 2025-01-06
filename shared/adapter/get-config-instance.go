package adapter

import (
	"context"
	"fmt"
)

const (
	GetConfigInstanceCommandString string = HyperdriveModuleCommand + " get-config-instance"
)

// Get the module config instance from the adapter
func (c *AdapterClient) GetConfigInstance(ctx context.Context) (map[string]any, error) {
	request := &KeyedRequest{
		Key: c.key,
	}
	configMap := map[string]any{}

	// Get the config from the adapter
	err := runCommand(c, ctx, GetConfigInstanceCommandString, request, &configMap)
	if err != nil {
		return nil, fmt.Errorf("error getting configuration instance: %w", err)
	}
	return configMap, nil
}

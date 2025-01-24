package adapter

import (
	"context"
	"fmt"
)

const (
	SetConfigCommandString string = HyperdriveModuleCommand + " set-config"
)

// Request format for `set-config`
type SetConfigRequest struct {
	KeyedRequest

	// The config instance to process
	Config map[string]any `json:"config"`
}

// Have the adapter set the module config
func (c *AdapterClient) SetConfig(ctx context.Context, instance map[string]any) error {
	request := &SetConfigRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Config: instance,
	}
	err := runCommand[SetConfigRequest, struct{}](c, ctx, SetConfigCommandString, request, nil)
	if err != nil {
		return fmt.Errorf("error processing configuration: %w", err)
	}
	return nil
}

package adapter

import (
	"context"
	"fmt"
)

const (
	SetSettingsCommandString string = HyperdriveModuleCommand + " set-settings"
)

// Request format for `set-settings`
type SetSettingsRequest struct {
	KeyedRequest

	// The Hyperdrive config to process
	Settings map[string]any `json:"settings"`
}

// Have the adapter set the module settings based on the full Hyperdriver configuration.
func (c *AdapterClient) SetSettings(ctx context.Context, settings map[string]any) error {
	request := &SetSettingsRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Settings: settings,
	}
	err := runCommand[SetSettingsRequest, struct{}](c, ctx, SetSettingsCommandString, request, nil)
	if err != nil {
		return fmt.Errorf("error setting module settings: %w", err)
	}
	return nil
}

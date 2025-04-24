package adapter

import (
	"context"
	"fmt"

	"github.com/nodeset-org/hyperdrive/config"
)

const (
	SetSettingsCommandString string = HyperdriveModuleCommand + " set-settings"
)

// Request format for `set-settings`
type SetSettingsRequest struct {
	KeyedRequest

	// The Hyperdrive config to process
	Settings *config.HyperdriveSettings `json:"settings"`
}

// Have the adapter set the module settings based on the full Hyperdrive configuration.
func (c *AdapterClient) SetSettings(ctx context.Context, settings *config.HyperdriveSettings) error {
	request := &SetSettingsRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Settings: settings,
	}
	err := RunCommand[SetSettingsRequest, struct{}](c, ctx, SetSettingsCommandString, request, nil)
	if err != nil {
		return fmt.Errorf("error setting module settings: %w", err)
	}
	return nil
}

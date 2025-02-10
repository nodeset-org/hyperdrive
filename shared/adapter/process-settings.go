package adapter

import (
	"context"
	"fmt"
	//"github.com/nodeset-org/hyperdrive/shared/config"
)

const (
	ProcessSettingsCommandString string = HyperdriveModuleCommand + " process-settings"
)

// Request format for `process-settings`
type ProcessSettingsRequest struct {
	// The Hyperdrive config settings to process
	Settings map[string]any `json:"settings"`
}

// Response format for `process-settings`
type ProcessSettingsResponse struct {
	// A list of errors that occurred during processing, if any
	Errors []string `json:"errors"`

	// A list of ports that will be exposed
	Ports map[string]uint16 `json:"ports"`
}

// Have the adapter process the module settings based on the full Hyperdriver configuration.
func (c *AdapterClient) ProcessSettings(ctx context.Context, settings map[string]any) (ProcessSettingsResponse, error) {
	request := &ProcessSettingsRequest{
		Settings: settings,
	}
	response := ProcessSettingsResponse{}
	err := runCommand(c, ctx, ProcessSettingsCommandString, request, &response)
	if err != nil {
		return response, fmt.Errorf("error processing module settings: %w", err)
	}
	return response, nil
}

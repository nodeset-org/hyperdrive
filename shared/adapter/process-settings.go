package adapter

import (
	"context"
	"fmt"

	"github.com/nodeset-org/hyperdrive/config"
)

const (
	ProcessSettingsCommandString string = HyperdriveModuleCommand + " process-settings"
)

// Request format for `process-settings`
type ProcessSettingsRequest struct {
	// The current config settings
	CurrentSettings *config.HyperdriveSettings `json:"currentSettings"`

	// The new (proposed) config settings
	NewSettings *config.HyperdriveSettings `json:"newSettings"`
}

// Response format for `process-settings`
type ProcessSettingsResponse struct {
	// A list of errors that occurred during processing, if any
	Errors []string `json:"errors"`

	// A list of ports that will be exposed
	Ports map[string]uint16 `json:"ports"`

	// A list of services that need to be restarted as a result of the new settings
	ServicesToRestart []string `json:"servicesToRestart"`
}

// Have the adapter process the module settings based on the full Hyperdriver configuration.
func (c *AdapterClient) ProcessSettings(ctx context.Context, oldSettings *config.HyperdriveSettings, newSettings *config.HyperdriveSettings) (ProcessSettingsResponse, error) {
	request := &ProcessSettingsRequest{
		CurrentSettings: oldSettings,
		NewSettings:     newSettings,
	}
	response := ProcessSettingsResponse{}
	err := runCommand(c, ctx, ProcessSettingsCommandString, request, &response)
	if err != nil {
		return response, fmt.Errorf("error processing module settings: %w", err)
	}
	return response, nil
}

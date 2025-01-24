package adapter

import (
	"context"
	"fmt"
	//"github.com/nodeset-org/hyperdrive/shared/config"
)

const (
	ProcessConfigCommandString string = HyperdriveModuleCommand + " process-config"
)

// Request format for `process-config`
type ProcessConfigRequest struct {
	KeyedRequest

	// The config instance to process
	Config map[string]any `json:"config"`
}

// Response format for `process-config`
type ProcessConfigResponse struct {
	// A list of errors that occurred during processing, if any
	Errors []string `json:"errors"`

	// A list of ports that will be exposed
	Ports map[string]uint16 `json:"ports"`
}

// Have the adapter process the module config.
func (c *AdapterClient) ProcessConfig(ctx context.Context, instance map[string]any) (ProcessConfigResponse, error) {
	request := &ProcessConfigRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Config: instance,
	}
	response := ProcessConfigResponse{}
	err := runCommand(c, ctx, ProcessConfigCommandString, request, &response)
	if err != nil {
		return response, fmt.Errorf("error processing configuration: %w", err)
	}
	return response, nil
}

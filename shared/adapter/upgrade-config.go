package adapter

import (
	"context"
	"fmt"

	"github.com/nodeset-org/hyperdrive/modules/config"
)

const (
	UpgradeConfigCommandString string = HyperdriveModuleCommand + " upgrade-config"
)

// Request format for `upgrade-config`
type UpgradeConfigRequest struct {
	KeyedRequest

	// The config instance to process
	Config map[string]any `json:"config"`
}

// Have the adapter process the module config.
func (c *AdapterClient) UpgradeConfig(ctx context.Context, instance map[string]any) (*config.ModuleInstance, error) {
	request := &UpgradeConfigRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Config: instance,
	}
	response := config.ModuleInstance{}
	err := runCommand(c, ctx, UpgradeConfigCommandString, request, &response)
	if err != nil {
		return &response, fmt.Errorf("error upgrading configuration: %w", err)
	}
	return &response, nil
}

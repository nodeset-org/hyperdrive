package adapter

import (
	"context"
	"fmt"

	"github.com/nodeset-org/hyperdrive/modules/config"
)

const (
	UpgradeInstanceCommandString string = HyperdriveModuleCommand + " upgrade-instance"
)

// Request format for `upgrade-instance`
type UpgradeInstanceRequest struct {
	KeyedRequest

	// The currently saved instance to upgrade
	Instance *config.ModuleInstance `json:"instance"`
}

// Send an instance upgrade request to the adapter.
func (c *AdapterClient) UpgradeInstance(ctx context.Context, instance *config.ModuleInstance) (*config.ModuleInstance, error) {
	request := &UpgradeInstanceRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Instance: instance,
	}
	response := config.ModuleInstance{}
	err := runCommand(c, ctx, UpgradeInstanceCommandString, request, &response)
	if err != nil {
		return &response, fmt.Errorf("error upgrading settings: %w", err)
	}
	return &response, nil
}

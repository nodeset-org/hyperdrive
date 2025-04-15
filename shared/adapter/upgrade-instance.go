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
	// The currently saved instance to upgrade
	Instance *config.ModuleInstance `json:"instance"`
}

// Send an instance upgrade request to the adapter.
func (c *AdapterClient) UpgradeInstance(ctx context.Context, instance *config.ModuleInstance) (*config.ModuleInstance, error) {
	request := &UpgradeInstanceRequest{
		Instance: instance,
	}
	response := config.ModuleInstance{}
	err := RunCommand(c, ctx, UpgradeInstanceCommandString, request, &response)
	if err != nil {
		return &response, fmt.Errorf("error upgrading settings: %w", err)
	}
	return &response, nil
}

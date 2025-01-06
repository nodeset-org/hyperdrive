package adapter

import (
	"context"
)

const (
	RunCommandString string = HyperdriveModuleCommand + " run"
)

// Run a command on the adapter
func (c *AdapterClient) Run(ctx context.Context, command string) error {
	// TODO
	return nil
}

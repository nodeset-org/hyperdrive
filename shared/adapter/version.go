package adapter

import (
	"context"
	"fmt"
)

const (
	VersionCommandString string = HyperdriveModuleCommand + " version"
)

// Response for the version command
type VersionResponse struct {
	Version string `json:"version"`
}

// Get the version of the adapter
func (c *AdapterClient) GetVersion(ctx context.Context) (string, error) {
	var version VersionResponse
	err := RunCommand[struct{}](c, ctx, VersionCommandString, nil, &version)
	if err != nil {
		return "", fmt.Errorf("error getting version: %w", err)
	}
	return version.Version, nil
}

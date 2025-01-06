package adapter

import (
	"context"
	"fmt"
)

const (
	GetContainersCommandString string = HyperdriveModuleCommand + " get-containers"
)

// Response format for `get-containers`
type GetContainersResponse struct {
	// The list of containers owned by this module
	Containers []string `json:"containers"`
}

// Get the list of containers owned by this module
func (c *AdapterClient) GetContainers(ctx context.Context) ([]string, error) {
	request := &KeyedRequest{
		Key: c.key,
	}
	response := &GetContainersResponse{}
	err := runCommand(c, ctx, GetContainersCommandString, request, response)
	if err != nil {
		return nil, fmt.Errorf("error getting containers: %w", err)
	}
	return response.Containers, nil
}

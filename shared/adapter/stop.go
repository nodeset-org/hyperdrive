package adapter

import (
	"context"
)

const (
	StopCommandString string = HyperdriveModuleCommand + " stop"
)

type StopRequest struct {
	KeyedRequest

	// The compose project name
	ComposeProjectName string `json:"composeProjectName"`

	// The services to stop. If empty, all services will be stopped.
	Services []string `json:"services"`
}

// Have the adapter stop the module.
func (c *AdapterClient) Stop(ctx context.Context, composeProjectName string, services []string) error {
	if services == nil {
		services = []string{}
	}
	request := &StopRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		ComposeProjectName: composeProjectName,
		Services:           services,
	}
	err := runCommand[StopRequest, struct{}](c, ctx, StopCommandString, request, nil)
	if err != nil {
		return err
	}
	return nil
}

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
}

// Have the adapter stop the module.
func (c *AdapterClient) Stop(ctx context.Context, composeProjectName string) error {
	request := &StopRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		ComposeProjectName: composeProjectName,
	}
	err := runCommand[StopRequest, struct{}](c, ctx, StopCommandString, request, nil)
	if err != nil {
		return err
	}
	return nil
}

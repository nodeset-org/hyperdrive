package adapter

import (
	"context"

	"github.com/nodeset-org/hyperdrive/shared/config"
)

const (
	StartCommandString string = HyperdriveModuleCommand + " start"
)

type StartRequest struct {
	KeyedRequest

	// The Hyperdrive config to process
	Settings *config.HyperdriveSettings `json:"settings"`

	// The compose project name
	ComposeProjectName string `json:"composeProjectName"`
}

// Have the adapter start the module based on the full Hyperdrive configuration.
func (c *AdapterClient) Start(ctx context.Context, settings *config.HyperdriveSettings, composeProjectName string) error {
	request := &StartRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Settings:           settings,
		ComposeProjectName: composeProjectName,
	}
	err := runCommand[StartRequest, struct{}](c, ctx, StartCommandString, request, nil)
	if err != nil {
		return err
	}
	return nil
}

package adapter

import (
	"context"

	"github.com/nodeset-org/hyperdrive/config"
)

const (
	StartCommandString string = HyperdriveModuleCommand + " start"
)

type StartRequest struct {
	KeyedRequest

	// The Hyperdrive config to process
	Settings *config.HyperdriveSettings `json:"settings"`
}

// Have the adapter start the module based on the full Hyperdrive configuration.
func (c *AdapterClient) Start(ctx context.Context, settings *config.HyperdriveSettings) error {
	request := &StartRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Settings: settings,
	}
	err := RunCommand[StartRequest, struct{}](c, ctx, StartCommandString, request, nil)
	if err != nil {
		return err
	}
	return nil
}

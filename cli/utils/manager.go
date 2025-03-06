package utils

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/management"
	"github.com/urfave/cli/v2"
)

// Create a new Hyperdrive manager from the CLI context
func NewHyperdriveManagerFromCtx(c *cli.Context) (*management.HyperdriveManager, error) {
	hdCtx := management.GetHyperdriveContext(c)
	if hdCtx == nil {
		return nil, fmt.Errorf("Hyperdrive CLI context has not been created")
	}
	return management.NewHyperdriveManager(hdCtx)
}

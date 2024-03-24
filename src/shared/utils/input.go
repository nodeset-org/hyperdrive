package utils

import (
	"github.com/rocket-pool/node-manager-core/utils/input"
	"github.com/urfave/cli/v2"
)

// Validate command argument count
func ValidateArgCount(c *cli.Context, expectedCount int) error {
	return input.ValidateArgCount(c.Args().Len(), expectedCount)
}

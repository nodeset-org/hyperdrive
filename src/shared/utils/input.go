package utils

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Validate command argument count
func ValidateArgCount(c *cli.Context, expectedCount int) error {
	argCount := c.Args().Len()
	if argCount != expectedCount {
		return fmt.Errorf("incorrect argument count; expected %d but have %d", expectedCount, argCount)
	}
	return nil
}

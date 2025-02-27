package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/nodeset-org/hyperdrive/cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/cli/input"
	"github.com/urfave/cli/v2"
)

// Validate command argument count
func ValidateArgCount(c *cli.Context, expectedCount int) {
	err := input.ValidateArgCount(c.Args().Len(), expectedCount)
	if err != nil {
		// Handle invalid arg count
		var argCountErr *input.InvalidArgCountError
		if errors.As(err, &argCountErr) {
			fmt.Fprintf(os.Stderr, "%s%s%s\n\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			cli.ShowSubcommandHelpAndExit(c, 1)
		}

		// Handle other errors
		fmt.Fprintf(os.Stderr, "%s%s%s\n\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
		os.Exit(1)
	}
}

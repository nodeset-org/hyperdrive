package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/cli/input"
	"github.com/urfave/cli/v2"
)

// Validate command argument count - only used by the CLI
// TODO: refactor CLI arg validation and move it out of shared
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

// Validate a token type
func ValidateTokenType(name, value string) (string, error) {
	// Check if this is a token address
	// This was taken from the Ethereum library: https://github.com/ethereum/go-ethereum/blob/master/common/types.go
	if strings.HasPrefix(value, "0x") {
		// Remove the 0x prefix
		val := value[2:]

		// Zero pad if it's an odd number of chars
		if len(val)%2 == 1 {
			val = "0" + val
		}

		// Attempt parsing
		_, err := hex.DecodeString(val)
		if err != nil {
			return "", fmt.Errorf("Invalid %s '%s' - could not parse address: %w", name, value, err)
		}

		// If it passes, return the original value
		return value, nil
	}

	// Not a token address, check against the well-known names
	val := strings.ToLower(value)
	if !(val == "eth") {
		return "", fmt.Errorf("Invalid %s '%s' - valid token names are 'ETH'", name, value)
	}
	return val, nil
}

package context

import (
	"math/big"

	"github.com/urfave/cli/v2"
)

const (
	contextMetadataName string = "hd-context"
)

// Context for global settings
type HyperdriveContext struct {
	// The path to the configuration file
	ConfigPath string

	// The max fee for transactions
	MaxFee float64

	// The max priority fee for transactions
	MaxPriorityFee float64

	// The nonce for the first transaction, if set
	Nonce *big.Int

	// True if debug mode is enabled
	DebugEnabled bool

	// True if this is a secure session
	SecureSession bool
}

// Add the Hyperdrive context into a CLI context
func SetHyperdriveContext(c *cli.Context, hdCtx *HyperdriveContext) {
	c.App.Metadata[contextMetadataName] = hdCtx
}

// Get the Hyperdrive context from a CLI context
func GetHyperdriveContext(c *cli.Context) *HyperdriveContext {
	return c.App.Metadata[contextMetadataName].(*HyperdriveContext)
}

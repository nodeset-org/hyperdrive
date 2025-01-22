package context

import (
	"math/big"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	contextMetadataName string = "hd-context"
)

// Context for global settings
type HyperdriveContext struct {
	*InstallationInfo

	// The path to the Hyperdrive user directory
	UserDirPath string

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

	// The address and URL of the API server
	ApiUrl *url.URL

	// The HTTP trace file if tracing is enabled
	HttpTraceFile *os.File
}

// Creates a new Hyperdrive context. If installationInfo is nil, a new one will be created using the system configuration.
func NewHyperdriveContext(userDirPath string, installationInfo *InstallationInfo) *HyperdriveContext {
	if installationInfo == nil {
		installationInfo = NewInstallationInfo()
	}
	return &HyperdriveContext{
		UserDirPath:      userDirPath,
		InstallationInfo: installationInfo,
	}
}

// Add the Hyperdrive context into a CLI context
func SetHyperdriveContext(c *cli.Context, hd *HyperdriveContext) {
	c.App.Metadata[contextMetadataName] = hd
}

// Get the Hyperdrive context from a CLI context
func GetHyperdriveContext(c *cli.Context) *HyperdriveContext {
	return c.App.Metadata[contextMetadataName].(*HyperdriveContext)
}

package management

import (
	"net/url"
	"os"

	"path/filepath"

	"github.com/urfave/cli/v2"
)

const (
	// Environment variable to set the system path for unit tests
	TestSystemDirEnvVar string = "HYPERDRIVE_TEST_SYSTEM_DIR"

	// Subfolders under the system dir
	TemplatesDir      string = "templates"
	OverrideSourceDir string = "override"
	ModulesDir        string = "modules"

	// Key for getting the Hyperdrive context from a CLI context
	contextMetadataName string = "hd-context"
)

// Context for Hyperdrive clients and commands
type HyperdriveContext struct {
	// The path to the Hyperdrive user directory
	UserDirPath string

	// The path to the Hyperdrive system directory
	SystemDirPath string

	// True if debug mode is enabled
	DebugEnabled bool

	// True if this is a secure session
	SecureSession bool

	// The address and URL of the API server
	ApiUrl *url.URL

	// The HTTP trace file if tracing is enabled
	HttpTraceFile *os.File
}

// Creates a new Hyperdrive context.
func NewHyperdriveContext(userDirPath string, systemDirPath string) *HyperdriveContext {
	systemDir := os.Getenv(TestSystemDirEnvVar)
	if systemDir == "" {
		systemDir = systemDirPath
	}

	return &HyperdriveContext{
		UserDirPath:   userDirPath,
		SystemDirPath: systemDir,
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

// The full path of the base templates directory
func (c HyperdriveContext) TemplatesDir() string {
	return filepath.Join(c.SystemDirPath, TemplatesDir)
}

// The full path of the override compose file source directory
func (c HyperdriveContext) OverrideSourceDir() string {
	return filepath.Join(c.SystemDirPath, OverrideSourceDir)
}

// The full path of the modules directory
func (c HyperdriveContext) ModulesDir() string {
	return filepath.Join(c.SystemDirPath, ModulesDir)
}

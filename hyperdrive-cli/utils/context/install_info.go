package context

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	// Environment variable to set the system path for unit tests
	TestSystemDirEnvVar string = "HYPERDRIVE_TEST_SYSTEM_DIR"

	// System dir path for Linux
	linuxSystemDir string = "/usr/share/hyperdrive"

	// Subfolders under the system dir
	scriptsDir        string = "scripts"
	templatesDir      string = "templates"
	overrideSourceDir string = "override"
	networksDir       string = "networks"
)

// Holds information about Hyperdrive's installation on the system
type InstallationInfo struct {
	// The system path for Hyperdrive scripts used in the Docker containers
	ScriptsDir string

	// The system path for Hyperdrive templates
	TemplatesDir string

	// The system path for the source files to put in the user's override directory
	OverrideSourceDir string

	// The system path for built-in network settings and resource definitions
	NetworksDir string
}

// Creates a new installation info instance
func NewInstallationInfo() *InstallationInfo {
	systemDir := os.Getenv(TestSystemDirEnvVar)
	if systemDir == "" {
		switch runtime.GOOS {
		// This is where to add different paths for different OS's like macOS
		default:
			// By default just use the Linux path
			systemDir = linuxSystemDir
		}
	}

	return &InstallationInfo{
		ScriptsDir:        filepath.Join(systemDir, scriptsDir),
		TemplatesDir:      filepath.Join(systemDir, templatesDir),
		OverrideSourceDir: filepath.Join(systemDir, overrideSourceDir),
		NetworksDir:       filepath.Join(systemDir, networksDir),
	}
}

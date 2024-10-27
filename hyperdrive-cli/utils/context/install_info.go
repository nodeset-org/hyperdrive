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
	LinuxSystemDir string = "/usr/share/hyperdrive"

	// Subfolders under the system dir
	ScriptsDir        string = "scripts"
	TemplatesDir      string = "templates"
	OverrideSourceDir string = "override"
	NetworksDir       string = "networks"
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
			systemDir = LinuxSystemDir
		}
	}

	return NewInstallationInfoForSystemDir(systemDir)
}

// Creates a new installation info instance with the given system directory
func NewInstallationInfoForSystemDir(systemDir string) *InstallationInfo {
	return &InstallationInfo{
		ScriptsDir:        filepath.Join(systemDir, ScriptsDir),
		TemplatesDir:      filepath.Join(systemDir, TemplatesDir),
		OverrideSourceDir: filepath.Join(systemDir, OverrideSourceDir),
		NetworksDir:       filepath.Join(systemDir, NetworksDir),
	}
}

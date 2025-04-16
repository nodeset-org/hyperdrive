package tests

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/mholt/archiver/v4"
	"github.com/nodeset-org/hyperdrive-daemon/shared"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/stretchr/testify/require"
)

const (
	autocompleteSrc string = "../../install/autocomplete"
	deploySrc       string = "../../install/deploy"
)

func TestManualInstall(t *testing.T) {
	err := testMgr.RevertToBaseline()
	if err != nil {
		fail("Error reverting to baseline snapshot: %v", err)
	}
	defer handle_panics()

	// Make sure the install package file doesn't exist yet
	oshaDir := testMgr.GetTestDir()
	pkgFilePath := filepath.Join(oshaDir, "hyperdrive-install.tar.xz")
	_, err = os.Stat(pkgFilePath)
	require.True(t, errors.Is(err, fs.ErrNotExist))
	t.Logf("Install package file does not exist as expected: %s", pkgFilePath)

	// Create the install package
	err = createInstallPackage(pkgFilePath)
	require.NoError(t, err)
	t.Logf("Install package created: %s", pkgFilePath)

	// Make sure the install package file exists now
	_, err = os.Stat(pkgFilePath)
	require.NoError(t, err)
	t.Logf("Install package file now exists: %s", pkgFilePath)

	// Run the installer
	installerScriptPath := filepath.Join("..", "..", "install", "install.sh")
	installPath := filepath.Join(oshaDir, "install_usr", "share")
	err = os.MkdirAll(installPath, 0755)
	require.NoError(t, err)
	runtimePath := filepath.Join(oshaDir, "install_var", "lib")
	err = os.MkdirAll(runtimePath, 0755)
	require.NoError(t, err)
	autocompletePath := filepath.Join(oshaDir, "install_usr", "share", "bash-completion", "completions")
	err = os.MkdirAll(autocompletePath, 0755)
	require.NoError(t, err)
	err = client.InstallService(client.InstallOptions{
		RequireEscalation:       false,
		Verbose:                 true,
		NoDeps:                  true,
		Version:                 shared.HyperdriveVersion,
		InstallPath:             installPath,
		RuntimePath:             runtimePath,
		LocalInstallScriptPath:  installerScriptPath,
		LocalInstallPackagePath: pkgFilePath,
		BashCompletionPath:      autocompletePath,
	})
	require.NoError(t, err)
	t.Logf("Service installed")

	// Make sure the files were deployed correctly
	err = filepath.WalkDir(deploySrc, func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err)
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(deploySrc, path)
		require.NoError(t, err)
		destPath := filepath.Join(installPath, relPath)
		_, err = os.Stat(destPath)
		require.NoError(t, err)
		t.Logf("File deployed: %s", destPath)
		return nil
	})
	require.NoError(t, err)

	// Check the autocomplete file
	autocompleteFile := filepath.Join(autocompletePath, "hyperdrive")
	_, err = os.Stat(autocompleteFile)
	require.NoError(t, err)
	t.Logf("File deployed: %s", autocompleteFile)
	t.Logf("All files deployed correctly")
}

// Create the install package
func createInstallPackage(pkgFilePath string) error {
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		deploySrc:       "install/deploy",
		autocompleteSrc: "install/autocomplete",
	})
	if err != nil {
		return fmt.Errorf("error traversing deploy directory: %w", err)
	}

	format := archiver.CompressedArchive{
		Compression: archiver.Xz{},
		Archival:    archiver.Tar{},
	}

	// Make the install package file
	pkgFile, err := os.Create(pkgFilePath)
	if err != nil {
		return fmt.Errorf("error creating install package file [%s]: %w", pkgFilePath, err)
	}

	err = format.Archive(context.Background(), pkgFile, files)
	if err != nil {
		return fmt.Errorf("error creating install package: %w", err)
	}
	return nil
}

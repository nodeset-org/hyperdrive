package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Saves a file to disk in a safe manner by writing to a temporary file and then replacing the original file.
// If something precludes a successful save, such as a full disk, the original file is left untouched.
// Requires the file's parent directory must already exist.
// Adapted from Patches's original Smart Node implementation.
func SafeSaveFile(data []byte, path string, mode os.FileMode) error {
	directory := filepath.Dir(path)
	filename := filepath.Base(path)

	// Create a temporary file in the same directory
	now := time.Now().Local()
	date := now.Format(time.DateOnly)
	timestamp := strings.Replace(now.Format(time.TimeOnly), ":", "-", -1)
	tempFilename := fmt.Sprintf("%s_%s_%s.tmp", filename, date, timestamp)
	tempFile, err := os.CreateTemp(directory, tempFilename)
	if err != nil {
		return fmt.Errorf("error creating temporary file [%s]: %w", tempFilename, err)
	}

	defer func() {
		// Close the file
		_ = tempFile.Close()

		// Remove the temporary file
		_ = os.Remove(tempFile.Name())
	}()

	// Write the data to the temporary file
	_, err = tempFile.Write(data)
	if err != nil {
		return fmt.Errorf("error writing data to temporary file [%s]: %w", tempFilename, err)
	}
	err = tempFile.Close()
	if err != nil {
		return fmt.Errorf("error closing temporary file [%s]: %w", tempFilename, err)
	}

	// Replace the original file with the temporary file
	err = os.Rename(tempFile.Name(), path)
	if err != nil {
		return fmt.Errorf("error replacing file [%s] with temporary file [%s]: %w", filename, tempFilename, err)
	}
	err = os.Chmod(path, mode)
	if err != nil {
		return fmt.Errorf("error setting permissions of file [%s]: %w", filename, err)
	}
	return nil
}

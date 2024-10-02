package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

const (
	// The default length of the API authorization key, in bytes
	DefaultKeyLength int = 32

	// The permissions to set on the API authorization key file
	KeyPermissions fs.FileMode = 0600
)

// Generates a new authorization secret key if it's not already on disk.
// If the key already exists, this does nothing.
func GenerateAuthKeyIfNotPresent(path string, keyLengthInBytes int) error {
	// Check if the file exists
	exists := true
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		exists = false
	} else if err != nil {
		return fmt.Errorf("error checking if key [%s] exists: %w", path, err)
	}

	// Break if it already exists
	if exists {
		return nil
	}

	// Generate the key
	buffer := make([]byte, keyLengthInBytes)
	_, err = rand.Read(buffer)
	if err != nil {
		return fmt.Errorf("error generating random key: %w", err)
	}

	// Write the key to disk
	err = os.WriteFile(path, buffer, KeyPermissions)
	if err != nil {
		return fmt.Errorf("error writing key to [%s]: %w", path, err)
	}
	return nil
}

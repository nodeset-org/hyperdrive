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
	DefaultKeyLength int = 384 / 8

	// The permissions to set on the API authorization key file
	KeyPermissions fs.FileMode = 0600
)

// Generates a new authorization secret key if it's not already on disk.
// If the key already exists, this does nothing.
// NOTE: key length must be 48 bytes or higher for security.
// See https://github.com/nodeset-org/hyperdrive-daemon/pull/32#discussion_r1784899614
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
	if keyLengthInBytes < DefaultKeyLength {
		return fmt.Errorf("key length must be at least %d bytes", DefaultKeyLength)
	}
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

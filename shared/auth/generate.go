package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	// The default length of the API authorization key, in bytes
	// Since the auth manager uses HS384, the default key length is 384 bits (equal to the output size of the underlying hash function).
	// See https://datatracker.ietf.org/doc/html/rfc2104#section-3 for more info.
	DefaultKeyLength int = 384 / 8

	// The permissions to set on the API authorization key file
	KeyPermissions fs.FileMode = 0600

	// The permissions to set on the API authorization key directory
	KeyDirPermissions fs.FileMode = 0700
)

// Generates a new authorization secret key if it's not already on disk.
// If the key already exists, this does nothing.
// NOTE: key length must be 48 bytes (hash size of HS384) or higher for security.
// See https://datatracker.ietf.org/doc/html/rfc2104#section-3
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
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating key directory [%s]: %w", dir, err)
	}
	err = os.WriteFile(path, buffer, KeyDirPermissions)
	if err != nil {
		return fmt.Errorf("error writing key to [%s]: %w", path, err)
	}
	return nil
}

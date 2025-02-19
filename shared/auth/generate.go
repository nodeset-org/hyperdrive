package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
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

// Generates a new random secret key and encodes it as hex; useful for JWT authorization.
func GenerateAuthKey(keyLengthInBytes int) (string, error) {
	buffer := make([]byte, keyLengthInBytes)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("error generating random key: %w", err)
	}

	return hex.EncodeToString(buffer), nil
}

// Writes the given key to the specified file path.
func CreateKeyFile(keyPath string, keyLengthInBytes int) error {
	key, err := GenerateAuthKey(keyLengthInBytes)
	if err != nil {
		return fmt.Errorf("error generating key: %w", err)
	}

	err = os.WriteFile(keyPath, []byte(key), KeyPermissions)
	if err != nil {
		return fmt.Errorf("error writing key file [%s]: %w", keyPath, err)
	}

	return nil
}

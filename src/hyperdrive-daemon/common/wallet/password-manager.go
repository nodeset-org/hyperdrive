package wallet

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

const (
	passwordFileMode fs.FileMode = 0600
)

// Simple class to wrap the node's password file
type PasswordManager struct {
	path string
}

// Creates a new password manager
func NewPasswordManager(path string) *PasswordManager {
	return &PasswordManager{
		path: path,
	}
}

// Gets the password saved on disk. Returns nil if the password file doesn't exist.
func (m *PasswordManager) GetPasswordFromDisk() (string, bool, error) {
	_, err := os.Stat(m.path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", false, nil
	}

	bytes, err := os.ReadFile(m.path)
	if err != nil {
		return "", false, fmt.Errorf("error reading password file [%s]: %w", m.path, err)
	}
	return string(bytes), true, nil
}

// Save the password to disk
func (m *PasswordManager) SavePassword(password string) error {
	err := os.WriteFile(m.path, []byte(password), passwordFileMode)
	if err != nil {
		return fmt.Errorf("error savingpassword to [%s]: %w", err)
	}
	return nil
}

// Delete the password from disk
func (m *PasswordManager) DeletePassword() error {
	err := os.Remove(m.path)
	if err != nil {
		return fmt.Errorf("error deleting password [%s]: %w", m.path, err)
	}
	return nil
}

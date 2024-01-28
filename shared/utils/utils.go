package utils

import (
	"strings"

	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/sethvargo/go-password/password"
)

const (
	hexPrefix string = "0x"
)

// Add a prefix to a hex string if not present
func AddPrefix(value string) string {
	if !strings.HasPrefix(value, hexPrefix) {
		return hexPrefix + value
	}
	return value
}

// Remove a prefix from a hex string if present
func RemovePrefix(value string) string {
	return strings.TrimPrefix(value, hexPrefix)
}

// Check if the node wallet is ready for transacting
func IsWalletReady(status types.WalletStatus) bool {
	return status.Address.HasAddress &&
		status.Wallet.IsLoaded &&
		status.Address.NodeAddress == status.Wallet.WalletAddress
}

// Generates a random password
func GenerateRandomPassword() (string, error) {
	// Generate a random 32-character password
	password, err := password.Generate(32, 6, 6, false, false)
	if err != nil {
		return "", err
	}

	return password, nil
}

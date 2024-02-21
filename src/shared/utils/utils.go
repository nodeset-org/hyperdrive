package utils

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/sethvargo/go-password/password"
)

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

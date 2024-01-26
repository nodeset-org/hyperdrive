package utils

import "github.com/nodeset-org/hyperdrive/shared/types"

// Add a prefix to a hex string if not present
func AddPrefix(value string) string {
	if len(value) < 2 || value[0:2] != "0x" {
		return "0x" + value
	}
	return value
}

// Remove a prefix from a hex string if present
func RemovePrefix(value string) string {
	if len(value) >= 2 && value[0:2] == "0x" {
		return value[2:]
	}
	return value
}

// Check if the node wallet is ready for transacting
func IsWalletReady(status types.WalletStatus) bool {
	return status.HasAddress &&
		status.HasKeystore &&
		status.HasPassword &&
		status.NodeAddress == status.KeystoreAddress
}

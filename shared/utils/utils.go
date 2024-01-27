package utils

import (
	"strings"

	"github.com/nodeset-org/hyperdrive/shared/types"
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
	return status.HasAddress &&
		status.HasKeystore &&
		status.HasPassword &&
		status.NodeAddress == status.KeystoreAddress
}

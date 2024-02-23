package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/wallet"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Convert a derivation path type to an actual path value
func GetDerivationPath(pathType types.DerivationPath) (string, error) {
	// Parse the derivation path
	switch pathType {
	case types.DerivationPath_Default:
		return wallet.DefaultNodeKeyPath, nil
	case types.DerivationPath_LedgerLive:
		return wallet.LedgerLiveNodeKeyPath, nil
	case types.DerivationPath_Mew:
		return wallet.MyEtherWalletNodeKeyPath, nil
	default:
		return "", fmt.Errorf("[%s] is not a valid derivation path type", string(pathType))
	}
}

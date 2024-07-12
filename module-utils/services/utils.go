package services

import (
	"errors"

	"github.com/rocket-pool/node-manager-core/wallet"
)

var (
	//lint:ignore ST1005 These are printed to the user and need to be in proper grammatical format
	ErrNodeAddressNotSet error = errors.New("The node currently does not have an address set. Please run 'hyperdrive wallet init' and try again.")

	//lint:ignore ST1005 These are printed to the user and need to be in proper grammatical format
	ErrNeedPassword error = errors.New("The node has a node wallet on disk but does not have the password for it loaded. Please run `hyperdrive wallet set-password` to load it.")

	//lint:ignore ST1005 These are printed to the user and need to be in proper grammatical format
	ErrWalletLoadFailure error = errors.New("The node has a node wallet and a password on disk but there was an error loading it - perhaps the password is incorrect? Please check the node logs for more information.")

	//lint:ignore ST1005 These are printed to the user and need to be in proper grammatical format
	ErrNoWallet error = errors.New("The node currently does not have a node wallet keystore. Please run 'hyperdrive wallet init' and try again.")

	//lint:ignore ST1005 These are printed to the user and need to be in proper grammatical format
	ErrWalletMismatch error = errors.New("The node's wallet keystore does not match the node address. This node is currently in read-only mode.")
)

func CheckIfWalletReady(status wallet.WalletStatus) error {
	if !status.Address.HasAddress {
		return ErrNodeAddressNotSet
	}
	if !status.Wallet.IsLoaded {
		if status.Wallet.IsOnDisk {
			if !status.Password.IsPasswordSaved {
				return ErrNeedPassword
			}
			return ErrWalletLoadFailure
		}
		return ErrNoWallet
	}
	if status.Wallet.WalletAddress != status.Address.NodeAddress {
		return ErrWalletMismatch
	}
	return nil
}

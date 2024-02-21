package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tyler-smith/go-bip39"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
	eth2util "github.com/wealdtech/go-eth2-util"
)

const (
	EthWithdrawalPrefix byte = 0x01
)

// Convert an address into 0x01-prefixed withdrawal credentials suitable for depositing into Beacon
func GetWithdrawalCredsFromAddress(address common.Address) common.Hash {
	addressBytes := address[:]
	withdrawalCreds := common.BytesToHash(addressBytes)
	withdrawalCreds[0] = EthWithdrawalPrefix // Set it to a 0x01 credential
	return withdrawalCreds
}

// Get a private BLS key from the mnemonic and path.
func GetPrivateKey(mnemonic string, path string) (*eth2types.BLSPrivateKey, error) {
	// Generate seed
	seed := bip39.NewSeed(mnemonic, "")

	// Initialize BLS support
	if err := InitializeBls(); err != nil {
		return nil, fmt.Errorf("Could not initialize BLS library: %w", err)
	}

	// Get private key
	privateKey, err := eth2util.PrivateKeyFromSeedAndPath(seed, path)
	if err != nil {
		return nil, fmt.Errorf("Could not get validator private key for [%s]: %w", path, err)
	}

	// Return
	return privateKey, nil
}

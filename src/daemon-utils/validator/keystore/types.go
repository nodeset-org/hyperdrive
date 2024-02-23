package keystore

import (
	"github.com/nodeset-org/eth-utils/beacon"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

const (
	DirectEIPVersion string = "EIP-2335"
)

// Validator keystore manager interface
type IKeystoreManager interface {
	// Store a validator key on disk
	StoreValidatorKey(key *eth2types.BLSPrivateKey, derivationPath string) error

	// Load a validator key from disk corresponding to the provided pubkey
	LoadValidatorKey(pubkey beacon.ValidatorPubkey) (*eth2types.BLSPrivateKey, error)

	// Get the path of the keystore directory managed by this manager
	GetKeystoreDir() string
}

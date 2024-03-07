package services

import (
	"fmt"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/node/validator/keystore"
	types "github.com/wealdtech/go-eth2-types/v2"
)

type ValidatorManager struct {
	keystoreManagers map[string]keystore.IKeystoreManager
}

func NewValidatorManager(moduleDir string) *ValidatorManager {
	// Get the validator storage path
	validatorPath := filepath.Join(moduleDir, config.ValidatorsDirectory)

	mgr := &ValidatorManager{
		keystoreManagers: map[string]keystore.IKeystoreManager{
			"lighthouse": keystore.NewLighthouseKeystoreManager(validatorPath),
			"lodestar":   keystore.NewLodestarKeystoreManager(validatorPath),
			"nimbus":     keystore.NewNimbusKeystoreManager(validatorPath),
			"prysm":      keystore.NewPrysmKeystoreManager(validatorPath),
			"teku":       keystore.NewTekuKeystoreManager(validatorPath),
		},
	}
	return mgr
}

func (m *ValidatorManager) StoreKey(key *types.BLSPrivateKey, derivationPath string) error {
	for name, mgr := range m.keystoreManagers {
		err := mgr.StoreValidatorKey(key, derivationPath)
		if err != nil {
			pubkey := beacon.ValidatorPubkey(key.PublicKey().Marshal())
			return fmt.Errorf("error saving validator key %s (path %s) to the %s keystore: %w", pubkey.HexWithPrefix(), derivationPath, name, err)
		}
	}
	return nil
}

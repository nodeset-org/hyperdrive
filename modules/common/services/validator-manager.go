package services

import (
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/modules/common/validator/keystore"
	modconfig "github.com/nodeset-org/hyperdrive/shared/config/modules"
)

type ValidatorManager struct {
	keystoreManagers map[string]keystore.IKeystoreManager
}

func NewValidatorManager(moduleDir string) *ValidatorManager {
	// Get the validator storage path
	validatorPath := filepath.Join(moduleDir, modconfig.ValidatorDir)

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

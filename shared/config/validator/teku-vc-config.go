package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Tags
	tekuVcTagTest string = "consensys/teku:24.1.1"
	tekuVcTagProd string = "consensys/teku:24.1.1"
)

// Configuration for Teku
type TekuVcConfig struct {
	// The Docker Hub tag for the Teku VC
	ContainerTag types.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags types.Parameter[string]
}

// Generates a new Teku VC configuration
func NewTekuVcConfig() *TekuVcConfig {
	return &TekuVcConfig{
		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
				Description:        "The tag name of the Teku container on Docker Hub you want to use for the Validator Client.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    tekuVcTagProd,
				types.Network_HoleskyDev: tekuVcTagTest,
				types.Network_Holesky:    tekuVcTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Teku Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}
}

// Get the title for the config
func (cfg *TekuVcConfig) GetTitle() string {
	return "Teku Validator Client"
}

// Get the parameters for this config
func (cfg *TekuVcConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *TekuVcConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

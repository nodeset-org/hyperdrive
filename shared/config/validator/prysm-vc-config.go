package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Tags
	prysmVcTagTest string = "rocketpool/prysm:v4.2.0"
	prysmVcTagProd string = "rocketpool/prysm:v4.2.0"
)

// Configuration for the Prysm VC
type PrysmVcConfig struct {
	// The Docker Hub tag for the Prysm BN
	ContainerTag types.Parameter[string]

	// Custom command line flags for the BN
	AdditionalFlags types.Parameter[string]
}

// Generates a new Prysm VC configuration
func NewPrysmVcConfig() *PrysmVcConfig {
	return &PrysmVcConfig{
		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
				Description:        "The tag name of the Prysm container on Docker Hub you want to use for the Validator Client.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    prysmVcTagProd,
				types.Network_HoleskyDev: prysmVcTagTest,
				types.Network_Holesky:    prysmVcTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Prysm Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// The the title for the config
func (cfg *PrysmVcConfig) GetTitle() string {
	return "Prysm Settings"
}

// Get the parameters for this config
func (cfg *PrysmVcConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *PrysmVcConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	lodestarVcTagTest string = "chainsafe/lodestar:v1.12.1"
	lodestarVcTagProd string = "chainsafe/lodestar:v1.12.1"
)

// Configuration for the Lodestar VC
type LodestarVcConfig struct {
	Title string

	// The Docker Hub tag for Lodestar VC
	ContainerTag types.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags types.Parameter[string]
}

// Generates a new Lodestar VC configuration
func NewLodestarVcConfig(cfg *HyperdriveConfig) *LodestarVcConfig {
	return &LodestarVcConfig{
		Title: "Lodestar Settings",

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Lodestar container from Docker Hub you want to use for the Validator Client.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    lodestarVcTagProd,
				types.Network_HoleskyDev: lodestarVcTagTest,
				types.Network_Holesky:    lodestarVcTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Lodestar's Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// Get the parameters for this config
func (cfg *LodestarVcConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// The the title for the config
func (cfg *LodestarVcConfig) GetConfigTitle() string {
	return cfg.Title
}

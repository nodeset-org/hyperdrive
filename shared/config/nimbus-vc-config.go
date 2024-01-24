package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Tags
	nimbusVcTagTest string = "statusim/nimbus-validator-client:multiarch-v24.1.1"
	nimbusVcTagProd string = "statusim/nimbus-validator-client:multiarch-v24.1.1"
)

// Configuration for Nimbus
type NimbusVcConfig struct {
	Title string

	// The Docker Hub tag for the VC
	ContainerTag types.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags types.Parameter[string]
}

// Generates a new Nimbus VC configuration
func NewNimbusVcConfig(cfg *HyperdriveConfig) *NimbusVcConfig {
	return &NimbusVcConfig{
		Title: "Nimbus Settings",

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Nimbus Validator Client container you want to use on Docker Hub.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    nimbusVcTagProd,
				types.Network_HoleskyDev: nimbusVcTagTest,
				types.Network_Holesky:    nimbusVcTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Nimbus's Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *NimbusVcConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the title for the config
func (cfg *NimbusVcConfig) GetTitle() string {
	return cfg.Title
}

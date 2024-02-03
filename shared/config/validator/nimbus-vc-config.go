package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Tags
	nimbusVcTagTest string = "statusim/nimbus-validator-client:multiarch-v24.1.1"
	nimbusVcTagProd string = "statusim/nimbus-validator-client:multiarch-v24.1.1"
)

// Configuration for Nimbus
type NimbusVcConfig struct {
	// The Docker Hub tag for the VC
	ContainerTag types.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags types.Parameter[string]
}

// Generates a new Nimbus VC configuration
func NewNimbusVcConfig() *NimbusVcConfig {
	return &NimbusVcConfig{
		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
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
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Nimbus Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *NimbusVcConfig) GetTitle() string {
	return "Nimbus Validator Client"
}

// Get the parameters for this config
func (cfg *NimbusVcConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *NimbusVcConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

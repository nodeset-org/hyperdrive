package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

const (
	// Tags
	nimbusVcTagTest string = "statusim/nimbus-validator-client:multiarch-v24.2.0"
	nimbusVcTagProd string = "statusim/nimbus-validator-client:multiarch-v24.2.0"
)

// Configuration for Nimbus
type NimbusVcConfig struct {
	// The Docker Hub tag for the VC
	ContainerTag config.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags config.Parameter[string]
}

// Generates a new Nimbus VC configuration
func NewNimbusVcConfig() *NimbusVcConfig {
	return &NimbusVcConfig{
		ContainerTag: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
				Description:        "The tag name of the Nimbus Validator Client container you want to use on Docker Hub.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[config.Network]string{
				config.Network_Mainnet:    nimbusVcTagProd,
				config.Network_HoleskyDev: nimbusVcTagTest,
				config.Network_Holesky:    nimbusVcTagTest,
			},
		},

		AdditionalFlags: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Nimbus Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_ValidatorClients},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: "",
			},
		},
	}
}

// Get the title for the config
func (cfg *NimbusVcConfig) GetTitle() string {
	return "Nimbus Validator Client"
}

// Get the parameters for this config
func (cfg *NimbusVcConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *NimbusVcConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{}
}

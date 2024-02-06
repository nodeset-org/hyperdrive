package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

const (
	// Tags
	prysmVcTagTest string = "rocketpool/prysm:v4.2.1"
	prysmVcTagProd string = "rocketpool/prysm:v4.2.1"
)

// Configuration for the Prysm VC
type PrysmVcConfig struct {
	// The Docker Hub tag for the Prysm BN
	ContainerTag config.Parameter[string]

	// Custom command line flags for the BN
	AdditionalFlags config.Parameter[string]
}

// Generates a new Prysm VC configuration
func NewPrysmVcConfig() *PrysmVcConfig {
	return &PrysmVcConfig{
		ContainerTag: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
				Description:        "The tag name of the Prysm container on Docker Hub you want to use for the Validator Client.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[config.Network]string{
				config.Network_Mainnet:    prysmVcTagProd,
				config.Network_HoleskyDev: prysmVcTagTest,
				config.Network_Holesky:    prysmVcTagTest,
			},
		},

		AdditionalFlags: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Prysm Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// The title for the config
func (cfg *PrysmVcConfig) GetTitle() string {
	return "Prysm Validator Client"
}

// Get the parameters for this config
func (cfg *PrysmVcConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *PrysmVcConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{}
}

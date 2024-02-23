package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

const (
	lodestarVcTagTest string = "chainsafe/lodestar:v1.15.0"
	lodestarVcTagProd string = "chainsafe/lodestar:v1.15.0"
)

// Configuration for the Lodestar VC
type LodestarVcConfig struct {
	// The Docker Hub tag for Lodestar VC
	ContainerTag config.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags config.Parameter[string]
}

// Generates a new Lodestar VC configuration
func NewLodestarVcConfig() *LodestarVcConfig {
	return &LodestarVcConfig{
		ContainerTag: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
				Description:        "The tag name of the Lodestar container from Docker Hub you want to use for the Validator Client.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[config.Network]string{
				config.Network_Mainnet:    lodestarVcTagProd,
				config.Network_HoleskyDev: lodestarVcTagTest,
				config.Network_Holesky:    lodestarVcTagTest,
			},
		},

		AdditionalFlags: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Lodestar Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *LodestarVcConfig) GetTitle() string {
	return "Lodestar Validator Client"
}

// Get the parameters for this config
func (cfg *LodestarVcConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *LodestarVcConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{}
}

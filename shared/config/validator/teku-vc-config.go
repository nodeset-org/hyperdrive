package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

const (
	// Tags
	tekuVcTagTest string = "consensys/teku:24.1.1"
	tekuVcTagProd string = "consensys/teku:24.1.1"
)

// Configuration for Teku
type TekuVcConfig struct {
	// The Docker Hub tag for the Teku VC
	ContainerTag config.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags config.Parameter[string]
}

// Generates a new Teku VC configuration
func NewTekuVcConfig() *TekuVcConfig {
	return &TekuVcConfig{
		ContainerTag: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
				Description:        "The tag name of the Teku container on Docker Hub you want to use for the Validator Client.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[config.Network]string{
				config.Network_Mainnet:    tekuVcTagProd,
				config.Network_HoleskyDev: tekuVcTagTest,
				config.Network_Holesky:    tekuVcTagTest,
			},
		},

		AdditionalFlags: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Teku Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *TekuVcConfig) GetTitle() string {
	return "Teku Validator Client"
}

// Get the parameters for this config
func (cfg *TekuVcConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *TekuVcConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{}
}

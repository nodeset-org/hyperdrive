package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Common configuration for all validator clients
type ValidatorClientCommonConfig struct {
	// Custom proposal graffiti
	Graffiti types.Parameter[string]

	// Toggle for enabling doppelganger detection
	DoppelgangerDetection types.Parameter[bool]

	// The port to expose VC metrics on
	MetricsPort types.Parameter[uint16]
}

// Generates a new common VC configuration
func NewValidatorClientCommonConfig() *ValidatorClientCommonConfig {
	return &ValidatorClientCommonConfig{
		Graffiti: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:          ids.GraffitiID,
				Name:        "Custom Graffiti",
				Description: "Add a short message to any blocks you propose, so the world can see what you have to say!\nIt has a 16 character limit.",
				MaxLength:   16,
				AffectsContainers: []types.ContainerID{
					types.ContainerID_ValidatorClients,
				},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		DoppelgangerDetection: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.DoppelgangerDetectionID,
				Name:               "Enable Doppelgänger Detection",
				Description:        "If enabled, your client will *intentionally* miss 1 or 2 attestations on startup to check if validator keys are already running elsewhere. If they are, it will disable validation duties for them to prevent you from being slashed.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: true,
			},
		},

		MetricsPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.MetricsPortID,
				Name:               "Validator Client Metrics Port",
				Description:        "The port your Validator Client should expose its metrics on, if metrics collection is enabled.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients, types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 9101,
			},
		},
	}
}

// Get the title for the config
func (cfg *ValidatorClientCommonConfig) GetTitle() string {
	return "Common Validator Client Settings"
}

// Get the parameters for this config
func (cfg *ValidatorClientCommonConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Graffiti,
		&cfg.DoppelgangerDetection,
		&cfg.MetricsPort,
	}
}

// Get the sections underneath this one
func (cfg *ValidatorClientCommonConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

// Common configuration for all validator clients
type ValidatorClientCommonConfig struct {
	// Custom proposal graffiti
	Graffiti config.Parameter[string]

	// Toggle for enabling doppelganger detection
	DoppelgangerDetection config.Parameter[bool]

	// The port to expose VC metrics on
	MetricsPort config.Parameter[uint16]
}

// Generates a new common VC configuration
func NewValidatorClientCommonConfig() *ValidatorClientCommonConfig {
	return &ValidatorClientCommonConfig{
		Graffiti: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:          ids.GraffitiID,
				Name:        "Custom Graffiti",
				Description: "Add a short message to any blocks you propose, so the world can see what you have to say!\nIt has a 16 character limit.",
				MaxLength:   16,
				AffectsContainers: []config.ContainerID{
					config.ContainerID_ValidatorClients,
				},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: "",
			},
		},

		DoppelgangerDetection: config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.DoppelgangerDetectionID,
				Name:               "Enable Doppelgänger Detection",
				Description:        "If enabled, your client will *intentionally* miss 1 or 2 attestations on startup to check if validator keys are already running elsewhere. If they are, it will disable validation duties for them to prevent you from being slashed.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]bool{
				config.Network_All: true,
			},
		},

		MetricsPort: config.Parameter[uint16]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.MetricsPortID,
				Name:               "Validator Client Metrics Port",
				Description:        "The port your Validator Client should expose its metrics on, if metrics collection is enabled.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_ValidatorClients, config.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]uint16{
				config.Network_All: 9101,
			},
		},
	}
}

// Get the title for the config
func (cfg *ValidatorClientCommonConfig) GetTitle() string {
	return "Common Validator Client"
}

// Get the parameters for this config
func (cfg *ValidatorClientCommonConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.Graffiti,
		&cfg.DoppelgangerDetection,
		&cfg.MetricsPort,
	}
}

// Get the sections underneath this one
func (cfg *ValidatorClientCommonConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{}
}

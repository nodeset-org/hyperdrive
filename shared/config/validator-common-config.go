package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	VcCommonGraffitiID              string = "graffiti"
	VcCommonDoppelgangerDetectionID string = "doppelgangerDetection"
)

// Common parameters shared by each Validator Client that Hyperdrive owns
type ValidatorCommonConfig struct {
	Title string

	// Custom proposal graffiti
	Graffiti types.Parameter[string]

	// Toggle for enabling doppelganger detection
	DoppelgangerDetection types.Parameter[bool]
}

// Create a new ValidatorCommonParams struct
func NewValidatorCommonConfig(cfg *HyperdriveConfig) *ValidatorCommonConfig {
	return &ValidatorCommonConfig{
		Title: "Common Validator Client Settings",

		Graffiti: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 VcCommonGraffitiID,
				Name:               "Custom Graffiti",
				Description:        "Add a short message to any blocks you propose, so the world can see what you have to say!\nIt has a 16 character limit.",
				MaxLength:          16,
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		DoppelgangerDetection: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 VcCommonDoppelgangerDetectionID,
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
	}
}

// Get the parameters for this config
func (cfg *ValidatorCommonConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Graffiti,
		&cfg.DoppelgangerDetection,
	}
}

// The the title for the config
func (cfg *ValidatorCommonConfig) GetConfigTitle() string {
	return cfg.Title
}

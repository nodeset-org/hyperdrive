package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/utils/sys"
)

const (
	// Tags
	lighthouseVcTagPortableTest string = "sigp/lighthouse:v4.5.0"
	lighthouseVcTagPortableProd string = "sigp/lighthouse:v4.5.0"
	lighthouseVcTagModernTest   string = "sigp/lighthouse:v4.5.0-modern"
	lighthouseVcTagModernProd   string = "sigp/lighthouse:v4.5.0-modern"
)

// Configuration for the Lighthouse VC
type LighthouseVcConfig struct {
	Title string

	// The Docker Hub tag for Lighthouse VC
	ContainerTag types.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags types.Parameter[string]
}

// Generates a new Lighthouse VC configuration
func NewLighthouseVcConfig(cfg *HyperdriveConfig) *LighthouseVcConfig {
	return &LighthouseVcConfig{
		Title: "Lighthouse Settings",

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Lighthouse container from Docker Hub you want to use for the Validator Client.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    getLighthouseVcTagProd(),
				types.Network_HoleskyDev: getLighthouseVcTagTest(),
				types.Network_Holesky:    getLighthouseVcTagTest(),
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Lighthouse's Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *LighthouseVcConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// The the title for the config
func (cfg *LighthouseVcConfig) GetConfigTitle() string {
	return cfg.Title
}

// Get the appropriate LH default tag for production
func getLighthouseVcTagProd() string {
	missingFeatures := sys.GetMissingModernCpuFeatures()
	if len(missingFeatures) > 0 {
		return lighthouseVcTagPortableProd
	}
	return lighthouseVcTagModernProd
}

// Get the appropriate LH default tag for testnets
func getLighthouseVcTagTest() string {
	missingFeatures := sys.GetMissingModernCpuFeatures()
	if len(missingFeatures) > 0 {
		return lighthouseVcTagPortableTest
	}
	return lighthouseVcTagModernTest
}

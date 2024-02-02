package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
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
	// The Docker Hub tag for Lighthouse VC
	ContainerTag types.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags types.Parameter[string]
}

// Generates a new Lighthouse VC configuration
func NewLighthouseVcConfig() *LighthouseVcConfig {
	return &LighthouseVcConfig{
		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
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
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Lighthouse Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// The title for the config
func (cfg *LighthouseVcConfig) GetTitle() string {
	return "Lighthouse Settings"
}

// Get the parameters for this config
func (cfg *LighthouseVcConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *LighthouseVcConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
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

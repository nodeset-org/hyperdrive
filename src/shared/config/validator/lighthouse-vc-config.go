package validator

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/utils/sys"
)

const (
	// Tags
	lighthouseVcTagPortableTest string = "sigp/lighthouse:v4.6.0"
	lighthouseVcTagPortableProd string = "sigp/lighthouse:v4.6.0"
	lighthouseVcTagModernTest   string = "sigp/lighthouse:v4.6.0-modern"
	lighthouseVcTagModernProd   string = "sigp/lighthouse:v4.6.0-modern"
)

// Configuration for the Lighthouse VC
type LighthouseVcConfig struct {
	// The Docker Hub tag for Lighthouse VC
	ContainerTag config.Parameter[string]

	// Custom command line flags for the VC
	AdditionalFlags config.Parameter[string]
}

// Generates a new Lighthouse VC configuration
func NewLighthouseVcConfig() *LighthouseVcConfig {
	return &LighthouseVcConfig{
		ContainerTag: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Validator Client Container Tag",
				Description:        "The tag name of the Lighthouse container from Docker Hub you want to use for the Validator Client.",
				AffectsContainers:  []config.ContainerID{config.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[config.Network]string{
				config.Network_Mainnet:    getLighthouseVcTagProd(),
				config.Network_HoleskyDev: getLighthouseVcTagTest(),
				config.Network_Holesky:    getLighthouseVcTagTest(),
			},
		},

		AdditionalFlags: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass the Lighthouse Validator Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *LighthouseVcConfig) GetTitle() string {
	return "Lighthouse Validator Client"
}

// Get the parameters for this config
func (cfg *LighthouseVcConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *LighthouseVcConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{}
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

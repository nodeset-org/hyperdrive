package stakewise

import (
	"github.com/nodeset-org/hyperdrive/shared/config/validator"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	StakewiseEnableID      string = "enable"
	OperatorContainerTagID string = "operatorContainerTag"
	AdditionalOpFlagsID    string = "additionalOpFlags"

	// Tags
	operatorTag string = "europe-west4-docker.pkg.dev/stakewiselabs/public/v3-operator:v1.0.8"
)

// Configuration for Stakewise
type StakewiseConfig struct {
	// Toggle for enabling access to the root filesystem (for multiple disk usage metrics)
	Enable types.Parameter[bool]

	// The Docker Hub tag for the Stakewise operator
	OperatorContainerTag types.Parameter[string]

	// Custom command line flags
	AdditionalOpFlags types.Parameter[string]

	// Validator client configs
	Lighthouse *validator.LighthouseVcConfig
	Lodestar   *validator.LodestarVcConfig
	Nimbus     *validator.NimbusVcConfig
	Prysm      *validator.PrysmVcConfig
	Teku       *validator.TekuVcConfig
}

// Generates a new Stakewise config
func NewStakewiseConfig() *StakewiseConfig {
	cfg := &StakewiseConfig{
		Enable: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 StakewiseEnableID,
				Name:               "Enable",
				Description:        "Enable support for Stakewise <placeholder description>",
				AffectsContainers:  []types.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: false,
			},
		},

		OperatorContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 OperatorContainerTagID,
				Name:               "Operator Container Tag",
				Description:        "The tag name of the Stakewise Operator image to use. See https://github.com/stakewise/v3-operator#using-docker for more details.",
				AffectsContainers:  []types.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_All: operatorTag,
			},
		},

		AdditionalOpFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalOpFlagsID,
				Name:               "Additional Operator Flags",
				Description:        "Additional custom command line flags you want to pass to the Operator container, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []types.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}

	cfg.Lighthouse = validator.NewLighthouseVcConfig()
	cfg.Lodestar = validator.NewLodestarVcConfig()
	cfg.Nimbus = validator.NewNimbusVcConfig()
	cfg.Prysm = validator.NewPrysmVcConfig()
	cfg.Teku = validator.NewTekuVcConfig()

	return cfg
}

// The the title for the config
func (cfg *StakewiseConfig) GetTitle() string {
	return "Stakewise Settings"
}

// Get the parameters for this config
func (cfg *StakewiseConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Enable,
		&cfg.OperatorContainerTag,
		&cfg.AdditionalOpFlags,
	}
}

// Get the sections underneath this one
func (cfg *StakewiseConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{
		"lighthouse": cfg.Lighthouse,
		"lodestar":   cfg.Lodestar,
		"nimbus":     cfg.Nimbus,
		"prysm":      cfg.Prysm,
		"teku":       cfg.Teku,
	}
}

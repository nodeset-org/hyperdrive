package swconfig

import (
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/validator"
)

const (
	// Param IDs
	StakewiseEnableID      string = "enable"
	OperatorContainerTagID string = "operatorContainerTag"
	AdditionalOpFlagsID    string = "additionalOpFlags"

	// Tags
	daemonTag   string = "nodeset/hyperdrive-stakewise:v" + shared.HyperdriveVersion
	operatorTag string = "europe-west4-docker.pkg.dev/stakewiselabs/public/v3-operator:v1.0.8"
)

// Configuration for Stakewise
type StakewiseConfig struct {
	hdCfg *config.HyperdriveConfig

	// Toggle for enabling access to the root filesystem (for multiple disk usage metrics)
	Enabled config.Parameter[bool]

	// The Docker Hub tag for the Stakewise operator
	OperatorContainerTag config.Parameter[string]

	// Custom command line flags
	AdditionalOpFlags config.Parameter[string]

	// Validator client configs
	VcCommon   *validator.ValidatorClientCommonConfig
	Lighthouse *validator.LighthouseVcConfig
	Lodestar   *validator.LodestarVcConfig
	Nimbus     *validator.NimbusVcConfig
	Prysm      *validator.PrysmVcConfig
	Teku       *validator.TekuVcConfig
}

// Generates a new Stakewise config
func NewStakewiseConfig(hdCfg *config.HyperdriveConfig) *StakewiseConfig {
	cfg := &StakewiseConfig{
		hdCfg: hdCfg,

		Enabled: config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 StakewiseEnableID,
				Name:               "Enable",
				Description:        "Enable support for Stakewise <placeholder description>",
				AffectsContainers:  []config.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]bool{
				config.Network_All: false,
			},
		},

		OperatorContainerTag: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 OperatorContainerTagID,
				Name:               "Operator Container Tag",
				Description:        "The tag name of the Stakewise Operator image to use. See https://github.com/stakewise/v3-operator#using-docker for more details.",
				AffectsContainers:  []config.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[config.Network]string{
				config.Network_All: operatorTag,
			},
		},

		AdditionalOpFlags: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 AdditionalOpFlagsID,
				Name:               "Additional Operator Flags",
				Description:        "Additional custom command line flags you want to pass to the Operator container, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []config.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: "",
			},
		},
	}

	cfg.VcCommon = validator.NewValidatorClientCommonConfig()
	cfg.Lighthouse = validator.NewLighthouseVcConfig()
	cfg.Lodestar = validator.NewLodestarVcConfig()
	cfg.Nimbus = validator.NewNimbusVcConfig()
	cfg.Prysm = validator.NewPrysmVcConfig()
	cfg.Teku = validator.NewTekuVcConfig()

	return cfg
}

// The title for the config
func (cfg *StakewiseConfig) GetTitle() string {
	return "Stakewise"
}

// Get the parameters for this config
func (cfg *StakewiseConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.Enabled,
		&cfg.OperatorContainerTag,
		&cfg.AdditionalOpFlags,
	}
}

// Get the sections underneath this one
func (cfg *StakewiseConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{
		"common":     cfg.VcCommon,
		"lighthouse": cfg.Lighthouse,
		"lodestar":   cfg.Lodestar,
		"nimbus":     cfg.Nimbus,
		"prysm":      cfg.Prysm,
		"teku":       cfg.Teku,
	}
}

// ===================
// === Module Info ===
// ===================

// The module name
func (cfg *StakewiseConfig) GetModuleName() string {
	return ModuleName
}

func (cfg *StakewiseConfig) GetValidatorContainerTagInfo() map[config.ContainerID]string {
	return map[config.ContainerID]string{
		ContainerID_StakewiseValidator: cfg.GetVcContainerTag(),
	}
}

func (cfg *StakewiseConfig) GetContainersToDeploy() []config.ContainerID {
	return []config.ContainerID{
		ContainerID_StakewiseDaemon,
		ContainerID_StakewiseOperator,
		ContainerID_StakewiseValidator,
	}
}

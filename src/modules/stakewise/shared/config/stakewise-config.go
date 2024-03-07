package swconfig

import (
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	nmc_config "github.com/rocket-pool/node-manager-core/config"
)

const (
	// Param IDs
	StakewiseEnableID      string = "enable"
	OperatorContainerTagID string = "operatorContainerTag"
	AdditionalOpFlagsID    string = "additionalOpFlags"
	VerifyDepositRootsID   string = "verifyDepositRoots"

	// Tags
	daemonTag   string = "nodeset/hyperdrive-stakewise:v" + shared.HyperdriveVersion
	operatorTag string = "europe-west4-docker.pkg.dev/stakewiselabs/public/v3-operator:v1.0.8"
)

// Configuration for Stakewise
type StakewiseConfig struct {
	hdCfg *config.HyperdriveConfig

	// Toggle for enabling access to the root filesystem (for multiple disk usage metrics)
	Enabled nmc_config.Parameter[bool]

	// Toggle for verifying deposit data Merkle roots before saving
	VerifyDepositsRoot nmc_config.Parameter[bool]

	// The Docker Hub tag for the Stakewise operator
	OperatorContainerTag nmc_config.Parameter[string]

	// Custom command line flags
	AdditionalOpFlags nmc_config.Parameter[string]

	// Validator client configs
	VcCommon   *nmc_config.ValidatorClientCommonConfig
	Lighthouse *nmc_config.LighthouseVcConfig
	Lodestar   *nmc_config.LodestarVcConfig
	Nimbus     *nmc_config.NimbusVcConfig
	Prysm      *nmc_config.PrysmVcConfig
	Teku       *nmc_config.TekuVcConfig
}

// Generates a new Stakewise config
func NewStakewiseConfig(hdCfg *config.HyperdriveConfig) *StakewiseConfig {
	cfg := &StakewiseConfig{
		hdCfg: hdCfg,

		Enabled: nmc_config.Parameter[bool]{
			ParameterCommon: &nmc_config.ParameterCommon{
				ID:                 StakewiseEnableID,
				Name:               "Enable",
				Description:        "Enable support for Stakewise (see more at https://docs.nodeset.io).",
				AffectsContainers:  []nmc_config.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[nmc_config.Network]bool{
				nmc_config.Network_All: false,
			},
		},

		VerifyDepositsRoot: nmc_config.Parameter[bool]{
			ParameterCommon: &nmc_config.ParameterCommon{
				ID:                 VerifyDepositRootsID,
				Name:               "Verify Deposits Root",
				Description:        "Enable this to verify that the Merkle root of aggregated deposit data returned by the NodeSet server matches the Merkle root stored in the NodeSet vault contract. This is a safety mechanism to ensure the Stakewise Operator container won't try to submit deposits for validators that the NodeSet vault hasn't verified yet.\n\n[orange]Don't disable this unless you know what you're doing.",
				AffectsContainers:  []nmc_config.ContainerID{ContainerID_StakewiseDaemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[nmc_config.Network]bool{
				nmc_config.Network_All: true,
			},
		},

		OperatorContainerTag: nmc_config.Parameter[string]{
			ParameterCommon: &nmc_config.ParameterCommon{
				ID:                 OperatorContainerTagID,
				Name:               "Operator Container Tag",
				Description:        "The tag name of the Stakewise Operator image to use. See https://github.com/stakewise/v3-operator#using-docker for more details.",
				AffectsContainers:  []nmc_config.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[nmc_config.Network]string{
				nmc_config.Network_All: operatorTag,
			},
		},

		AdditionalOpFlags: nmc_config.Parameter[string]{
			ParameterCommon: &nmc_config.ParameterCommon{
				ID:                 AdditionalOpFlagsID,
				Name:               "Additional Operator Flags",
				Description:        "Additional custom command line flags you want to pass to the Operator container, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []nmc_config.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[nmc_config.Network]string{
				nmc_config.Network_All: "",
			},
		},
	}

	cfg.VcCommon = nmc_config.NewValidatorClientCommonConfig()
	cfg.Lighthouse = nmc_config.NewLighthouseVcConfig()
	cfg.Lodestar = nmc_config.NewLodestarVcConfig()
	cfg.Nimbus = nmc_config.NewNimbusVcConfig()
	cfg.Prysm = nmc_config.NewPrysmVcConfig()
	cfg.Teku = nmc_config.NewTekuVcConfig()
	cfg.Lighthouse.ContainerTag.Default[config.Network_HoleskyDev] = cfg.Lighthouse.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Lodestar.ContainerTag.Default[config.Network_HoleskyDev] = cfg.Lodestar.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Nimbus.ContainerTag.Default[config.Network_HoleskyDev] = cfg.Nimbus.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Prysm.ContainerTag.Default[config.Network_HoleskyDev] = cfg.Prysm.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Teku.ContainerTag.Default[config.Network_HoleskyDev] = cfg.Teku.ContainerTag.Default[nmc_config.Network_Holesky]

	return cfg
}

// The title for the config
func (cfg *StakewiseConfig) GetTitle() string {
	return "Stakewise"
}

// Get the parameters for this config
func (cfg *StakewiseConfig) GetParameters() []nmc_config.IParameter {
	return []nmc_config.IParameter{
		&cfg.Enabled,
		&cfg.VerifyDepositsRoot,
		&cfg.OperatorContainerTag,
		&cfg.AdditionalOpFlags,
	}
}

// Get the sections underneath this one
func (cfg *StakewiseConfig) GetSubconfigs() map[string]nmc_config.IConfigSection {
	return map[string]nmc_config.IConfigSection{
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

func (cfg *StakewiseConfig) GetValidatorContainerTagInfo() map[nmc_config.ContainerID]string {
	return map[nmc_config.ContainerID]string{
		ContainerID_StakewiseValidator: cfg.GetVcContainerTag(),
	}
}

func (cfg *StakewiseConfig) GetContainersToDeploy() []nmc_config.ContainerID {
	return []nmc_config.ContainerID{
		ContainerID_StakewiseDaemon,
		ContainerID_StakewiseOperator,
		ContainerID_StakewiseValidator,
	}
}

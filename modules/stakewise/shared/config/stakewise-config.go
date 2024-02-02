package swconfig

import (
	"fmt"

	swshared "github.com/nodeset-org/hyperdrive/modules/stakewise/shared"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/validator"
	"github.com/nodeset-org/hyperdrive/shared/types"
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
	Enabled types.Parameter[bool]

	// The Docker Hub tag for the Stakewise operator
	OperatorContainerTag types.Parameter[string]

	// Custom command line flags
	AdditionalOpFlags types.Parameter[string]

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

		Enabled: types.Parameter[bool]{
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
func (cfg *StakewiseConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Enabled,
		&cfg.OperatorContainerTag,
		&cfg.AdditionalOpFlags,
	}
}

// Get the sections underneath this one
func (cfg *StakewiseConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{
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
	return "stakewise"
}

func (cfg *StakewiseConfig) GetValidatorContainerTagInfo() map[string]string {
	return map[string]string{
		string(ContainerID_StakewiseValidator): cfg.GetStakewiseVcContainerTag(),
	}
}

// ==================
// === Templating ===
// ==================

// The tag for the daemon container
func (cfg *StakewiseConfig) DaemonTag() string {
	return daemonTag
}

// Get the container tag of the selected VC
func (cfg *StakewiseConfig) GetStakewiseVcContainerTag() string {
	bn := cfg.hdCfg.GetSelectedBeaconNode()
	switch bn {
	case types.BeaconNode_Lighthouse:
		return cfg.Lighthouse.ContainerTag.Value
	case types.BeaconNode_Lodestar:
		return cfg.Lodestar.ContainerTag.Value
	case types.BeaconNode_Nimbus:
		return cfg.Nimbus.ContainerTag.Value
	case types.BeaconNode_Prysm:
		return cfg.Prysm.ContainerTag.Value
	case types.BeaconNode_Teku:
		return cfg.Teku.ContainerTag.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", bn))
	}
}

// Gets the additional flags of the selected VC
func (cfg *StakewiseConfig) GetStakewiseVcAdditionalFlags() string {
	bn := cfg.hdCfg.GetSelectedBeaconNode()
	switch bn {
	case types.BeaconNode_Lighthouse:
		return cfg.Lighthouse.AdditionalFlags.Value
	case types.BeaconNode_Lodestar:
		return cfg.Lodestar.AdditionalFlags.Value
	case types.BeaconNode_Nimbus:
		return cfg.Nimbus.AdditionalFlags.Value
	case types.BeaconNode_Prysm:
		return cfg.Prysm.AdditionalFlags.Value
	case types.BeaconNode_Teku:
		return cfg.Teku.AdditionalFlags.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", bn))
	}
}

// Check if any of the services have doppelganger detection enabled
// NOTE: update this with each new service that runs a VC!
func (cfg *StakewiseConfig) IsDoppelgangerEnabled() bool {
	return cfg.VcCommon.DoppelgangerDetection.Value
}

// Used by text/template to format validator.yml
func (cfg *StakewiseConfig) StakewiseGraffiti() (string, error) {
	prefix := cfg.hdCfg.GraffitiPrefix()
	customGraffiti := cfg.VcCommon.Graffiti.Value
	if customGraffiti == "" {
		return prefix, nil
	}
	return fmt.Sprintf("%s (%s)", prefix, customGraffiti), nil
}

func (cfg *StakewiseConfig) StakewiseFeeRecipient() string {
	res := swshared.NewStakewiseResources(cfg.hdCfg.Network.Value)
	return res.FeeRecipient.Hex()
}

func (cfg *StakewiseConfig) StakewiseVault() string {
	res := swshared.NewStakewiseResources(cfg.hdCfg.Network.Value)
	return res.Vault.Hex()
}

func (cfg *StakewiseConfig) StakewiseNetwork() string {
	res := swshared.NewStakewiseResources(cfg.hdCfg.Network.Value)
	return res.NodesetNetwork
}

func (cfg *StakewiseConfig) IsEnabled() bool {
	return cfg.Enabled.Value
}

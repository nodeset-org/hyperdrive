package constconfig

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
)

const (
	// Param IDs
	// OperatorContainerTagID string = "operatorContainerTag"
	AdditionalOpFlagsID string = "additionalOpFlags"
	// VerifyDepositRootsID   string = "verifyDepositRoots"

	// Tags
	// daemonTag   string = "nodeset/hyperdrive-stakewise:v" + shared.HyperdriveVersion
	// operatorTag string = "europe-west4-docker.pkg.dev/stakewiselabs/public/v3-operator:v1.0.8"
)

// Configuration for Constellation
type ConstellationConfig struct {
	hdCfg *config.HyperdriveConfig

	// Toggle for enabling access to the root filesystem (for multiple disk usage metrics)
	Enabled config.Parameter[bool]

	// Custom command line flags
	AdditionalOpFlags config.Parameter[string]
}

// The title for the config
func (cfg *ConstellationConfig) GetTitle() string {
	return "Constellation"
}

// Generates a new Constellation config
func NewConstellationConfig(hdCfg *config.HyperdriveConfig) *ConstellationConfig {
	cfg := &ConstellationConfig{
		hdCfg: hdCfg,

		Enabled: config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				Name:               "Enable",
				Description:        "Enable support for Stakewise (see more at https://docs.nodeset.io).",
				AffectsContainers:  []config.ContainerID{ContainerID_ConstellationDaemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]bool{
				config.Network_All: false,
			},
		},

		AdditionalOpFlags: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 AdditionalOpFlagsID,
				Name:               "Additional Operator Flags",
				Description:        "Additional custom command line flags you want to pass to the Operator container, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []config.ContainerID{ContainerID_ConstellationDaemon},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]string{
				config.Network_All: "",
			},
		},
	}
	return cfg
}

// The module name
func (cfg *ConstellationConfig) GetModuleName() string {
	return ModuleName
}

func (cfg *ConstellationConfig) GetValidatorContainerTagInfo() map[config.ContainerID]string {
	return map[config.ContainerID]string{
		ContainerID_ConstellationValidator: cfg.GetVcContainerTag(),
	}
}

func (cfg *ConstellationConfig) GetContainersToDeploy() []config.ContainerID {
	return []config.ContainerID{
		ContainerID_ConstellationDaemon,
		// ContainerID_ConstellationOperator,
		ContainerID_ConstellationValidator,
	}
}

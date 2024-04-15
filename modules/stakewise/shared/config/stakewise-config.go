package swconfig

import (
	"github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	hdids "github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	// Tags
	daemonTag   string = "nodeset/hyperdrive-stakewise:v" + shared.HyperdriveVersion
	operatorTag string = "europe-west4-docker.pkg.dev/stakewiselabs/public/v3-operator:v1.1.0"
)

// Configuration for Stakewise
type StakewiseConfig struct {
	// Toggle for enabling access to the root filesystem (for multiple disk usage metrics)
	Enabled config.Parameter[bool]

	// Toggle for verifying deposit data Merkle roots before saving
	VerifyDepositsRoot config.Parameter[bool]

	// The Docker Hub tag for the Stakewise operator
	OperatorContainerTag config.Parameter[string]

	// Custom command line flags
	AdditionalOpFlags config.Parameter[string]

	// Validator client configs
	VcCommon   *config.ValidatorClientCommonConfig
	Lighthouse *config.LighthouseVcConfig
	Lodestar   *config.LodestarVcConfig
	Nimbus     *config.NimbusVcConfig
	Prysm      *config.PrysmVcConfig
	Teku       *config.TekuVcConfig

	// Internal fields
	Version   string
	hdCfg     *hdconfig.HyperdriveConfig
	resources *StakewiseResources
}

// Generates a new Stakewise config
func NewStakewiseConfig(hdCfg *hdconfig.HyperdriveConfig) *StakewiseConfig {
	cfg := &StakewiseConfig{
		hdCfg: hdCfg,

		Enabled: config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.StakewiseEnableID,
				Name:               "Enable",
				Description:        "Enable support for Stakewise (see more at https://docs.nodeset.io).",
				AffectsContainers:  []config.ContainerID{ContainerID_StakewiseOperator},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]bool{
				config.Network_All: false,
			},
		},

		VerifyDepositsRoot: config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.VerifyDepositRootsID,
				Name:               "Verify Deposits Root",
				Description:        "Enable this to verify that the Merkle root of aggregated deposit data returned by the NodeSet server matches the Merkle root stored in the NodeSet vault contract. This is a safety mechanism to ensure the Stakewise Operator container won't try to submit deposits for validators that the NodeSet vault hasn't verified yet.\n\n[orange]Don't disable this unless you know what you're doing.",
				AffectsContainers:  []config.ContainerID{ContainerID_StakewiseDaemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]bool{
				config.Network_All: true,
			},
		},

		OperatorContainerTag: config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.OperatorContainerTagID,
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
				ID:                 ids.AdditionalOpFlagsID,
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

	cfg.VcCommon = config.NewValidatorClientCommonConfig()
	cfg.Lighthouse = config.NewLighthouseVcConfig()
	cfg.Lodestar = config.NewLodestarVcConfig()
	cfg.Nimbus = config.NewNimbusVcConfig()
	cfg.Prysm = config.NewPrysmVcConfig()
	cfg.Teku = config.NewTekuVcConfig()
	cfg.Lighthouse.ContainerTag.Default[hdconfig.Network_HoleskyDev] = cfg.Lighthouse.ContainerTag.Default[config.Network_Holesky]
	cfg.Lodestar.ContainerTag.Default[hdconfig.Network_HoleskyDev] = cfg.Lodestar.ContainerTag.Default[config.Network_Holesky]
	cfg.Nimbus.ContainerTag.Default[hdconfig.Network_HoleskyDev] = cfg.Nimbus.ContainerTag.Default[config.Network_Holesky]
	cfg.Prysm.ContainerTag.Default[hdconfig.Network_HoleskyDev] = cfg.Prysm.ContainerTag.Default[config.Network_Holesky]
	cfg.Teku.ContainerTag.Default[hdconfig.Network_HoleskyDev] = cfg.Teku.ContainerTag.Default[config.Network_Holesky]

	// Apply the default values for mainnet
	config.ApplyDefaults(cfg, hdCfg.Network.Value)
	cfg.updateResources()

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
		&cfg.VerifyDepositsRoot,
		&cfg.OperatorContainerTag,
		&cfg.AdditionalOpFlags,
	}
}

// Get the sections underneath this one
func (cfg *StakewiseConfig) GetSubconfigs() map[string]config.IConfigSection {
	return map[string]config.IConfigSection{
		ids.VcCommonID:   cfg.VcCommon,
		ids.LighthouseID: cfg.Lighthouse,
		ids.LodestarID:   cfg.Lodestar,
		ids.NimbusID:     cfg.Nimbus,
		ids.PrysmID:      cfg.Prysm,
		ids.TekuID:       cfg.Teku,
	}
}

// Changes the current network, propagating new parameter settings if they are affected
func (cfg *StakewiseConfig) ChangeNetwork(oldNetwork config.Network, newNetwork config.Network) {
	// Run the changes
	config.ChangeNetwork(cfg, oldNetwork, newNetwork)
	cfg.updateResources()
}

// Creates a copy of the configuration
func (cfg *StakewiseConfig) Clone() hdconfig.IModuleConfig {
	clone := NewStakewiseConfig(cfg.hdCfg)
	config.Clone(cfg, clone, cfg.hdCfg.Network.Value)
	clone.Version = cfg.Version
	clone.updateResources()
	return clone
}

// Get the Stakewise resources for the selected network
func (cfg *StakewiseConfig) GetStakewiseResources() *StakewiseResources {
	return cfg.resources
}

// Updates the default parameters based on the current network value
func (cfg *StakewiseConfig) UpdateDefaults(network config.Network) {
	config.UpdateDefaults(cfg, network)
}

// Checks to see if the current configuration is valid; if not, returns a list of errors
func (cfg *StakewiseConfig) Validate() []string {
	errors := []string{}
	return errors
}

// Serialize the module config to a map
func (cfg *StakewiseConfig) Serialize() map[string]any {
	cfgMap := config.Serialize(cfg)
	cfgMap[hdids.VersionID] = cfg.Version
	return cfgMap
}

// Deserialize the module config from a map
func (cfg *StakewiseConfig) Deserialize(configMap map[string]any, network config.Network) error {
	err := config.Deserialize(cfg, configMap, network)
	if err != nil {
		return err
	}
	version, exists := configMap[hdids.VersionID]
	if !exists {
		// Handle pre-version configs
		version = "0.0.1"
	}
	cfg.Version = version.(string)
	return nil
}

// Get the version of the module config
func (cfg *StakewiseConfig) GetVersion() string {
	return cfg.Version
}

// =====================
// === Field Helpers ===
// =====================

// Update the config's resource cache
func (cfg *StakewiseConfig) updateResources() {
	cfg.resources = newStakewiseResources(cfg.hdCfg.Network.Value)
}

// ===================
// === Module Info ===
// ===================

func (cfg *StakewiseConfig) GetHdClientLogFileName() string {
	return ClientLogName
}

func (cfg *StakewiseConfig) GetApiLogFileName() string {
	return hdconfig.ApiLogName
}

func (cfg *StakewiseConfig) GetTasksLogFileName() string {
	return hdconfig.TasksLogName
}

func (cfg *StakewiseConfig) GetLogNames() []string {
	return []string{
		cfg.GetHdClientLogFileName(),
		cfg.GetApiLogFileName(),
		cfg.GetTasksLogFileName(),
	}
}

// The module name
func (cfg *StakewiseConfig) GetModuleName() string {
	return ModuleName
}

// The module name
func (cfg *StakewiseConfig) GetShortName() string {
	return ShortModuleName
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

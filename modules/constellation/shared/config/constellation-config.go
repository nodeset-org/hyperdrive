package constconfig

import (
	"github.com/nodeset-org/hyperdrive/modules/constellation/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	hdids "github.com/nodeset-org/hyperdrive/shared/config/ids"

	"github.com/rocket-pool/node-manager-core/config"
)

const (
	// Tags
	daemonTag string = "nodeset/hyperdrive-constellation:v" + shared.HyperdriveVersion
)

// Configuration for Constellation
type ConstellationConfig struct {
	hdCfg *hdconfig.HyperdriveConfig

	// Toggle for enabling access to the root filesystem (for multiple disk usage metrics)
	Enabled config.Parameter[bool]

	// Validator client configs
	VcCommon   *config.ValidatorClientCommonConfig
	Lighthouse *config.LighthouseVcConfig
	Lodestar   *config.LodestarVcConfig
	Nimbus     *config.NimbusVcConfig
	Prysm      *config.PrysmVcConfig
	Teku       *config.TekuVcConfig

	// Internal fields
	Version   string
	resources *ConstellationResources
}

// Generates a new Constellation config
func NewConstellationConfig(hdCfg *hdconfig.HyperdriveConfig) *ConstellationConfig {
	cfg := &ConstellationConfig{
		hdCfg: hdCfg,

		Enabled: config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				ID:                 ids.ConstellationEnableID,
				Name:               "Enable",
				Description:        "Enable support for Constellation (see more at https://docs.nodeset.io).",
				AffectsContainers:  []config.ContainerID{ContainerID_ConstellationDaemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[config.Network]bool{
				config.Network_All: false,
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
func (cfg *ConstellationConfig) GetTitle() string {
	return "Constellation"
}

// Get the parameters for this config
func (cfg *ConstellationConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.Enabled,
	}
}

// Get the sections underneath this one
func (cfg *ConstellationConfig) GetSubconfigs() map[string]config.IConfigSection {
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
func (cfg *ConstellationConfig) ChangeNetwork(oldNetwork config.Network, newNetwork config.Network) {
	// Run the changes
	config.ChangeNetwork(cfg, oldNetwork, newNetwork)
}

// Creates a copy of the configuration
func (cfg *ConstellationConfig) Clone() hdconfig.IModuleConfig {
	clone := NewConstellationConfig(cfg.hdCfg)
	config.Clone(cfg, clone, cfg.hdCfg.Network.Value)
	return clone
}

// Updates the default parameters based on the current network value
func (cfg *ConstellationConfig) UpdateDefaults(network config.Network) {
	config.UpdateDefaults(cfg, network)
}

// Checks to see if the current configuration is valid; if not, returns a list of errors
func (cfg *ConstellationConfig) Validate() []string {
	errors := []string{}
	return errors
}

// Serialize the module config to a map
func (cfg *ConstellationConfig) Serialize() map[string]any {
	cfgMap := config.Serialize(cfg)
	cfgMap[hdids.VersionID] = cfg.Version
	return cfgMap
}

// Deserialize the module config from a map
func (cfg *ConstellationConfig) Deserialize(configMap map[string]any, network config.Network) error {
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

// ===================
// === Module Info ===
// ===================

func (cfg *ConstellationConfig) GetHdClientLogFileName() string {
	return ClientLogName
}

func (cfg *ConstellationConfig) GetLogNames() []string {
	return []string{
		cfg.GetHdClientLogFileName(),
		cfg.GetApiLogFileName(),
		cfg.GetTasksLogFileName(),
	}
}

func (cfg *ConstellationConfig) GetApiLogFileName() string {
	return hdconfig.ApiLogName
}

func (cfg *ConstellationConfig) GetTasksLogFileName() string {
	return hdconfig.TasksLogName
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
		ContainerID_ConstellationValidator,
	}
}

// The module name
func (cfg *ConstellationConfig) GetShortName() string {
	return ShortModuleName
}

// Get the version of the module config
func (cfg *ConstellationConfig) GetVersion() string {
	return cfg.Version
}

// =====================
// === Field Helpers ===
// =====================

// Update the config's resource cache
func (cfg *ConstellationConfig) updateResources() {
	cfg.resources = NewConstellationResources(cfg.hdCfg.Network.Value)
}

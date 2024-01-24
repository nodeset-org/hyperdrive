/*
Derived from the Rocket Pool Smartnode source code:
https://github.com/rocket-pool/smartnode
*/

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/config/migration"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/pbnjay/memory"

	"gopkg.in/yaml.v2"
)

// =========================
// === Hyperdrive Config ===
// =========================

const (
	// Param IDs
	HyperdriveDebugModeID  string = "debugMode"
	HyperdriveNetworkID    string = "network"
	HyperdriveClientModeID string = "clientMode"
	HyperdriveDirectoryID  string = "hdDir"

	// Tags
	HyperdriveTag string = "nodeset/hyperdrive:v" + shared.HyperdriveVersion
)

// The master configuration struct
type HyperdriveConfig struct {
	// General settings
	DebugMode  types.Parameter[bool]
	Network    types.Parameter[types.Network]
	ClientMode types.Parameter[types.ClientMode]

	// Execution client settings
	LocalExecutionConfig    *LocalExecutionConfig
	ExternalExecutionConfig *ExternalExecutionConfig

	// Beacon node settings
	LocalBeaconConfig    *LocalBeaconConfig
	ExternalBeaconConfig *ExternalBeaconConfig

	// Fallback clients
	Fallback *FallbackConfig

	// Metrics
	Metrics *MetricsConfig

	// Internal fields
	Version             string
	HyperdriveDirectory string
	chainID             map[types.Network]uint
}

// Load configuration settings from a file
func LoadFromFile(path string) (*HyperdriveConfig, error) {
	// Return nil if the file doesn't exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, nil
	}

	// Read the file
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read Hyperdrive settings file at %s: %w", shellescape.Quote(path), err)
	}

	// Attempt to parse it out into a settings map
	var settings map[string]any
	if err := yaml.Unmarshal(configBytes, &settings); err != nil {
		return nil, fmt.Errorf("could not parse settings file: %w", err)
	}

	// Deserialize it into a config object
	cfg := NewHyperdriveConfig(filepath.Dir(path))
	err = cfg.Deserialize(settings)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize settings file: %w", err)
	}

	return cfg, nil
}

// Creates a new Hyperdrive configuration instance
func NewHyperdriveConfig(hdDir string) *HyperdriveConfig {
	cfg := &HyperdriveConfig{
		HyperdriveDirectory: hdDir,

		// Parameters
		DebugMode: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                   HyperdriveDebugModeID,
				Name:                 "Debug Mode",
				Description:          "True to enable debug mode, which at some point will print extra stuff but doesn't right now.",
				AffectsContainers:    []types.ContainerID{types.ContainerID_Daemon},
				EnvironmentVariables: []string{"HD_DEBUG_MODE"},
				CanBeBlank:           false,
				OverwriteOnUpgrade:   false,
			},
			Default: map[types.Network]bool{
				types.Network_All: false,
			},
		},

		Network: types.Parameter[types.Network]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HyperdriveNetworkID,
				Name:               "Network",
				Description:        "The Ethereum network you want to use - select Prater Testnet or Holesky Testnet to practice with fake ETH, or Mainnet to stake on the real network using real ETH.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ExecutionClient, types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getNetworkOptions(),
			Default: map[types.Network]types.Network{
				types.Network_All: types.Network_Mainnet,
			},
		},

		ClientMode: types.Parameter[types.ClientMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HyperdriveClientModeID,
				Name:               "Client Mode",
				Description:        "Choose which mode to use for your Execution and Consensus clients - locally managed (Docker Mode), or externally managed (Hybrid Mode).",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ExecutionClient, types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*types.ParameterOption[types.ClientMode]{
				{
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Locally Managed",
						Description: "Allow the Smartnode to manage the Execution and Consensus clients for you (Docker Mode)",
					},
					Value: types.ClientMode_Local,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Externally Managed",
						Description: "Use existing Execution and Consensus clients that you manage on your own (Hybrid Mode)",
					},
					Value: types.ClientMode_External,
				}},
			Default: map[types.Network]types.ClientMode{
				types.Network_All: types.ClientMode_Local,
			},
		},

		// Internal fields
		chainID: map[types.Network]uint{
			types.Network_Mainnet:    1,     // Mainnet
			types.Network_HoleskyDev: 17000, // Holesky
			types.Network_Holesky:    17000, // Holesky
		},
	}

	// Create the subconfigs
	cfg.Fallback = NewFallbackConfig(cfg)
	cfg.LocalExecutionConfig = NewExecutionCommonConfig(cfg)
	cfg.ExternalExecutionConfig = NewExternalExecutionConfig(cfg)
	cfg.LocalBeaconConfig = NewLocalBeaconConfig(cfg)
	cfg.ExternalBeaconConfig = NewExternalBeaconConfig(cfg)
	cfg.Metrics = NewMetricsConfig(cfg)

	// Apply the default values for mainnet
	cfg.applyAllDefaults()

	return cfg
}

// Get the title for this config
func (cfg *HyperdriveConfig) GetTitle() string {
	return "Hyperdrive Settings"
}

// Get the parameters for this config
func (cfg *HyperdriveConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.DebugMode,
		&cfg.Network,
		&cfg.ClientMode,
	}
}

// Get the subconfigurations for this config
func (cfg *HyperdriveConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{
		"fallback":          cfg.Fallback,
		"localExecution":    cfg.LocalExecutionConfig,
		"externalExecution": cfg.ExternalExecutionConfig,
		"localBeacon":       cfg.LocalBeaconConfig,
		"externalBeacon":    cfg.ExternalBeaconConfig,
		"metrics":           cfg.Metrics,
	}
}

// Serializes the configuration into a map of maps, compatible with a settings file
func (cfg *HyperdriveConfig) Serialize() map[string]any {
	masterMap := map[string]any{}

	hdMap := serialize(cfg)
	masterMap[HyperdriveDirectoryID] = cfg.HyperdriveDirectory
	masterMap[ids.VersionID] = fmt.Sprintf("v%s", shared.HyperdriveVersion)
	masterMap[ids.RootConfigID] = hdMap

	return masterMap
}

// Deserializes a settings file into this config
func (cfg *HyperdriveConfig) Deserialize(masterMap map[string]any) error {
	// Upgrade the config to the latest version
	err := migration.UpdateConfig(masterMap)
	if err != nil {
		return fmt.Errorf("error upgrading configuration to v%s: %w", shared.HyperdriveVersion, err)
	}

	// Get the network
	network := types.Network_Mainnet
	hyperdriveParams, exists := masterMap[ids.RootConfigID]
	if !exists {
		return fmt.Errorf("config is missing the [%s] section", ids.RootConfigID)
	}
	hdMap, isMap := hyperdriveParams.(map[string]any)
	if !isMap {
		return fmt.Errorf("config has an entry named [%s] but it is not a map", ids.RootConfigID)
	}
	networkVal, exists := hdMap[cfg.Network.ID]
	if exists {
		networkString, isString := networkVal.(string)
		if !isString {
			return fmt.Errorf("expected [%s - %s] to be a string but it is not", ids.RootConfigID, cfg.Network.ID)
		}
		network = types.Network(networkString)
	}

	// Deserialize the params and subconfigs
	err = deserialize(cfg, hdMap, network)
	if err != nil {
		return fmt.Errorf("error deserializing [%s]: %w", ids.RootConfigID, err)
	}

	// Get the special fields
	cfg.HyperdriveDirectory = masterMap[HyperdriveDirectoryID].(string)
	cfg.Version = masterMap[ids.VersionID].(string)

	return nil
}

// Create a copy of this configuration
func (cfg *HyperdriveConfig) CreateCopy() *HyperdriveConfig {
	cfgCopy := NewHyperdriveConfig(cfg.HyperdriveDirectory)
	network := cfg.Network.Value
	clone(cfg, cfgCopy, network)
	return cfgCopy
}

// Handle a network change on all of the parameters
func (cfg *HyperdriveConfig) ChangeNetwork(newNetwork types.Network) {
	// Get the current network
	oldNetwork := cfg.Network.Value
	if oldNetwork == newNetwork {
		return
	}
	cfg.Network.Value = newNetwork
	changeNetwork(cfg, oldNetwork, newNetwork)
}

// Update the default settings for all overwrite-on-upgrade parameters
func (cfg *HyperdriveConfig) UpdateDefaults() {
	updateDefaults(cfg, cfg.Network.Value)
}

// Get all of the settings that have changed between an old config and this config, and get all of the containers that are affected by those changes - also returns whether or not the selected network was changed
func (cfg *HyperdriveConfig) GetChanges(oldConfig *HyperdriveConfig) (*types.ChangedSection, map[types.ContainerID]bool, bool) {
	// Get the changed parameters
	section, changeCount := getChangedSettings(oldConfig, cfg)
	if changeCount == 0 {
		return nil, nil, false
	}
	section.Name = cfg.GetTitle()

	// Get the affected containers
	containers := map[types.ContainerID]bool{}
	getAffectedContainers(section, containers)

	// Check if the network has changed
	changeNetworks := false
	if oldConfig.Network.Value != cfg.Network.Value {
		changeNetworks = true
	}

	// Return everything
	return section, containers, changeNetworks
}

// Checks to see if the current configuration is valid; if not, returns a list of errors
func (cfg *HyperdriveConfig) Validate() []string {
	errors := []string{}

	// Check for illegal blank strings
	/* TODO - this needs to be smarter and ignore irrelevant settings
	for _, param := range config.GetParameters() {
		if param.Type == ParameterType_String && !param.CanBeBlank && param.Value == "" {
			errors = append(errors, fmt.Sprintf("[%s] cannot be blank.", param.Name))
		}
	}

	for name, subconfig := range config.GetSubconfigs() {
		for _, param := range subconfig.GetParameters() {
			if param.Type == ParameterType_String && !param.CanBeBlank && param.Value == "" {
				errors = append(errors, fmt.Sprintf("[%s - %s] cannot be blank.", name, param.Name))
			}
		}
	}
	*/

	// Ensure the selected port numbers are unique. Keeps track of all the errors
	portMap := make(map[uint16]bool)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.LocalBeaconConfig.HttpPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.LocalBeaconConfig.P2pPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.LocalExecutionConfig.HttpPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.LocalExecutionConfig.WebsocketPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.LocalExecutionConfig.EnginePort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.LocalExecutionConfig.P2pPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.Metrics.EcMetricsPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.Metrics.BnMetricsPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.Metrics.Prometheus.Port, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.Metrics.ExporterMetricsPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.Metrics.Grafana.Port, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, cfg.Metrics.DaemonMetricsPort, errors)
	_, errors = addAndCheckForDuplicate(portMap, cfg.LocalBeaconConfig.Lighthouse.P2pQuicPort, errors)

	return errors
}

// ===============
// === Getters ===
// ===============

func (cfg *HyperdriveConfig) GetChainID() uint {
	return cfg.chainID[cfg.Network.Value]
}

// =====================
// === Field Helpers ===
// =====================

// Applies all of the defaults to all of the settings that have them defined
func (cfg *HyperdriveConfig) applyAllDefaults() error {
	network := cfg.Network.Value
	for _, param := range cfg.GetParameters() {
		err := param.SetToDefault(network)
		if err != nil {
			return fmt.Errorf("error setting parameter default: %w", err)
		}
	}

	return nil
}

// Get the list of options for networks to run on
func getNetworkOptions() []*types.ParameterOption[types.Network] {
	options := []*types.ParameterOption[types.Network]{
		{
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Ethereum Mainnet",
				Description: "This is the real Ethereum main network, using real ETH and real RPL to make real validators.",
			},
			Value: types.Network_Mainnet,
		},
		{
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Holesky Testnet",
				Description: "This is the Holešky (Holešovice) test network, which is the next generation of long-lived testnets for Ethereum. It uses free fake ETH and free fake RPL to make fake validators.\nUse this if you want to practice running the Smartnode in a free, safe environment before moving to Mainnet.",
			},
			Value: types.Network_Holesky,
		},
	}

	if strings.HasSuffix(shared.HyperdriveVersion, "-dev") {
		options = append(options, &types.ParameterOption[types.Network]{
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Devnet",
				Description: "This is a development network used by Rocket Pool engineers to test new features and contract upgrades before they are promoted to Holesky for staging. You should not use this network unless invited to do so by the developers.",
			},
			Value: types.Network_HoleskyDev,
		})
	}

	return options
}

// Get a more verbose client description, including warnings
func getAugmentedEcDescription(client types.ExecutionClient, originalDescription string) string {
	switch client {
	case types.ExecutionClient_Nethermind:
		totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
		if totalMemoryGB < 9 {
			return fmt.Sprintf("%s\n\n[red]WARNING: Nethermind currently requires over 8 GB of RAM to run smoothly. We do not recommend it for your system. This may be improved in a future release.", originalDescription)
		}
	case types.ExecutionClient_Besu:
		totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
		if totalMemoryGB < 9 {
			return fmt.Sprintf("%s\n\n[red]WARNING: Besu currently requires over 8 GB of RAM to run smoothly. We do not recommend it for your system. This may be improved in a future release.", originalDescription)
		}
	}

	return originalDescription
}

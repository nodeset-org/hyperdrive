package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/config/migration"
	modconfig "github.com/nodeset-org/hyperdrive/shared/config/modules"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/pbnjay/memory"

	"gopkg.in/yaml.v3"
)

// =========================
// === Hyperdrive Config ===
// =========================

const (
	// Param IDs
	DebugModeID          string = "debugMode"
	NetworkID            string = "network"
	ClientModeID         string = "clientMode"
	UserDataPathID       string = "hdUserDataDir"
	ProjectNameID        string = "projectName"
	AutoTxMaxFeeID       string = "autoTxMaxFee"
	MaxPriorityFeeID     string = "maxPriorityFee"
	AutoTxGasThresholdID string = "autoTxGasThreshold"

	// Tags
	hyperdriveTag string = "nodeset/hyperdrive:v" + shared.HyperdriveVersion

	// Internal fields
	userDirectoryKey string = "hdUserDir"
)

// The master configuration struct
type HyperdriveConfig struct {
	// General settings
	DebugMode          types.Parameter[bool]
	Network            types.Parameter[types.Network]
	ClientMode         types.Parameter[types.ClientMode]
	ProjectName        types.Parameter[string]
	UserDataPath       types.Parameter[string]
	AutoTxMaxFee       types.Parameter[float64]
	MaxPriorityFee     types.Parameter[float64]
	AutoTxGasThreshold types.Parameter[float64]

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

	// Modules
	Modules map[string]any

	// Internal fields
	Version                 string
	HyperdriveUserDirectory string
	chainID                 map[types.Network]uint
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
		HyperdriveUserDirectory: hdDir,

		ProjectName: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ProjectNameID,
				Name:               "Project Name",
				Description:        "This is the prefix that will be attached to all of the Docker containers managed by Hyperdrive.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode, types.ContainerID_Daemon, types.ContainerID_ExecutionClient, types.ContainerID_Exporter, types.ContainerID_Grafana, types.ContainerID_Prometheus, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "hyperdrive",
			},
		},

		Network: types.Parameter[types.Network]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 NetworkID,
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
				ID:                 ClientModeID,
				Name:               "Client Mode",
				Description:        "Choose which mode to use for your Execution Client and Beacon Node - locally managed (Docker Mode), or externally managed (Hybrid Mode).",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ExecutionClient, types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*types.ParameterOption[types.ClientMode]{
				{
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Locally Managed",
						Description: "Allow Hyperdrive to manage the Execution Client and Beacon Node for you (Docker Mode)",
					},
					Value: types.ClientMode_Local,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Externally Managed",
						Description: "Use an existing Execution Client and Beacon Node that you manage on your own (Hybrid Mode)",
					},
					Value: types.ClientMode_External,
				}},
			Default: map[types.Network]types.ClientMode{
				types.Network_All: types.ClientMode_Local,
			},
		},

		AutoTxMaxFee: types.Parameter[float64]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AutoTxMaxFeeID,
				Name:               "Auto TX Max Fee",
				Description:        "Set this if you want all of Hyperdrive's automatic transactions to use this specific max fee value (in gwei), which is the most you'd be willing to pay (*including the priority fee*).\n\nA value of 0 will use the suggested max fee based on the current network conditions.\n\nAny other value will ignore the network suggestion and use this value instead.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]float64{
				types.Network_All: float64(0),
			},
		},

		MaxPriorityFee: types.Parameter[float64]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MaxPriorityFeeID,
				Name:               "Max Priority Fee",
				Description:        "The default value for the priority fee (in gwei) for all of your transactions, including automatic ones. This describes how much you're willing to pay *above the network's current base fee* - the higher this is, the more ETH you give to the validators for including your transaction, which generally means it will be included in a block faster (as long as your max fee is sufficiently high to cover the current network conditions).\n\nMust be larger than 0.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]float64{
				types.Network_All: float64(1),
			},
		},

		AutoTxGasThreshold: types.Parameter[float64]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AutoTxGasThresholdID,
				Name:               "Automatic TX Gas Threshold",
				Description:        "The threshold (in gwei) that the recommended network gas price must be under in order for automated transactions to be submitted when due.\n\nA value of 0 will disable non-essential automatic transactions.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]float64{
				types.Network_All: float64(100),
			},
		},

		UserDataPath: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 UserDataPathID,
				Name:               "User Data Path",
				Description:        "The absolute path of your personal `data` folder that contains secrets such as your node wallet's encrypted file, the password for your node wallet, and all of the validator keys for any Hyperdrive modules.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: filepath.Join(hdDir, "data"),
			},
		},

		DebugMode: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 DebugModeID,
				Name:               "Debug Mode",
				Description:        "Enable debug log printing in the daemon.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: false,
			},
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
	cfg.Network.Value = types.Network_Mainnet
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
		&cfg.ProjectName,
		&cfg.Network,
		&cfg.ClientMode,
		&cfg.AutoTxMaxFee,
		&cfg.MaxPriorityFee,
		&cfg.AutoTxGasThreshold,
		&cfg.UserDataPath,
		&cfg.DebugMode,
	}
}

// Get the subconfigurations for this config
func (cfg *HyperdriveConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{
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

	hdMap := Serialize(cfg)
	masterMap[userDirectoryKey] = cfg.HyperdriveUserDirectory
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
		return fmt.Errorf("config has an entry named [%s] but it is not a map, it's a %s", ids.RootConfigID, reflect.TypeOf(hyperdriveParams))
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
	err = Deserialize(cfg, hdMap, network)
	if err != nil {
		return fmt.Errorf("error deserializing [%s]: %w", ids.RootConfigID, err)
	}

	// Get the special fields
	udKey, exists := masterMap[userDirectoryKey]
	if !exists {
		return fmt.Errorf("expected a user directory parameter named [%s] but it was not found", userDirectoryKey)
	}
	cfg.HyperdriveUserDirectory = udKey.(string)
	version, exists := masterMap[ids.VersionID]
	if !exists {
		return fmt.Errorf("expected a version parameter named [%s] but it was not found", ids.VersionID)
	}
	cfg.Version = version.(string)

	// Handle modules
	modules, exists := masterMap[modconfig.ModulesName]
	if exists {
		if modMap, ok := modules.(map[string]any); ok {
			cfg.Modules = modMap
		} else {
			return fmt.Errorf("config has an entry named [%s] but it is not a map, it's a %s", modconfig.ModulesName, reflect.TypeOf(modules))
		}
	} else {
		cfg.Modules = map[string]any{}
	}

	return nil
}

// Create a copy of this configuration
func (cfg *HyperdriveConfig) CreateCopy() *HyperdriveConfig {
	cfgCopy := NewHyperdriveConfig(cfg.HyperdriveUserDirectory)
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
		return nil, map[types.ContainerID]bool{}, false
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

// =====================
// === Field Helpers ===
// =====================

// Applies all of the defaults to all of the settings that have them defined
func (cfg *HyperdriveConfig) applyAllDefaults() {
	network := cfg.Network.Value
	applyDefaults(cfg, network)
}

// Get the list of options for networks to run on
func getNetworkOptions() []*types.ParameterOption[types.Network] {
	options := []*types.ParameterOption[types.Network]{
		/*{
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Ethereum Mainnet",
				Description: "This is the real Ethereum main network, using real ETH to make real validators.",
			},
			Value: types.Network_Mainnet,
		},*/
		{
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Holesky Testnet",
				Description: "This is the Holešky (Holešovice) test network, which is the next generation of long-lived testnets for Ethereum. It uses free fake ETH to make fake validators.\nUse this if you want to practice running Hyperdrive in a free, safe environment before moving to Mainnet.",
			},
			Value: types.Network_Holesky,
		},
	}

	if strings.HasSuffix(shared.HyperdriveVersion, "-dev") {
		options = append(options, &types.ParameterOption[types.Network]{
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Devnet",
				Description: "This is a development network used by Hyperdrive engineers to test new features and contract upgrades before they are promoted to Holesky for staging. You should not use this network unless invited to do so by the developers.",
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

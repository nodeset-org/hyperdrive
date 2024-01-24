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
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/pbnjay/memory"

	"gopkg.in/yaml.v2"
)

// =========================
// === Hyperdrive Config ===
// =========================

const (
	// Param IDs
	HyperdriveDebugModeID          string = "debugMode"
	HyperdriveNetworkID            string = "network"
	HyperdriveClientModeID         string = "clientMode"
	HyperdriveUseFallbackClientsID string = "useFallbackClients"
	HyperdriveExecutionClientID    string = "executionClient"
	HyperdriveBeaconNodeID         string = "beaconNode"

	// Tags
	HyperdriveTag string = "nodeset/hyperdrive:v" + shared.HyperdriveVersion

	// Internal settings
	HyperdriveDaemonSocketPath string = "data/sockets/daemon.sock"
)

// The master configuration struct
type HyperdriveConfig struct {
	// General settings
	DebugMode          types.Parameter[bool]
	Network            types.Parameter[types.Network]
	ClientMode         types.Parameter[types.ClientMode]
	UseFallbackClients types.Parameter[bool]
	Fallback           *FallbackConfig

	// Execution client settings
	ExecutionClient         types.Parameter[types.ExecutionClient]
	ExecutionCommon         *ExecutionCommonConfig
	Geth                    *GethConfig
	Nethermind              *NethermindConfig
	Besu                    *BesuConfig
	ExternalExecutionConfig *ExternalExecutionConfig

	// Beacon node settings
	BeaconNode              types.Parameter[types.BeaconNode]
	ConsensusCommon         *ConsensusCommonConfig
	Lighthouse              *LighthouseBnConfig
	Lodestar                *LodestarBnConfig
	Nimbus                  *NimbusBnConfig
	Prysm                   *PrysmBnConfig
	Teku                    *TekuBnConfig
	ExternalConsensusConfig *ExternalBeaconConfig

	// Metrics
	Grafana           *GrafanaConfig
	Prometheus        *PrometheusConfig
	Exporter          *ExporterConfig
	BitflyNodeMetrics *BitflyNodeMetricsConfig

	// Internal fields
	Version             string
	HyperdriveDirectory string
	chainID             map[types.Network]uint
}

// Load configuration settings from a file
func LoadFromFile(path string, isDaemon bool) (*HyperdriveConfig, error) {
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
	var settings map[string]string
	if err := yaml.Unmarshal(configBytes, &settings); err != nil {
		return nil, fmt.Errorf("could not parse settings file: %w", err)
	}

	// Deserialize it into a config object
	cfg := NewHyperdriveConfig(filepath.Dir(path))
	err = cfg.Deserialize(settings, isDaemon)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize settings file: %w", err)
	}

	return cfg, nil
}

// Creates a new Hyperdrive configuration instance
func NewHyperdriveConfig(hdDir string) *HyperdriveConfig {
	/*
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "~"
		}
	*/
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

		UseFallbackClients: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HyperdriveUseFallbackClientsID,
				Name:               "Use Fallback Clients",
				Description:        "Enable this if you would like to specify a fallback Execution and Consensus Client, which will temporarily be used by the Smartnode and your Validator Client if your primary Execution / Consensus client pair ever go offline (e.g. if you switch, prune, or resync your clients).",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: false,
			},
		},

		ExecutionClient: types.Parameter[types.ExecutionClient]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HyperdriveExecutionClientID,
				Name:               "Execution Client",
				Description:        "Select which Execution client you would like to run.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*types.ParameterOption[types.ExecutionClient]{
				{
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Geth",
						Description: "Geth is one of the three original implementations of the Ethereum protocol. It is written in Go, fully open source and licensed under the GNU LGPL v3.",
					},
					Value: types.ExecutionClient_Geth,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Nethermind",
						Description: getAugmentedEcDescription(types.ExecutionClient_Nethermind, "Nethermind is a high-performance full Ethereum protocol client with very fast sync speeds. Nethermind is built with proven industrial technologies such as .NET 6 and the Kestrel web server. It is fully open source."),
					},
					Value: types.ExecutionClient_Nethermind,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Besu",
						Description: getAugmentedEcDescription(types.ExecutionClient_Besu, "Hyperledger Besu is a robust full Ethereum protocol client. It uses a novel system called \"Bonsai Trees\" to store its chain data efficiently, which allows it to access block states from the past and does not require pruning. Besu is fully open source and written in Java."),
					},
					Value: types.ExecutionClient_Besu,
				}},
			Default: map[types.Network]types.ExecutionClient{
				types.Network_All: types.ExecutionClient_Geth},
		},

		BeaconNode: types.Parameter[types.BeaconNode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HyperdriveBeaconNodeID,
				Name:               "Consensus Client",
				Description:        "Select which Consensus client you would like to use.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*types.ParameterOption[types.BeaconNode]{
				{
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Lighthouse",
						Description: "Lighthouse is a Consensus client with a heavy focus on speed and security. The team behind it, Sigma Prime, is an information security and software engineering firm who have funded Lighthouse along with the Ethereum Foundation, Consensys, and private individuals. Lighthouse is built in Rust and offered under an Apache 2.0 License.",
					},
					Value: types.ConsensusClient_Lighthouse,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Lodestar",
						Description: "Lodestar is the fifth open-source Ethereum consensus client. It is written in Typescript maintained by ChainSafe Systems. Lodestar, their flagship product, is a production-capable Beacon Chain and Validator Client uniquely situated as the go-to for researchers and developers for rapid prototyping and browser usage.",
					},
					Value: types.ConsensusClient_Lodestar,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Nimbus",
						Description: "Nimbus is a Consensus client implementation that strives to be as lightweight as possible in terms of resources used. This allows it to perform well on embedded systems, resource-restricted devices -- including Raspberry Pis and mobile devices -- and multi-purpose servers.",
					},
					Value: types.ConsensusClient_Nimbus,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Prysm",
						Description: "Prysm is a Go implementation of Ethereum Consensus protocol with a focus on usability, security, and reliability. Prysm is developed by Prysmatic Labs, a company with the sole focus on the development of their client. Prysm is written in Go and released under a GPL-3.0 license.",
					},
					Value: types.ConsensusClient_Prysm,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Teku",
						Description: "PegaSys Teku (formerly known as Artemis) is a Java-based Ethereum 2.0 client designed & built to meet institutional needs and security requirements. PegaSys is an arm of ConsenSys dedicated to building enterprise-ready clients and tools for interacting with the core Ethereum platform. Teku is Apache 2 licensed and written in Java, a language notable for its maturity & ubiquity.",
					},
					Value: types.ConsensusClient_Teku,
				}},
			Default: map[types.Network]types.BeaconNode{
				types.Network_All: types.ConsensusClient_Nimbus,
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
	cfg.ExecutionCommon = NewExecutionCommonConfig(cfg)
	cfg.Geth = NewGethConfig(cfg)
	cfg.Nethermind = NewNethermindConfig(cfg)
	cfg.Besu = NewBesuConfig(cfg)
	cfg.ExternalExecutionConfig = NewExternalExecutionConfig(cfg)
	cfg.ConsensusCommon = NewConsensusCommonConfig(cfg)
	cfg.Lighthouse = NewLighthouseBnConfig(cfg)
	cfg.Lodestar = NewLodestarBnConfig(cfg)
	cfg.Nimbus = NewNimbusBnConfig(cfg)
	cfg.Prysm = NewPrysmBnConfig(cfg)
	cfg.Teku = NewTekuBnConfig(cfg)
	cfg.ExternalConsensusConfig = NewExternalBeaconConfig(cfg)
	cfg.Grafana = NewGrafanaConfig(cfg)
	cfg.Prometheus = NewPrometheusConfig(cfg)
	cfg.Exporter = NewExporterConfig(cfg)
	cfg.BitflyNodeMetrics = NewBitflyNodeMetricsConfig(cfg)

	// Apply the default values for mainnet
	cfg.applyAllDefaults()

	return cfg
}

// Serializes the configuration into a map of maps, compatible with a settings file
func (cfg *HyperdriveConfig) Serialize() map[string]string {
	masterMap := map[string]string{}
	for _, param := range cfg.GetParameters() {
		masterMap[param.GetCommon().ID] = param.GetValueAsString()
	}
	masterMap["hdDir"] = cfg.HyperdriveDirectory
	masterMap["version"] = fmt.Sprintf("v%s", shared.HyperdriveVersion)

	return masterMap
}

// Deserializes a settings file into this config
func (cfg *HyperdriveConfig) Deserialize(masterMap map[string]string, isDaemon bool) error {
	// Get the network
	network := cfg.Network.Value
	// Deserialize root params
	for _, param := range cfg.GetParameters() {
		id := param.GetCommon().ID
		serializedValue, exists := masterMap[id]
		if !exists {
			fmt.Printf("WARN: Parameter [%s] was not found in the config, setting it to the network default.\n", id)
			param.SetToDefault(network)
		} else {
			err := param.Deserialize(serializedValue, network)
			if err != nil {
				return fmt.Errorf("error deserializing parameter [%s]: %w", id, err)
			}
		}
	}

	cfg.HyperdriveDirectory = masterMap["hdDir"]
	cfg.Version = masterMap["version"]
	return nil
}

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

// Get the parameters for this config
func (cfg *HyperdriveConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.DebugMode,
	}
}

// Generates a collection of environment variables based on this config's settings
func (cfg *HyperdriveConfig) GenerateEnvironmentVariables() map[string]string {
	envVars := map[string]string{}

	// Basic variables
	envVars["HYPERDRIVE_IMAGE"] = HyperdriveTag
	envVars["HD_INSTALL_PATH"] = cfg.HyperdriveDirectory

	// Settings
	for _, param := range cfg.GetParameters() {
		vars := param.GetCommon().EnvironmentVariables
		for _, envVar := range vars {
			envVars[envVar] = param.GetValueAsString()
		}
	}

	return envVars
}

func (cfg *HyperdriveConfig) GetChainID() uint {
	return cfg.chainID[cfg.Network.Value]
}

// Get all of the changed settings between an old and new config
func getChangedSettingsMap(oldConfig *HyperdriveConfig, newConfig *HyperdriveConfig) []types.ChangedSetting {
	changedSettings := []types.ChangedSetting{}

	// Root settings
	oldRootParams := oldConfig.GetParameters()
	newRootParams := newConfig.GetParameters()
	changedSettings = getChangedSettings(oldRootParams, newRootParams, newConfig)

	return changedSettings
}

// Get all of the settings that have changed between the given parameter lists.
// Assumes the parameter lists represent identical parameters (e.g. they have the same number of elements and
// each element has the same ID).
func getChangedSettings(oldParams []types.IParameter, newParams []types.IParameter, newConfig *HyperdriveConfig) []types.ChangedSetting {
	changedSettings := []types.ChangedSetting{}

	for i, param := range newParams {
		oldValString := oldParams[i].GetValueAsString()
		newValString := param.GetValueAsString()
		if oldValString != newValString {
			changedSettings = append(changedSettings, types.ChangedSetting{
				Name:               param.GetCommon().Name,
				OldValue:           oldValString,
				NewValue:           newValString,
				AffectedContainers: getAffectedContainers(param, newConfig),
			})
		}
	}

	return changedSettings
}

// Get a list of containers that will be need to be restarted after this change is applied
func getAffectedContainers(param types.IParameter, cfg *HyperdriveConfig) map[types.ContainerID]bool {
	affectedContainers := map[types.ContainerID]bool{}
	for _, container := range param.GetCommon().AffectsContainers {
		affectedContainers[container] = true
	}
	return affectedContainers
}

// Get the possible RPC port mode options
func getPortModes(warningOverride string) []*types.ParameterOption[types.RpcPortMode] {
	if warningOverride == "" {
		warningOverride = "Allow connections from external hosts. This is safe if you're running your node on your local network. If you're a VPS user, this would expose your node to the internet"
	}

	return []*types.ParameterOption[types.RpcPortMode]{
		{
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Closed",
				Description: "Do not allow connections to the port",
			},
			Value: types.RpcPortMode_Closed,
		}, {
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Open to Localhost",
				Description: "Allow connections from this host only",
			},
			Value: types.RpcPortMode_OpenLocalhost,
		}, {
			ParameterOptionCommon: &types.ParameterOptionCommon{
				Name:        "Open to External hosts",
				Description: warningOverride,
			},
			Value: types.RpcPortMode_OpenExternal,
		},
	}
}

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

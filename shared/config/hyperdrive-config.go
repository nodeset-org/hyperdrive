/*
Derived from the Rocket Pool Smartnode source code:
https://github.com/rocket-pool/smartnode
*/

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alessio/shellescape"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/types"

	"gopkg.in/yaml.v2"
)

// =========================
// === Hyperdrive Config ===
// =========================

const (
	HyperdriveTag              string = "nodeset/hyperdrive-stakewise-daemon:v" + shared.HyperdriveStakewiseDaemonVersion
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
	ExternalExecutionConfig *ExternalExecutionConfig

	// Consensus client settings
	ConsensusClient         types.Parameter[types.ConsensusClient]
	ConsensusCommon         *ConsensusCommonConfig
	Nimbus                  *NimbusConfig
	Teku                    *TekuConfig
	ExternalConsensusConfig *ExternalBeaconConfig

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

		DebugMode: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                   "debugMode",
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
	}

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
	masterMap["version"] = fmt.Sprintf("v%s", shared.HyperdriveStakewiseDaemonVersion)

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

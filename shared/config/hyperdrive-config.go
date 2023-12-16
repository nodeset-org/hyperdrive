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
	"github.com/nodeset-org/hyperdrive/shared/types/config"
	rpcfg "github.com/rocket-pool/smartnode/shared/services/config"
	rptypes "github.com/rocket-pool/smartnode/shared/types/config"
	"gopkg.in/yaml.v2"
)

// ===============================
// === Smartnode Config Status ===
// ===============================

type SmartnodeStatus uint8

const (
	SmartnodeStatus_Unknown SmartnodeStatus = iota
	SmartnodeStatus_Loaded
	SmartnodeStatus_MissingCfg
	SmartnodeStatus_InvalidConfig
	SmartnodeStatus_InvalidDir
	SmartnodeStatus_EmptyDir
)

// =========================
// === Hyperdrive Config ===
// =========================

const (
	SmartnodeSettingsFilename  string = "user-settings.yml"
	HyperdriveDaemonSocketPath string = "data/sockets/daemon.sock"
)

// The master configuration struct
type HyperdriveConfig struct {
	Version string `yaml:"-"`

	HyperdriveDirectory string `yaml:"-"`

	SmartnodeStatus SmartnodeStatus `yaml:"-"`

	SmartnodeConfigLoadErrorMessage string `yaml:"-"`

	SmartnodeConfig *rpcfg.RocketPoolConfig `yaml:"-"`

	SmartnodeDirectory *config.Parameter[string] `yaml:"smartnodeDir,omitempty"`

	DaemonSocketPath *config.Parameter[string] `yaml:"daemonSocketPath,omitempty"`

	DebugMode *config.Parameter[bool] `yaml:"debug,omitempty"`
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
	var settings map[string]string
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}
	cfg := &HyperdriveConfig{
		HyperdriveDirectory: hdDir,

		SmartnodeDirectory: &config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                   "smartnodeDir",
				Name:                 "Smartnode Directory",
				Description:          "The directory of the Smartnode installation on this machine that you want to use with Hyperdrive.",
				IsChoice:             false,
				AffectsContainers:    []config.ContainerID{config.ContainerID_Daemon},
				EnvironmentVariables: []string{"HD_SMARTNODE_DIR"},
				CanBeBlank:           false,
				OverwriteOnUpgrade:   false,
			},
			Default: map[rptypes.Network]string{
				rptypes.Network_All: filepath.Join(homeDir, ".rocketpool"),
			},
		},

		DaemonSocketPath: &config.Parameter[string]{
			ParameterCommon: &config.ParameterCommon{
				ID:                   "daemonSocketPath",
				Name:                 "Daemon Socket Path",
				Description:          "The path of the socket file the Daemon will create and listen for API requests on.",
				IsChoice:             false,
				AffectsContainers:    []config.ContainerID{config.ContainerID_Daemon},
				EnvironmentVariables: []string{"HD_DAEMON_SOCKET_PATH"},
				CanBeBlank:           false,
				OverwriteOnUpgrade:   false,
			},
			Default: map[rptypes.Network]string{
				rptypes.Network_All: filepath.Join(hdDir, HyperdriveDaemonSocketPath),
			},
		},

		DebugMode: &config.Parameter[bool]{
			ParameterCommon: &config.ParameterCommon{
				ID:                   "debugMode",
				Name:                 "Debug Mode",
				Description:          "The path of the socket file the Daemon will create and listen for API requests on.",
				IsChoice:             false,
				AffectsContainers:    []config.ContainerID{config.ContainerID_Daemon},
				EnvironmentVariables: []string{"HD_DEBUG_MODE"},
				CanBeBlank:           false,
				OverwriteOnUpgrade:   false,
			},
			Default: map[rptypes.Network]bool{
				rptypes.Network_All: false,
			},
		},
	}

	// Parse the Smartnode config
	cfg.parseSmartnodeConfig()

	// Apply the default values for mainnet
	cfg.applyAllDefaults()

	return cfg
}

// Deserializes a settings file into this config
func (cfg *HyperdriveConfig) Deserialize(masterMap map[string]string) error {
	// Get the network
	network := rptypes.Network_Mainnet
	if cfg.SmartnodeConfig != nil {
		network = cfg.SmartnodeConfig.Smartnode.Network.Value.(rptypes.Network)
	}

	// Deserialize root params
	for _, param := range cfg.GetParameters() {
		serializedValue, exists := masterMap[param.GetCommon().ID]
		if !exists {
			param.SetToDefault(network)
		}
		err := param.Deserialize(serializedValue, network)
		if err != nil {
			return fmt.Errorf("error deserializing root config: %wd", err)
		}
	}

	cfg.Version = masterMap["version"]
	return nil
}

func (cfg *HyperdriveConfig) parseSmartnodeConfig() {
	smartnodeDir := cfg.SmartnodeDirectory.Value

	if smartnodeDir == "" {
		cfg.SmartnodeStatus = SmartnodeStatus_EmptyDir
		return
	}

	// Make sure it exists
	info, err := os.Stat(smartnodeDir)
	if os.IsNotExist(err) {
		cfg.SmartnodeStatus = SmartnodeStatus_InvalidDir
		cfg.SmartnodeConfigLoadErrorMessage = "Directory does not exist."
		return
	}
	if err != nil {
		cfg.SmartnodeStatus = SmartnodeStatus_InvalidDir
		cfg.SmartnodeConfigLoadErrorMessage = err.Error()
		return
	}
	if !info.IsDir() {
		cfg.SmartnodeStatus = SmartnodeStatus_InvalidDir
		cfg.SmartnodeConfigLoadErrorMessage = "The provided Smartnode path is not a directory."
		return
	}

	// Check the Smartnode's config
	smartnodeCfgPath := filepath.Join(smartnodeDir, SmartnodeSettingsFilename)
	info, err = os.Stat(smartnodeCfgPath)
	if os.IsNotExist(err) {
		cfg.SmartnodeStatus = SmartnodeStatus_MissingCfg
		return
	}
	if err != nil {
		cfg.SmartnodeStatus = SmartnodeStatus_InvalidConfig
		cfg.SmartnodeConfigLoadErrorMessage = err.Error()
		return
	}
	if info.IsDir() {
		cfg.SmartnodeStatus = SmartnodeStatus_InvalidConfig
		cfg.SmartnodeConfigLoadErrorMessage = "The Smartnode config path is a directory, not a file."
		return
	}

	// Load the Smartnode's config
	cfg.SmartnodeConfig, err = rpcfg.LoadFromFile(smartnodeCfgPath)
	if err != nil {
		cfg.SmartnodeStatus = SmartnodeStatus_InvalidConfig
		cfg.SmartnodeConfigLoadErrorMessage = err.Error()
		return
	}
	cfg.SmartnodeStatus = SmartnodeStatus_Loaded
}

// Applies all of the defaults to all of the settings that have them defined
func (cfg *HyperdriveConfig) applyAllDefaults() error {
	network := rptypes.Network_Mainnet
	if cfg.SmartnodeConfig != nil {
		network = cfg.SmartnodeConfig.Smartnode.Network.Value.(rptypes.Network)
	}
	for _, param := range cfg.GetParameters() {
		err := param.SetToDefault(network)
		if err != nil {
			return fmt.Errorf("error setting parameter default: %w", err)
		}
	}

	return nil
}

// Get the parameters for this config
func (cfg *HyperdriveConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		cfg.SmartnodeDirectory,
		cfg.DebugMode,
	}
}

// Get all of the changed settings between an old and new config
func getChangedSettingsMap(oldConfig *HyperdriveConfig, newConfig *HyperdriveConfig) []config.ChangedSetting {
	changedSettings := []config.ChangedSetting{}

	// Root settings
	oldRootParams := oldConfig.GetParameters()
	newRootParams := newConfig.GetParameters()
	changedSettings = getChangedSettings(oldRootParams, newRootParams, newConfig)

	return changedSettings
}

// Get all of the settings that have changed between the given parameter lists.
// Assumes the parameter lists represent identical parameters (e.g. they have the same number of elements and
// each element has the same ID).
func getChangedSettings(oldParams []config.IParameter, newParams []config.IParameter, newConfig *HyperdriveConfig) []config.ChangedSetting {
	changedSettings := []config.ChangedSetting{}

	for i, param := range newParams {
		oldValString := oldParams[i].GetValueAsString()
		newValString := param.GetValueAsString()
		if oldValString != newValString {
			changedSettings = append(changedSettings, config.ChangedSetting{
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
func getAffectedContainers(param config.IParameter, cfg *HyperdriveConfig) map[config.ContainerID]bool {
	affectedContainers := map[config.ContainerID]bool{}
	for _, container := range param.GetCommon().AffectsContainers {
		affectedContainers[container] = true
	}
	return affectedContainers
}

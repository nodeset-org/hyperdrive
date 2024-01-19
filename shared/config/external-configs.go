package config

import (
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/types"
)

const (
	// Param IDs
	HttpUrlID string = "httpUrl"
)

// Configuration for external Execution clients
type ExternalExecutionConfig struct {
	Title string

	// The URL of the HTTP endpoint
	HttpUrl types.Parameter[string]
}

// Configuration for external Consensus clients
type ExternalBeaconConfig struct {
	Title string

	// The URL of the HTTP endpoint
	HttpUrl types.Parameter[string]
}

// Generates a new ExternalExecutionConfig configuration
func NewExternalExecutionConfig(cfg *HyperdriveConfig) *ExternalExecutionConfig {
	return &ExternalExecutionConfig{
		Title: "External Execution Client Settings",

		HttpUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HttpUrlID,
				Name:               "HTTP URL",
				Description:        "The URL of the HTTP RPC endpoint for your external Execution client.\nNOTE: If you are running it on the same machine as the Smartnode, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead, for example 'http://192.168.1.100:8545'.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}
}

// Generates a new ExternalBeaconConfig configuration
func NewExternalBeaconConfig(cfg *HyperdriveConfig) *ExternalBeaconConfig {
	return &ExternalBeaconConfig{
		Title: "External Beacon Node Settings",

		HttpUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HttpUrlID,
				Name:               "HTTP URL",
				Description:        "The URL of the HTTP Beacon API endpoint for your external client.\nNOTE: If you are running it on the same machine as the Smartnode, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ValidatorClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}
}

// Get the parameters for this config
func (cfg *ExternalExecutionConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.HttpUrl,
	}
}

// Get the parameters for this config
func (cfg *ExternalBeaconConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.HttpUrl,
	}
}

// The the title for the config
func (cfg *ExternalExecutionConfig) GetConfigTitle() string {
	return cfg.Title
}

// The the title for the config
func (cfg *ExternalBeaconConfig) GetConfigTitle() string {
	return cfg.Title
}

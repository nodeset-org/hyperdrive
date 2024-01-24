package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	HttpUrlID     string = "httpUrl"
	PrysmRpcUrlID string = "prysmRpcUrl"
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

	// The URL of the Prysm gRPC endpoint (only needed if using Prysm VCs)
	PrysmRpcUrl types.Parameter[string]
}

// Generates a new ExternalExecutionConfig configuration
func NewExternalExecutionConfig(cfg *HyperdriveConfig) *ExternalExecutionConfig {
	return &ExternalExecutionConfig{
		Title: "External Execution Client Settings",

		HttpUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HttpUrlID,
				Name:               "HTTP URL",
				Description:        "The URL of the HTTP RPC endpoint for your external Execution client.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead, for example 'http://192.168.1.100:8545'.",
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
				Description:        "The URL of the HTTP Beacon API endpoint for your external client.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		PrysmRpcUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 PrysmRpcUrlID,
				Name:               "RPC URL (Prysm Only)",
				Description:        "**Only used if you have Prysm selected as a Validator Client in one of Hyperdrive's modules.**\n\nThe URL of Prysm's gRPC API endpoint for your external Beacon Node. Prysm's Validator Client will need this in order to connect to it.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
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
		&cfg.PrysmRpcUrl,
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

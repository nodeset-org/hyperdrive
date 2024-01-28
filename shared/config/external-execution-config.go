package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	WebsocketUrlID string = "wsUrl"
)

// Configuration for external Execution clients
type ExternalExecutionConfig struct {
	// The selected EC
	ExecutionClient types.Parameter[types.ExecutionClient]

	// The URL of the HTTP endpoint
	HttpUrl types.Parameter[string]

	// The URL of the Websocket endpoint
	WebsocketUrl types.Parameter[string]

	// Internal Fields
	parent *HyperdriveConfig
}

// Generates a new ExternalExecutionConfig configuration
func NewExternalExecutionConfig(parent *HyperdriveConfig) *ExternalExecutionConfig {
	return &ExternalExecutionConfig{
		parent: parent,

		ExecutionClient: types.Parameter[types.ExecutionClient]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.EcID,
				Name:               "Execution Client",
				Description:        "Select which Execution client your external client is.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*types.ParameterOption[types.ExecutionClient]{
				{
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Geth",
						Description: "Select if your external client is Geth.",
					},
					Value: types.ExecutionClient_Geth,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Nethermind",
						Description: "Select if your external client is Nethermind.",
					},
					Value: types.ExecutionClient_Nethermind,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Besu",
						Description: "Select if your external client is Besu.",
					},
					Value: types.ExecutionClient_Besu,
				}},
			Default: map[types.Network]types.ExecutionClient{
				types.Network_All: types.ExecutionClient_Geth},
		},

		HttpUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.HttpUrlID,
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

		WebsocketUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 WebsocketUrlID,
				Name:               "Websocket URL",
				Description:        "The URL of the Websocket RPC endpoint for your external Execution client.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead, for example 'http://192.168.1.100:8545'.",
				AffectsContainers:  []types.ContainerID{},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}
}

// The the title for the config
func (cfg *ExternalExecutionConfig) GetTitle() string {
	return "External Execution Client Settings"
}

// Get the parameters for this config
func (cfg *ExternalExecutionConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ExecutionClient,
		&cfg.HttpUrl,
		&cfg.WebsocketUrl,
	}
}

// Get the sections underneath this one
func (cfg *ExternalExecutionConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

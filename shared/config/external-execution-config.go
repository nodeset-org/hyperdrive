package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

const (
	// Param IDs
	WebsocketUrlID string = "wsUrl"
)

// Configuration for external Execution clients
type ExternalExecutionConfig struct {
	// The selected EC
	ExecutionClient Parameter[ExecutionClient]

	// The URL of the HTTP endpoint
	HttpUrl Parameter[string]

	// The URL of the Websocket endpoint
	WebsocketUrl Parameter[string]

	// Internal Fields
	parent *HyperdriveConfig
}

// Generates a new ExternalExecutionConfig configuration
func NewExternalExecutionConfig(parent *HyperdriveConfig) *ExternalExecutionConfig {
	return &ExternalExecutionConfig{
		parent: parent,

		ExecutionClient: Parameter[ExecutionClient]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.EcID,
				Name:               "Execution Client",
				Description:        "Select which Execution client your external client is.",
				AffectsContainers:  []ContainerID{ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*ParameterOption[ExecutionClient]{
				{
					ParameterOptionCommon: &ParameterOptionCommon{
						Name:        "Geth",
						Description: "Select if your external client is Geth.",
					},
					Value: ExecutionClient_Geth,
				}, {
					ParameterOptionCommon: &ParameterOptionCommon{
						Name:        "Nethermind",
						Description: "Select if your external client is Nethermind.",
					},
					Value: ExecutionClient_Nethermind,
				}, {
					ParameterOptionCommon: &ParameterOptionCommon{
						Name:        "Besu",
						Description: "Select if your external client is Besu.",
					},
					Value: ExecutionClient_Besu,
				}},
			Default: map[Network]ExecutionClient{
				Network_All: ExecutionClient_Geth},
		},

		HttpUrl: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.HttpUrlID,
				Name:               "HTTP URL",
				Description:        "The URL of the HTTP RPC endpoint for your external Execution client.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead, for example 'http://192.168.1.100:8545'.",
				AffectsContainers:  []ContainerID{ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},

		WebsocketUrl: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 WebsocketUrlID,
				Name:               "Websocket URL",
				Description:        "The URL of the Websocket RPC endpoint for your external Execution client.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead, for example 'http://192.168.1.100:8545'.",
				AffectsContainers:  []ContainerID{},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},
	}
}

// The title for the config
func (cfg *ExternalExecutionConfig) GetTitle() string {
	return "External Execution Client"
}

// Get the parameters for this config
func (cfg *ExternalExecutionConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.ExecutionClient,
		&cfg.HttpUrl,
		&cfg.WebsocketUrl,
	}
}

// Get the sections underneath this one
func (cfg *ExternalExecutionConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

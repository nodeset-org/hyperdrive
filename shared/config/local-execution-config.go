package config

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	EcWebsocketPortID string = "wsPort"
	EcEnginePortID    string = "enginePort"
	EcOpenApiPortsID  string = "openApiPorts"
)

// Configuration for the Execution client
type LocalExecutionConfig struct {
	// The selected EC
	ExecutionClient types.Parameter[types.ExecutionClient]

	// The HTTP API port
	HttpPort types.Parameter[uint16]

	// The Websocket API port
	WebsocketPort types.Parameter[uint16]

	// The Engine API port
	EnginePort types.Parameter[uint16]

	// Toggle for forwarding the HTTP API port outside of Docker
	OpenApiPorts types.Parameter[types.RpcPortMode]

	// P2P traffic port
	P2pPort types.Parameter[uint16]

	// Subconfigs
	Geth       *GethConfig
	Nethermind *NethermindConfig
	Besu       *BesuConfig

	// Internal Fields
	parent *HyperdriveConfig
}

// Create a new LocalExecutionConfig struct
func NewExecutionCommonConfig(parent *HyperdriveConfig) *LocalExecutionConfig {
	cfg := &LocalExecutionConfig{
		parent: parent,

		ExecutionClient: types.Parameter[types.ExecutionClient]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.EcID,
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

		HttpPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.HttpPortID,
				Name:               "HTTP API Port",
				Description:        "The port your Execution client should use for its HTTP API endpoint (also known as HTTP RPC API endpoint).",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ExecutionClient, types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 8545,
			},
		},

		WebsocketPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 EcWebsocketPortID,
				Name:               "Websocket API Port",
				Description:        "The port your Execution client should use for its Websocket API endpoint (also known as Websocket RPC API endpoint).",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 8546,
			},
		},

		EnginePort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 EcEnginePortID,
				Name:               "Engine API Port",
				Description:        "The port your Execution client should use for its Engine API endpoint (the endpoint the Beacon Node will connect to post-merge).",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient, types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 8551,
			},
		},

		OpenApiPorts: types.Parameter[types.RpcPortMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 EcOpenApiPortsID,
				Name:               "Expose API Ports",
				Description:        "Expose the HTTP and Websocket API ports to other processes on your machine, or to your local network so other machines can access your Execution Client's API endpoints.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getPortModes(""),
			Default: map[types.Network]types.RpcPortMode{
				types.Network_All: types.RpcPortMode_Closed,
			},
		},

		P2pPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.P2pPortID,
				Name:               "P2P Port",
				Description:        "The port the Execution Client should use for P2P (blockchain) traffic to communicate with other nodes.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 30303,
			},
		},
	}

	// Create the subconfigs
	cfg.Geth = NewGethConfig(cfg)
	cfg.Nethermind = NewNethermindConfig(cfg)
	cfg.Besu = NewBesuConfig(cfg)

	return cfg
}

// Get the title for the config
func (cfg *LocalExecutionConfig) GetTitle() string {
	return "Local Execution Client Settings"
}

// Get the parameters for this config
func (cfg *LocalExecutionConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.ExecutionClient,
		&cfg.HttpPort,
		&cfg.WebsocketPort,
		&cfg.EnginePort,
		&cfg.OpenApiPorts,
		&cfg.P2pPort,
	}
}

// Get the sections underneath this one
func (cfg *LocalExecutionConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{
		"besu":       cfg.Besu,
		"geth":       cfg.Geth,
		"nethermind": cfg.Nethermind,
	}
}

// ==================
// === Templating ===
// ==================

// Get the Docker mapping for the selected API port mode
func (cfg *LocalExecutionConfig) getOpenApiPortMapping() string {
	rpcMode := cfg.OpenApiPorts.Value
	if !rpcMode.IsOpen() {
		return ""
	}
	httpMapping := rpcMode.DockerPortMapping(cfg.HttpPort.Value)
	wsMapping := rpcMode.DockerPortMapping(cfg.WebsocketPort.Value)
	return fmt.Sprintf(", \"%s\", \"%s\"", httpMapping, wsMapping)
}

// Gets the max peers of the selected EC
func (cfg *LocalExecutionConfig) getMaxPeers() uint16 {
	switch cfg.ExecutionClient.Value {
	case types.ExecutionClient_Geth:
		return cfg.Geth.MaxPeers.Value
	case types.ExecutionClient_Nethermind:
		return cfg.Nethermind.MaxPeers.Value
	case types.ExecutionClient_Besu:
		return cfg.Besu.MaxPeers.Value
	default:
		panic(fmt.Sprintf("Unknown Execution Client %s", string(cfg.ExecutionClient.Value)))
	}
}

// Get the container tag of the selected EC
func (cfg *LocalExecutionConfig) getContainerTag() string {
	switch cfg.ExecutionClient.Value {
	case types.ExecutionClient_Geth:
		return cfg.Geth.ContainerTag.Value
	case types.ExecutionClient_Nethermind:
		return cfg.Nethermind.ContainerTag.Value
	case types.ExecutionClient_Besu:
		return cfg.Besu.ContainerTag.Value
	default:
		panic(fmt.Sprintf("Unknown Execution Client %s", string(cfg.ExecutionClient.Value)))
	}
}

// Gets the additional flags of the selected EC
func (cfg *LocalExecutionConfig) getAdditionalFlags() string {
	switch cfg.ExecutionClient.Value {
	case types.ExecutionClient_Geth:
		return cfg.Geth.AdditionalFlags.Value
	case types.ExecutionClient_Nethermind:
		return cfg.Nethermind.AdditionalFlags.Value
	case types.ExecutionClient_Besu:
		return cfg.Besu.AdditionalFlags.Value
	default:
		panic(fmt.Sprintf("Unknown Execution Client %s", string(cfg.ExecutionClient.Value)))
	}
}

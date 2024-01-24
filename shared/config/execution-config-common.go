package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	EcEnginePortID string = "enginePort"
)

// Configuration for the Execution client
type ExecutionCommonConfig struct {
	Title string

	// The HTTP API port
	HttpPort types.Parameter[uint16]

	// The Engine API port
	EnginePort types.Parameter[uint16]

	// Toggle for forwarding the HTTP API port outside of Docker
	OpenRpcPort types.Parameter[types.RpcPortMode]

	// P2P traffic port
	P2pPort types.Parameter[uint16]
}

// Create a new ExecutionCommonConfig struct
func NewExecutionCommonConfig(cfg *HyperdriveConfig) *ExecutionCommonConfig {
	return &ExecutionCommonConfig{
		Title: "Common Execution Client Settings",

		HttpPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HttpPortID,
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

		EnginePort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 EcEnginePortID,
				Name:               "Engine API Port",
				Description:        "The port your Execution client should use for its Engine API endpoint (the endpoint the Consensus client will connect to post-merge).",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient, types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 8551,
			},
		},

		OpenRpcPort: types.Parameter[types.RpcPortMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 OpenRpcPortID,
				Name:               "Expose RPC Ports",
				Description:        "Expose the HTTP and Websocket RPC ports to other processes on your machine, or to your local network so other machines can access your Execution Client's RPC endpoint.",
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
				ID:                 P2pPortID,
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
}

// Get the parameters for this config
func (cfg *ExecutionCommonConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.HttpPort,
		&cfg.EnginePort,
		&cfg.OpenRpcPort,
		&cfg.P2pPort,
	}
}

// Get the title for the config
func (cfg *ExecutionCommonConfig) GetConfigTitle() string {
	return cfg.Title
}

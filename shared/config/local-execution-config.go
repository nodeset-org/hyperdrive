package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	EcEnginePortID string = "enginePort"
)

// Configuration for the Execution client
type LocalExecutionConfig struct {
	// The HTTP API port
	HttpPort types.Parameter[uint16]

	// The Engine API port
	EnginePort types.Parameter[uint16]

	// Toggle for forwarding the HTTP API port outside of Docker
	OpenRpcPort types.Parameter[types.RpcPortMode]

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
				ID:                 OpenHttpPortID,
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
		&cfg.HttpPort,
		&cfg.EnginePort,
		&cfg.OpenRpcPort,
		&cfg.P2pPort,
	}
}

// Get the sections underneath this one
func (cfg *LocalExecutionConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{
		"besu":       cfg.Besu,
		"geth":       cfg.Geth,
		"nethermind": cfg.Nethermind,
	}
}

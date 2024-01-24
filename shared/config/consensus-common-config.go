package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	BnCheckpointSyncUrlID string = "checkpointSyncUrl"
)

// Common parameters shared by all of the Beacon Clients
type ConsensusCommonConfig struct {
	Title string

	// The checkpoint sync URL if used
	CheckpointSyncProvider types.Parameter[string]

	// The port to use for gossip traffic
	P2pPort types.Parameter[uint16]

	// The port to expose the HTTP API on
	HttpPort types.Parameter[uint16]

	// Toggle for forwarding the HTTP API port outside of Docker
	OpenRpcPort types.Parameter[types.RpcPortMode]
}

// Create a new ConsensusCommonParams struct
func NewConsensusCommonConfig(cfg *HyperdriveConfig) *ConsensusCommonConfig {
	return &ConsensusCommonConfig{
		Title: "Common Consensus Client Settings",

		CheckpointSyncProvider: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:   BnCheckpointSyncUrlID,
				Name: "Checkpoint Sync URL",
				Description: "If you would like to instantly sync using an existing Beacon node, enter its URL.\n" +
					"Example: https://<project ID>:<secret>@eth2-beacon-prater.infura.io\n" +
					"Leave this blank if you want to sync normally from the start of the chain.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		P2pPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 P2pPortID,
				Name:               "P2P Port",
				Description:        "The port to use for P2P (blockchain) traffic.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 9001,
			},
		},

		HttpPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HttpPortID,
				Name:               "HTTP API Port",
				Description:        "The port your Consensus client should run its HTTP API on.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients, types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 5052,
			},
		},

		OpenRpcPort: types.Parameter[types.RpcPortMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 OpenRpcPortID,
				Name:               "Expose API Port",
				Description:        "Select an option to expose your Consensus client's API port to your localhost or external hosts on the network, so other machines can access it too.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getPortModes("Allow connections from external hosts. This is safe if you're running your node on your local network. If you're a VPS user, this would expose your node to the internet and could make it vulnerable to MEV/tips theft"),
			Default: map[types.Network]types.RpcPortMode{
				types.Network_All: types.RpcPortMode_Closed,
			},
		},
	}
}

// Get the parameters for this config
func (cfg *ConsensusCommonConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.CheckpointSyncProvider,
		&cfg.P2pPort,
		&cfg.HttpPort,
		&cfg.OpenRpcPort,
	}
}

// The the title for the config
func (cfg *ConsensusCommonConfig) GetConfigTitle() string {
	return cfg.Title
}

package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	BnCheckpointSyncUrlID string = "checkpointSyncUrl"
)

// Common parameters shared by all of the Beacon Clients
type LocalBeaconConfig struct {
	// The checkpoint sync URL if used
	CheckpointSyncProvider types.Parameter[string]

	// The port to use for gossip traffic
	P2pPort types.Parameter[uint16]

	// The port to expose the HTTP API on
	HttpPort types.Parameter[uint16]

	// Toggle for forwarding the HTTP API port outside of Docker
	OpenHttpPort types.Parameter[types.RpcPortMode]

	// Subconfigs
	Lighthouse *LighthouseBnConfig
	Lodestar   *LodestarBnConfig
	Nimbus     *NimbusBnConfig
	Prysm      *PrysmBnConfig
	Teku       *TekuBnConfig

	// Internal Fields
	parent *HyperdriveConfig
}

// Create a new LocalBeaconConfig struct
func NewLocalBeaconConfig(parent *HyperdriveConfig) *LocalBeaconConfig {
	cfg := &LocalBeaconConfig{
		parent: parent,

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

		OpenHttpPort: types.Parameter[types.RpcPortMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 OpenHttpPortID,
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

	cfg.Lighthouse = NewLighthouseBnConfig(cfg)
	cfg.Lodestar = NewLodestarBnConfig(cfg)
	cfg.Nimbus = NewNimbusBnConfig(cfg)
	cfg.Prysm = NewPrysmBnConfig(cfg)
	cfg.Teku = NewTekuBnConfig(cfg)

	return cfg
}

// The the title for the config
func (cfg *LocalBeaconConfig) GetTitle() string {
	return "Local Beacon Node Settings"
}

// Get the parameters for this config
func (cfg *LocalBeaconConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.CheckpointSyncProvider,
		&cfg.P2pPort,
		&cfg.HttpPort,
		&cfg.OpenHttpPort,
	}
}

// Get the sections underneath this one
func (cfg *LocalBeaconConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{
		"lighthouse": cfg.Lighthouse,
		"lodestar":   cfg.Lodestar,
		"nimbus":     cfg.Nimbus,
		"prysm":      cfg.Prysm,
		"teku":       cfg.Teku,
	}
}

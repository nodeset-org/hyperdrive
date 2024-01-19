package config

import (
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/types"
)

const (
	// Param IDs
	BnGraffitiID              string = "graffiti"
	BnCheckpointSyncUrlID     string = "checkpointSyncUrl"
	BnDoppelgangerDetectionID string = "doppelgangerDetection"
	BnP2pQuicPortID           string = "p2pQuicPort"

	// Reserved for future use
	defaultP2pQuicPort uint16 = 8001
)

// Common parameters shared by all of the Beacon Clients
type ConsensusCommonConfig struct {
	Title string

	// Custom proposal graffiti
	Graffiti types.Parameter[string]

	// The checkpoint sync URL if used
	CheckpointSyncProvider types.Parameter[string]

	// The port to use for gossip traffic
	P2pPort types.Parameter[uint16]

	// The port to expose the HTTP API on
	HttpPort types.Parameter[uint16]

	// Toggle for forwarding the HTTP API port outside of Docker
	OpenRpcPort types.Parameter[types.RpcPortMode]

	// Toggle for enabling doppelganger detection
	DoppelgangerDetection types.Parameter[bool]
}

// Create a new ConsensusCommonParams struct
func NewConsensusCommonConfig(cfg *HyperdriveConfig) *ConsensusCommonConfig {
	return &ConsensusCommonConfig{
		Title: "Common Consensus Client Settings",

		Graffiti: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BnGraffitiID,
				Name:               "Custom Graffiti",
				Description:        "Add a short message to any blocks you propose, so the world can see what you have to say!\nIt has a 16 character limit.",
				MaxLength:          16,
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClient},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

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
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_BeaconNode, types.ContainerID_ValidatorClient, types.ContainerID_Prometheus},
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

		DoppelgangerDetection: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BnDoppelgangerDetectionID,
				Name:               "Enable Doppelgänger Detection",
				Description:        "If enabled, your client will *intentionally* miss 1 or 2 attestations on startup to check if validator keys are already running elsewhere. If they are, it will disable validation duties for them to prevent you from being slashed.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: true,
			},
		},
	}
}

// Get the parameters for this config
func (cfg *ConsensusCommonConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Graffiti,
		&cfg.CheckpointSyncProvider,
		&cfg.P2pPort,
		&cfg.HttpPort,
		&cfg.OpenRpcPort,
		&cfg.DoppelgangerDetection,
	}
}

// The the title for the config
func (cfg *ConsensusCommonConfig) GetConfigTitle() string {
	return cfg.Title
}

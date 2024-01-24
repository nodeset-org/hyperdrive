package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	PrysmRpcPortID     string = "rpcPort"
	PrysmOpenRpcPortID string = "openRpcPort"

	// Tags
	prysmBnTagTest string = "rocketpool/prysm:v4.2.0"
	prysmBnTagProd string = "rocketpool/prysm:v4.2.0"
)

// Configuration for the Prysm BN
type PrysmBnConfig struct {
	Title string

	// The max number of P2P peers to connect to
	MaxPeers types.Parameter[uint16]

	// The RPC port for BN / VC connections
	RpcPort types.Parameter[uint16]

	// Toggle for forwarding the RPC API outside of Docker
	OpenRpcPort types.Parameter[types.RpcPortMode]

	// The Docker Hub tag for the Prysm BN
	ContainerTag types.Parameter[string]

	// Custom command line flags for the BN
	AdditionalFlags types.Parameter[string]
}

// Generates a new Prysm BN configuration
func NewPrysmBnConfig(cfg *HyperdriveConfig) *PrysmBnConfig {
	return &PrysmBnConfig{
		Title: "Prysm Settings",

		MaxPeers: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers your client should try to maintain. You can try lowering this if you have a low-resource system or a constrained network.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 70,
			},
		},

		RpcPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 PrysmRpcPortID,
				Name:               "RPC Port",
				Description:        "The port Prysm should run its JSON-RPC API on.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 5053,
			},
		},

		OpenRpcPort: types.Parameter[types.RpcPortMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 PrysmOpenRpcPortID,
				Name:               "Expose RPC Port",
				Description:        "Expose Prysm's JSON-RPC port to other processes on your machine, or to your local network so other machines can access it too.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getPortModes("Allow connections from external hosts. This is safe if you're running your node on your local network. If you're a VPS user, this would expose your node to the internet and could make it vulnerable to MEV/tips theft"),
			Default: map[types.Network]types.RpcPortMode{
				types.Network_All: types.RpcPortMode_Closed,
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Prysm Beacon Node container on Docker Hub you want to use for the Beacon Node.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    prysmBnTagProd,
				types.Network_HoleskyDev: prysmBnTagTest,
				types.Network_Holesky:    prysmBnTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Prysm's Beacon Node, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}
}

// Get the parameters for this config
func (cfg *PrysmBnConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.MaxPeers,
		&cfg.RpcPort,
		&cfg.OpenRpcPort,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// The the title for the config
func (cfg *PrysmBnConfig) GetConfigTitle() string {
	return cfg.Title
}

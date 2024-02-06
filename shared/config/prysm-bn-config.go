package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

const (
	// Param IDs
	PrysmRpcPortID     string = "rpcPort"
	PrysmOpenRpcPortID string = "openRpcPort"

	// Tags
	prysmBnTagTest string = "rocketpool/prysm:v4.2.1"
	prysmBnTagProd string = "rocketpool/prysm:v4.2.1"
)

// Configuration for the Prysm BN
type PrysmBnConfig struct {
	// The max number of P2P peers to connect to
	MaxPeers Parameter[uint16]

	// The RPC port for BN / VC connections
	RpcPort Parameter[uint16]

	// Toggle for forwarding the RPC API outside of Docker
	OpenRpcPort Parameter[RpcPortMode]

	// The Docker Hub tag for the Prysm BN
	ContainerTag Parameter[string]

	// Custom command line flags for the BN
	AdditionalFlags Parameter[string]

	// Internal Fields
	parent *LocalBeaconConfig
}

// Generates a new Prysm BN configuration
func NewPrysmBnConfig(parent *LocalBeaconConfig) *PrysmBnConfig {
	return &PrysmBnConfig{
		parent: parent,

		MaxPeers: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.MaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers your client should try to maintain. You can try lowering this if you have a low-resource system or a constrained network.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 70,
			},
		},

		RpcPort: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 PrysmRpcPortID,
				Name:               "RPC Port",
				Description:        "The port Prysm should run its JSON-RPC API on.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode, ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 5053,
			},
		},

		OpenRpcPort: Parameter[RpcPortMode]{
			ParameterCommon: &ParameterCommon{
				ID:                 PrysmOpenRpcPortID,
				Name:               "Expose RPC Port",
				Description:        "Expose Prysm's JSON-RPC port to other processes on your machine, or to your local network so other machines can access it too.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getPortModes("Allow connections from external hosts. This is safe if you're running your node on your local network. If you're a VPS user, this would expose your node to the internet and could make it vulnerable to MEV/tips theft"),
			Default: map[Network]RpcPortMode{
				Network_All: RpcPortMode_Closed,
			},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Prysm Beacon Node container on Docker Hub you want to use for the Beacon Node.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_Mainnet:    prysmBnTagProd,
				Network_HoleskyDev: prysmBnTagTest,
				Network_Holesky:    prysmBnTagTest,
			},
		},

		AdditionalFlags: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Prysm's Beacon Node, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},
	}
}

// The title for the config
func (cfg *PrysmBnConfig) GetTitle() string {
	return "Prysm Beacon Node"
}

// Get the parameters for this config
func (cfg *PrysmBnConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.MaxPeers,
		&cfg.RpcPort,
		&cfg.OpenRpcPort,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *PrysmBnConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

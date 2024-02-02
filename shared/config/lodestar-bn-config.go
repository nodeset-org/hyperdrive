package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	lodestarBnTagTest string = "chainsafe/lodestar:v1.12.1"
	lodestarBnTagProd string = "chainsafe/lodestar:v1.12.1"
)

// Configuration for the Lodestar BN
type LodestarBnConfig struct {
	// The max number of P2P peers to connect to
	MaxPeers types.Parameter[uint16]

	// The Docker Hub tag for Lodestar BN
	ContainerTag types.Parameter[string]

	// Custom command line flags for the BN
	AdditionalFlags types.Parameter[string]

	// Internal Fields
	parent *LocalBeaconConfig
}

// Generates a new Lodestar BN configuration
func NewLodestarBnConfig(parent *LocalBeaconConfig) *LodestarBnConfig {
	return &LodestarBnConfig{
		parent: parent,

		MaxPeers: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.MaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers your client should try to maintain. You can try lowering this if you have a low-resource system or a constrained network.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 50,
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Lodestar container from Docker Hub you want to use for the Beacon Node.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    lodestarBnTagProd,
				types.Network_HoleskyDev: lodestarBnTagTest,
				types.Network_Holesky:    lodestarBnTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Lodestar's Beacon Client, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// The title for the config
func (cfg *LodestarBnConfig) GetTitle() string {
	return "Lodestar Settings"
}

// Get the parameters for this config
func (cfg *LodestarBnConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.MaxPeers,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *LodestarBnConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

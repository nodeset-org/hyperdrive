package config

import (
	"fmt"
	"runtime"

	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	NimbusPruningModeID string = "pruningMode"

	// Tags
	nimbusBnTagTest string = "statusim/nimbus-eth2:multiarch-v24.1.1"
	nimbusBnTagProd string = "statusim/nimbus-eth2:multiarch-v24.1.1"
)

// Nimbus's pruning mode
type Nimbus_PruningMode string

const (
	Nimbus_PruningMode_Archive Nimbus_PruningMode = "archive"
	Nimbus_PruningMode_Pruned  Nimbus_PruningMode = "pruned"
)

// Configuration for Nimbus
type NimbusBnConfig struct {
	// The max number of P2P peers to connect to
	MaxPeers types.Parameter[uint16]

	// The Docker Hub tag for the BN
	ContainerTag types.Parameter[string]

	// The pruning mode to use in the BN
	PruningMode types.Parameter[Nimbus_PruningMode]

	// Custom command line flags for the BN
	AdditionalFlags types.Parameter[string]

	// Internal Fields
	parent *LocalBeaconConfig
}

// Generates a new Nimbus configuration
func NewNimbusBnConfig(parent *LocalBeaconConfig) *NimbusBnConfig {
	return &NimbusBnConfig{
		parent: parent,

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
				types.Network_All: getNimbusDefaultPeers(),
			},
		},

		PruningMode: types.Parameter[Nimbus_PruningMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 NimbusPruningModeID,
				Name:               "Pruning Mode",
				Description:        "Choose how Nimbus will prune its database. Highlight each option to learn more about it.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*types.ParameterOption[Nimbus_PruningMode]{
				{
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Pruned",
						Description: "Nimbus will only keep the last 5 months of data available, and will delete everything older than that. This will make Nimbus use less disk space overall, but you won't be able to access state older than 5 months (such as regenerating old rewards trees).\n\n[orange]WARNING: Pruning an *existing* database will take a VERY long time when Nimbus first starts. If you change from Archive to Pruned, you should delete your old chain data and do a checkpoint sync using `rocketpool service resync-eth2`. Make sure you have a checkpoint sync provider specified first!",
					},
					Value: Nimbus_PruningMode_Pruned,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Archive",
						Description: "Nimbus will download the entire Beacon Chain history and store it forever. This is healthier for the overall network, since people will be able to sync the entire chain from scratch using your node.",
					},
					Value: Nimbus_PruningMode_Archive,
				},
			},
			Default: map[types.Network]Nimbus_PruningMode{
				types.Network_All: Nimbus_PruningMode_Pruned,
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Nimbus Beacon Node container you want to use on Docker Hub.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    nimbusBnTagProd,
				types.Network_HoleskyDev: nimbusBnTagTest,
				types.Network_Holesky:    nimbusBnTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Nimbus's Beacon Client, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
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

// Get the title for the config
func (cfg *NimbusBnConfig) GetTitle() string {
	return "Nimbus Settings"
}

// Get the parameters for this config
func (cfg *NimbusBnConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.MaxPeers,
		&cfg.ContainerTag,
		&cfg.PruningMode,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *NimbusBnConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

// Get the default number of peers
func getNimbusDefaultPeers() uint16 {
	switch runtime.GOARCH {
	case "arm64":
		return 100
	case "amd64":
		return 160
	default:
		panic(fmt.Sprintf("unsupported architecture %s", runtime.GOARCH))
	}
}

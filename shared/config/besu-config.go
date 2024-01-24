package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Constants
const (
	// Param IDs
	BesuJvmHeapSizeID   string = "jvmHeapSize"
	BesuMaxPeersID      string = "maxPeers"
	BesuMaxBackLayersID string = "maxBackLayers"

	// Tags
	besuTagTest string = "hyperledger/besu:24.1.0"
	besuTagProd string = "hyperledger/besu:24.1.0"
)

// Configuration for Besu
type BesuConfig struct {
	// Max number of P2P peers to connect to
	JvmHeapSize types.Parameter[uint64]

	// Max number of P2P peers to connect to
	MaxPeers types.Parameter[uint16]

	// Historical state block regeneration limit
	MaxBackLayers types.Parameter[uint64]

	// The Docker Hub tag for Besu
	ContainerTag types.Parameter[string]

	// Custom command line flags
	AdditionalFlags types.Parameter[string]

	// Internal Fields
	parent *LocalExecutionConfig
}

// Generates a new Besu configuration
func NewBesuConfig(parent *LocalExecutionConfig) *BesuConfig {
	return &BesuConfig{
		parent: parent,

		JvmHeapSize: types.Parameter[uint64]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BesuJvmHeapSizeID,
				Name:               "JVM Heap Size",
				Description:        "The max amount of RAM, in MB, that Besu's JVM should limit itself to. Setting this lower will cause Besu to use less RAM, though it will always use more than this limit.\n\nUse 0 for automatic allocation.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint64{
				types.Network_All: uint64(0),
			},
		},

		MaxPeers: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BesuMaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers Besu should connect to. This can be lowered to improve performance on low-power systems or constrained networks. We recommend keeping it at 12 or higher.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 25,
			},
		},

		MaxBackLayers: types.Parameter[uint64]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BesuMaxBackLayersID,
				Name:               "Historical Block Replay Limit",
				Description:        "Besu has the ability to revisit the state of any historical block on the chain by \"replaying\" all of the previous blocks to get back to the target. This limit controls how many blocks you can replay - in other words, how far back Besu can go in time. Normal Execution client processing will be paused while a replay is in progress.\n\n[orange]NOTE: If you try to replay a state from a long time ago, it may take Besu several minutes to rebuild the state!",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint64{
				types.Network_All: uint64(512),
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Besu container you want to use on Docker Hub.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    besuTagProd,
				types.Network_HoleskyDev: besuTagTest,
				types.Network_Holesky:    besuTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass to Besu, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}
}

// The the title for the config
func (cfg *BesuConfig) GetTitle() string {
	return "Besu Settings"
}

// Get the parameters for this config
func (cfg *BesuConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.JvmHeapSize,
		&cfg.MaxPeers,
		&cfg.MaxBackLayers,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *BesuConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

// Constants
const (
	// Param IDs
	BesuJvmHeapSizeID   string = "jvmHeapSize"
	BesuMaxPeersID      string = "maxPeers"
	BesuMaxBackLayersID string = "maxBackLayers"

	// Tags
	besuTagTest string = "hyperledger/besu:24.1.1"
	besuTagProd string = "hyperledger/besu:24.1.1"
)

// Configuration for Besu
type BesuConfig struct {
	// Max number of P2P peers to connect to
	JvmHeapSize Parameter[uint64]

	// Max number of P2P peers to connect to
	MaxPeers Parameter[uint16]

	// Historical state block regeneration limit
	MaxBackLayers Parameter[uint64]

	// The Docker Hub tag for Besu
	ContainerTag Parameter[string]

	// Custom command line flags
	AdditionalFlags Parameter[string]
}

// Generates a new Besu configuration
func NewBesuConfig() *BesuConfig {
	return &BesuConfig{
		JvmHeapSize: Parameter[uint64]{
			ParameterCommon: &ParameterCommon{
				ID:                 BesuJvmHeapSizeID,
				Name:               "JVM Heap Size",
				Description:        "The max amount of RAM, in MB, that Besu's JVM should limit itself to. Setting this lower will cause Besu to use less RAM, though it will always use more than this limit.\n\nUse 0 for automatic allocation.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint64{
				Network_All: uint64(0),
			},
		},

		MaxPeers: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 BesuMaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers Besu should connect to. This can be lowered to improve performance on low-power systems or constrained networks. We recommend keeping it at 12 or higher.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 25,
			},
		},

		MaxBackLayers: Parameter[uint64]{
			ParameterCommon: &ParameterCommon{
				ID:                 BesuMaxBackLayersID,
				Name:               "Historical Block Replay Limit",
				Description:        "Besu has the ability to revisit the state of any historical block on the chain by \"replaying\" all of the previous blocks to get back to the target. This limit controls how many blocks you can replay - in other words, how far back Besu can go in time. Normal Execution client processing will be paused while a replay is in progress.\n\n[orange]NOTE: If you try to replay a state from a long time ago, it may take Besu several minutes to rebuild the state!",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint64{
				Network_All: uint64(512),
			},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Besu container you want to use on Docker Hub.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_Mainnet:    besuTagProd,
				Network_HoleskyDev: besuTagTest,
				Network_Holesky:    besuTagTest,
			},
		},

		AdditionalFlags: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass to Besu, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
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
func (cfg *BesuConfig) GetTitle() string {
	return "Besu"
}

// Get the parameters for this config
func (cfg *BesuConfig) GetParameters() []IParameter {
	return []IParameter{
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

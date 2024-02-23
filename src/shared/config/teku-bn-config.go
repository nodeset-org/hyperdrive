package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/pbnjay/memory"
)

const (
	// Param IDs
	TekuJvmHeapSizeID string = "jvmHeapSize"
	TekuArchiveModeID string = "archiveMode"

	// Tags
	tekuBnTagTest string = "consensys/teku:24.1.1"
	tekuBnTagProd string = "consensys/teku:24.1.1"
)

// Configuration for Teku
type TekuBnConfig struct {
	// Max number of P2P peers to connect to
	JvmHeapSize Parameter[uint64]

	// The max number of P2P peers to connect to
	MaxPeers Parameter[uint16]

	// The archive mode flag
	ArchiveMode Parameter[bool]

	// The Docker Hub tag for the Teku BN
	ContainerTag Parameter[string]

	// Custom command line flags for the BN
	AdditionalFlags Parameter[string]
}

// Generates a new Teku BN configuration
func NewTekuBnConfig() *TekuBnConfig {
	return &TekuBnConfig{
		JvmHeapSize: Parameter[uint64]{
			ParameterCommon: &ParameterCommon{
				ID:                 TekuJvmHeapSizeID,
				Name:               "JVM Heap Size",
				Description:        "The max amount of RAM, in MB, that Teku's JVM should limit itself to. Setting this lower will cause Teku to use less RAM, though it will always use more than this limit.\n\nUse 0 for automatic allocation.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint64{
				Network_All: getTekuHeapSize(),
			},
		},

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
				Network_All: 100,
			},
		},

		ArchiveMode: Parameter[bool]{
			ParameterCommon: &ParameterCommon{
				ID:                 TekuArchiveModeID,
				Name:               "Enable Archive Mode",
				Description:        "When enabled, Teku will run in \"archive\" mode which means it can recreate the state of the Beacon chain for a previous block. This is required for manually generating the Merkle rewards tree.\n\nIf you are sure you will never be manually generating a tree, you can disable archive mode.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]bool{
				Network_All: false,
			},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Teku container on Docker Hub you want to use for the Beacon Node.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_Mainnet:    tekuBnTagProd,
				Network_HoleskyDev: tekuBnTagTest,
				Network_Holesky:    tekuBnTagTest,
			},
		},

		AdditionalFlags: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Teku's Beacon Node, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// Get the title for the config
func (cfg *TekuBnConfig) GetTitle() string {
	return "Teku Beacon Node"
}

// Get the parameters for this config
func (cfg *TekuBnConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.JvmHeapSize,
		&cfg.MaxPeers,
		&cfg.ArchiveMode,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *TekuBnConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

// Get the recommended heap size for Teku
func getTekuHeapSize() uint64 {
	totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
	if totalMemoryGB < 9 {
		return 2048
	}
	return 0
}

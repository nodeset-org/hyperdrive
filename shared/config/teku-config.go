package config

import (
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/types"
	"github.com/pbnjay/memory"
)

const (
	// Param IDs
	TekuJvmHeapSizeID string = "jvmHeapSize"
	TekuArchiveModeID string = "archiveMode"

	// Tags
	tekuTagTest string = "consensys/teku:24.1.0"
	tekuTagProd string = "consensys/teku:24.1.0"
)

// Configuration for Teku
type TekuConfig struct {
	Title string

	// Max number of P2P peers to connect to
	JvmHeapSize types.Parameter[uint64]

	// The max number of P2P peers to connect to
	MaxPeers types.Parameter[uint16]

	// The archive mode flag
	ArchiveMode types.Parameter[bool]

	// The Docker Hub tag for Teku
	ContainerTag types.Parameter[string]

	// Custom command line flags for the BN
	AdditionalBnFlags types.Parameter[string]

	// Custom command line flags for the VC
	AdditionalVcFlags types.Parameter[string]
}

// Generates a new Teku configuration
func NewTekuConfig(cfg *HyperdriveConfig) *TekuConfig {
	return &TekuConfig{
		Title: "Teku Settings",

		JvmHeapSize: types.Parameter[uint64]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 TekuJvmHeapSizeID,
				Name:               "JVM Heap Size",
				Description:        "The max amount of RAM, in MB, that Teku's JVM should limit itself to. Setting this lower will cause Teku to use less RAM, though it will always use more than this limit.\n\nUse 0 for automatic allocation.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint64{
				types.Network_All: getTekuHeapSize(),
			},
		},

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
				types.Network_All: 100,
			},
		},

		ArchiveMode: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 TekuArchiveModeID,
				Name:               "Enable Archive Mode",
				Description:        "When enabled, Teku will run in \"archive\" mode which means it can recreate the state of the Beacon chain for a previous block. This is required for manually generating the Merkle rewards tree.\n\nIf you are sure you will never be manually generating a tree, you can disable archive mode.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: false,
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Teku container you want to use on Docker Hub.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode, types.ContainerID_ValidatorClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    tekuTagProd,
				types.Network_HoleskyDev: tekuTagTest,
				types.Network_Holesky:    tekuTagTest,
			},
		},

		AdditionalBnFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalBnFlagsID,
				Name:               "Additional Beacon Node Flags",
				Description:        "Additional custom command line flags you want to pass Teku's Beacon Node, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		AdditionalVcFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalVcFlagsID,
				Name:               "Additional Validator Client Flags",
				Description:        "Additional custom command line flags you want to pass Teku's Validator Client, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClient},
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
func (cfg *TekuConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.JvmHeapSize,
		&cfg.MaxPeers,
		&cfg.ArchiveMode,
		&cfg.ContainerTag,
		&cfg.AdditionalBnFlags,
		&cfg.AdditionalVcFlags,
	}
}

// Get the recommended heap size for Teku
func getTekuHeapSize() uint64 {
	totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
	if totalMemoryGB < 9 {
		return 2048
	}
	return 0
}

// Get the title for the config
func (cfg *TekuConfig) GetConfigTitle() string {
	return cfg.Title
}

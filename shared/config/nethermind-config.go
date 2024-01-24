package config

import (
	"fmt"
	"runtime"

	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/pbnjay/memory"
)

// Constants
const (
	// Param IDs
	NethermindCacheSizeID         string = "cacheSize"
	NethermindPruneMemSizeID      string = "pruneMemSize"
	NethermindAdditionalModulesID string = "additionalModules"
	NethermindAdditionalUrlsID    string = "additionalUrls"

	// Tags
	nethermindTagProd string = "nethermind/nethermind:1.25.1"
	nethermindTagTest string = "nethermind/nethermind:1.25.1"
)

// Configuration for Nethermind
type NethermindConfig struct {
	Title string

	// Nethermind's cache memory hint
	CacheSize types.Parameter[uint64]

	// Max number of P2P peers to connect to
	MaxPeers types.Parameter[uint16]

	// Nethermind's memory for pruning
	PruneMemSize types.Parameter[uint64]

	// Additional modules to enable on the primary JSON RPC endpoint
	AdditionalModules types.Parameter[string]

	// Additional JSON RPC URLs
	AdditionalUrls types.Parameter[string]

	// The Docker Hub tag for Nethermind
	ContainerTag types.Parameter[string]

	// Custom command line flags
	AdditionalFlags types.Parameter[string]
}

// Generates a new Nethermind configuration
func NewNethermindConfig(cfg *HyperdriveConfig) *NethermindConfig {
	return &NethermindConfig{
		Title: "Nethermind Settings",

		CacheSize: types.Parameter[uint64]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 NethermindCacheSizeID,
				Name:               "Cache (Memory Hint) Size",
				Description:        "The amount of RAM (in MB) you want to suggest for Nethermind's cache. While there is no guarantee that Nethermind will stay under this limit, lower values are preferred for machines with less RAM.\n\nThe default value for this will be calculated dynamically based on your system's available RAM, but you can adjust it manually.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint64{
				types.Network_All: calculateNethermindCache(),
			},
		},

		MaxPeers: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers Nethermind should connect to. This can be lowered to improve performance on low-power systems or constrained types.Networks. We recommend keeping it at 12 or higher.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: calculateNethermindPeers(),
			},
		},

		PruneMemSize: types.Parameter[uint64]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 NethermindPruneMemSizeID,
				Name:               "In-Memory Pruning Cache Size",
				Description:        "The amount of RAM (in MB) you want to dedicate to Nethermind for its in-memory pruning system. Higher values mean less writes to your SSD and slower overall database growth.\n\nThe default value for this will be calculated dynamically based on your system's available RAM, but you can adjust it manually.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint64{
				types.Network_All: calculateNethermindPruneMemSize(),
			},
		},

		AdditionalModules: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 NethermindAdditionalModulesID,
				Name:               "Additional Modules",
				Description:        "Additional modules you want to add to the primary JSON-RPC route. The defaults are Eth,Net,Personal,Web3. You can add any additional ones you need here; separate multiple modules with commas, and do not use spaces.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		AdditionalUrls: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 NethermindAdditionalUrlsID,
				Name:               "Additional URLs",
				Description:        "Additional JSON-RPC URLs you want to run alongside the primary URL. These will be added to the \"--JsonRpc.AdditionalRpcUrls\" argument. Wrap each additional URL in quotes, and separate multiple URLs with commas (no spaces). Please consult the Nethermind documentation for more information on this flag, its intended usage, and its expected formatting.\n\nFor advanced users only.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Nethermind container you want to use on Docker Hub.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    nethermindTagProd,
				types.Network_HoleskyDev: nethermindTagTest,
				types.Network_Holesky:    nethermindTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass to Nethermind, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
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

// Calculate the recommended size for Nethermind's cache based on the amount of system RAM
func calculateNethermindCache() uint64 {
	totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024

	if totalMemoryGB == 0 {
		return 0
	} else if totalMemoryGB < 9 {
		return 512
	} else if totalMemoryGB < 13 {
		return 512
	} else if totalMemoryGB < 17 {
		return 1024
	} else if totalMemoryGB < 25 {
		return 1024
	} else if totalMemoryGB < 33 {
		return 1024
	} else {
		return 2048
	}
}

// Calculate the recommended size for Nethermind's in-memory pruning based on the amount of system RAM
func calculateNethermindPruneMemSize() uint64 {
	totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024

	if totalMemoryGB == 0 {
		return 0
	} else if totalMemoryGB < 9 {
		return 512
	} else if totalMemoryGB < 13 {
		return 512
	} else if totalMemoryGB < 17 {
		return 1024
	} else if totalMemoryGB < 25 {
		return 1024
	} else if totalMemoryGB < 33 {
		return 1024
	} else {
		return 1024
	}
}

// Calculate the default number of Nethermind peers
func calculateNethermindPeers() uint16 {
	switch runtime.GOARCH {
	case "arm64":
		return 25
	case "amd64":
		return 50
	default:
		panic(fmt.Sprintf("unsupported architecture %s", runtime.GOARCH))
	}
}

// Get the parameters for this config
func (cfg *NethermindConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.CacheSize,
		&cfg.MaxPeers,
		&cfg.PruneMemSize,
		&cfg.AdditionalModules,
		&cfg.AdditionalUrls,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the title for the config
func (cfg *NethermindConfig) GetConfigTitle() string {
	return cfg.Title
}

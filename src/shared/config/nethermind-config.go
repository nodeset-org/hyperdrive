package config

import (
	"fmt"
	"runtime"

	"github.com/nodeset-org/hyperdrive/shared/config/ids"
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
	nethermindTagProd string = "nethermind/nethermind:1.25.3"
	nethermindTagTest string = "nethermind/nethermind:1.25.3"
)

// Configuration for Nethermind
type NethermindConfig struct {
	// Nethermind's cache memory hint
	CacheSize Parameter[uint64]

	// Max number of P2P peers to connect to
	MaxPeers Parameter[uint16]

	// Nethermind's memory for pruning
	PruneMemSize Parameter[uint64]

	// Additional modules to enable on the primary JSON RPC endpoint
	AdditionalModules Parameter[string]

	// Additional JSON RPC URLs
	AdditionalUrls Parameter[string]

	// The Docker Hub tag for Nethermind
	ContainerTag Parameter[string]

	// Custom command line flags
	AdditionalFlags Parameter[string]
}

// Generates a new Nethermind configuration
func NewNethermindConfig() *NethermindConfig {
	return &NethermindConfig{
		CacheSize: Parameter[uint64]{
			ParameterCommon: &ParameterCommon{
				ID:                 NethermindCacheSizeID,
				Name:               "Cache (Memory Hint) Size",
				Description:        "The amount of RAM (in MB) you want to suggest for Nethermind's cache. While there is no guarantee that Nethermind will stay under this limit, lower values are preferred for machines with less RAM.\n\nThe default value for this will be calculated dynamically based on your system's available RAM, but you can adjust it manually.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint64{
				Network_All: calculateNethermindCache(),
			},
		},

		MaxPeers: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.MaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers Nethermind should connect to. This can be lowered to improve performance on low-power systems or constrained Networks. We recommend keeping it at 12 or higher.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: calculateNethermindPeers(),
			},
		},

		PruneMemSize: Parameter[uint64]{
			ParameterCommon: &ParameterCommon{
				ID:                 NethermindPruneMemSizeID,
				Name:               "In-Memory Pruning Cache Size",
				Description:        "The amount of RAM (in MB) you want to dedicate to Nethermind for its in-memory pruning system. Higher values mean less writes to your SSD and slower overall database growth.\n\nThe default value for this will be calculated dynamically based on your system's available RAM, but you can adjust it manually.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint64{
				Network_All: calculateNethermindPruneMemSize(),
			},
		},

		AdditionalModules: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 NethermindAdditionalModulesID,
				Name:               "Additional Modules",
				Description:        "Additional modules you want to add to the primary JSON-RPC route. The defaults are Eth,Net,Personal,Web3. You can add any additional ones you need here; separate multiple modules with commas, and do not use spaces.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},

		AdditionalUrls: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 NethermindAdditionalUrlsID,
				Name:               "Additional URLs",
				Description:        "Additional JSON-RPC URLs you want to run alongside the primary URL. These will be added to the \"--JsonRpc.AdditionalRpcUrls\" argument. Wrap each additional URL in quotes, and separate multiple URLs with commas (no spaces). Please consult the Nethermind documentation for more information on this flag, its intended usage, and its expected formatting.\n\nFor advanced users only.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Nethermind container you want to use on Docker Hub.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_Mainnet:    nethermindTagProd,
				Network_HoleskyDev: nethermindTagTest,
				Network_Holesky:    nethermindTagTest,
			},
		},

		AdditionalFlags: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass to Nethermind, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// Get the title for the config
func (cfg *NethermindConfig) GetTitle() string {
	return "Nethermind"
}

// Get the parameters for this config
func (cfg *NethermindConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.CacheSize,
		&cfg.MaxPeers,
		&cfg.PruneMemSize,
		&cfg.AdditionalModules,
		&cfg.AdditionalUrls,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *NethermindConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
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

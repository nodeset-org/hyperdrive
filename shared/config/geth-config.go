package config

import (
	"fmt"
	"runtime"

	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Constants
const (
	// Param IDs
	GethEnablePbssID string = "enablePbss"

	// Tags
	gethTagProd string = "ethereum/client-go:v1.13.10"
	gethTagTest string = "ethereum/client-go:v1.13.10"
)

// Configuration for Geth
type GethConfig struct {
	// The flag for enabling PBSS
	EnablePbss types.Parameter[bool]

	// Max number of P2P peers to connect to
	MaxPeers types.Parameter[uint16]

	// The Docker Hub tag for Geth
	ContainerTag types.Parameter[string]

	// Custom command line flags
	AdditionalFlags types.Parameter[string]

	// Internal Fields
	parent *LocalExecutionConfig
}

// Generates a new Geth configuration
func NewGethConfig(parent *LocalExecutionConfig) *GethConfig {
	return &GethConfig{
		parent: parent,

		EnablePbss: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 GethEnablePbssID,
				Name:               "Enable PBSS",
				Description:        "Enable Geth's new path-based state scheme. With this enabled, you will no longer need to manually prune Geth; it will automatically prune its database in real-time.\n\n[orange]NOTE:\nEnabling this will require you to remove and resync your Geth DB using `hyperdrive service resync-eth1`.\nYou will need a synced fallback node configured before doing this, or you will no longer be able to attest until it has finished resyncing!",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: true,
			},
		},

		MaxPeers: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.MaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers Geth should connect to. This can be lowered to improve performance on low-power systems or constrained types.Networks. We recommend keeping it at 12 or higher.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{types.Network_All: calculateGethPeers()},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Geth container you want to use on Docker Hub.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    gethTagProd,
				types.Network_HoleskyDev: gethTagTest,
				types.Network_Holesky:    gethTagTest,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass to Geth, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// Get the title for the config
func (cfg *GethConfig) GetTitle() string {
	return "Geth"
}

// Get the parameters for this config
func (cfg *GethConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.EnablePbss,
		&cfg.MaxPeers,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *GethConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

// Calculate the default number of Geth peers
func calculateGethPeers() uint16 {
	switch runtime.GOARCH {
	case "arm64":
		return 25
	case "amd64":
		return 50
	default:
		panic(fmt.Sprintf("unsupported architecture %s", runtime.GOARCH))
	}
}

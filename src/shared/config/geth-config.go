package config

import (
	"fmt"
	"runtime"

	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

// Constants
const (
	// Param IDs
	GethEnablePbssID string = "enablePbss"

	// Tags
	gethTagProd string = "ethereum/client-go:v1.13.11"
	gethTagTest string = "ethereum/client-go:v1.13.11"
)

// Configuration for Geth
type GethConfig struct {
	// The flag for enabling PBSS
	EnablePbss Parameter[bool]

	// Max number of P2P peers to connect to
	MaxPeers Parameter[uint16]

	// The Docker Hub tag for Geth
	ContainerTag Parameter[string]

	// Custom command line flags
	AdditionalFlags Parameter[string]
}

// Generates a new Geth configuration
func NewGethConfig() *GethConfig {
	return &GethConfig{
		EnablePbss: Parameter[bool]{
			ParameterCommon: &ParameterCommon{
				ID:                 GethEnablePbssID,
				Name:               "Enable PBSS",
				Description:        "Enable Geth's new path-based state scheme. With this enabled, you will no longer need to manually prune Geth; it will automatically prune its database in real-time.\n\n[orange]NOTE:\nEnabling this will require you to remove and resync your Geth DB using `hyperdrive service resync-eth1`.\nYou will need a synced fallback node configured before doing this, or you will no longer be able to attest until it has finished resyncing!",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]bool{
				Network_All: true,
			},
		},

		MaxPeers: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.MaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers Geth should connect to. This can be lowered to improve performance on low-power systems or constrained Networks. We recommend keeping it at 12 or higher.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{Network_All: calculateGethPeers()},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Geth container you want to use on Docker Hub.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_Mainnet:    gethTagProd,
				Network_HoleskyDev: gethTagTest,
				Network_Holesky:    gethTagTest,
			},
		},

		AdditionalFlags: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass to Geth, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *GethConfig) GetTitle() string {
	return "Geth"
}

// Get the parameters for this config
func (cfg *GethConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.EnablePbss,
		&cfg.MaxPeers,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *GethConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
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

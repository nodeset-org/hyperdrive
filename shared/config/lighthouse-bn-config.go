package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/utils/sys"
)

const (
	// Param IDs
	LhQuicPortID string = "p2pQuicPort"

	// Tags
	lighthouseBnTagPortableTest string = "sigp/lighthouse:v4.5.0"
	lighthouseBnTagPortableProd string = "sigp/lighthouse:v4.5.0"
	lighthouseBnTagModernTest   string = "sigp/lighthouse:v4.5.0-modern"
	lighthouseBnTagModernProd   string = "sigp/lighthouse:v4.5.0-modern"
)

// Configuration for the Lighthouse BN
type LighthouseBnConfig struct {
	// The port to use for gossip traffic using the QUIC protocol
	P2pQuicPort types.Parameter[uint16]

	// The max number of P2P peers to connect to
	MaxPeers types.Parameter[uint16]

	// The Docker Hub tag for Lighthouse BN
	ContainerTag types.Parameter[string]

	// Custom command line flags for the BN
	AdditionalFlags types.Parameter[string]

	// Internal Flags
	parent *LocalBeaconConfig
}

// Generates a new Lighthouse BN configuration
func NewLighthouseBnConfig(parent *LocalBeaconConfig) *LighthouseBnConfig {
	return &LighthouseBnConfig{
		parent: parent,

		P2pQuicPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 LhQuicPortID,
				Name:               "P2P QUIC Port",
				Description:        "The port to use for P2P (blockchain) traffic using the QUIC protocol.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 8001,
			},
		},

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
				types.Network_All: 80,
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Lighthouse container from Docker Hub you want to use for the Beacon Node.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_Mainnet:    getLighthouseBnTagProd(),
				types.Network_HoleskyDev: getLighthouseBnTagTest(),
				types.Network_Holesky:    getLighthouseBnTagTest(),
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass Lighthouse's Beacon Node, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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

// The the title for the config
func (cfg *LighthouseBnConfig) GetTitle() string {
	return "Lighthouse Settings"
}

// Get the parameters for this config
func (cfg *LighthouseBnConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.MaxPeers,
		&cfg.P2pQuicPort,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *LighthouseBnConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

// Get the appropriate LH default tag for production
func getLighthouseBnTagProd() string {
	missingFeatures := sys.GetMissingModernCpuFeatures()
	if len(missingFeatures) > 0 {
		return lighthouseBnTagPortableProd
	}
	return lighthouseBnTagModernProd
}

// Get the appropriate LH default tag for testnets
func getLighthouseBnTagTest() string {
	missingFeatures := sys.GetMissingModernCpuFeatures()
	if len(missingFeatures) > 0 {
		return lighthouseBnTagPortableTest
	}
	return lighthouseBnTagModernTest
}

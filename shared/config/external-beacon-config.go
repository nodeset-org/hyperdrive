package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	PrysmRpcUrlID string = "prysmRpcUrl"
)

// Configuration for external Consensus clients
type ExternalBeaconConfig struct {
	Title string

	// The URL of the HTTP endpoint
	HttpUrl types.Parameter[string]

	// The URL of the Prysm gRPC endpoint (only needed if using Prysm VCs)
	PrysmRpcUrl types.Parameter[string]

	// Internal Fields
	parent *HyperdriveConfig
}

// Generates a new ExternalBeaconConfig configuration
func NewExternalBeaconConfig(parent *HyperdriveConfig) *ExternalBeaconConfig {
	return &ExternalBeaconConfig{
		parent: parent,

		HttpUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 HttpUrlID,
				Name:               "HTTP URL",
				Description:        "The URL of the HTTP Beacon API endpoint for your external client.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		PrysmRpcUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 PrysmRpcUrlID,
				Name:               "RPC URL (Prysm Only)",
				Description:        "**Only used if you have Prysm selected as a Validator Client in one of Hyperdrive's modules.**\n\nThe URL of Prysm's gRPC API endpoint for your external Beacon Node. Prysm's Validator Client will need this in order to connect to it.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}
}

// The the title for the config
func (cfg *ExternalBeaconConfig) GetTitle() string {
	return "External Beacon Node Settings"
}

// Get the parameters for this config
func (cfg *ExternalBeaconConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.HttpUrl,
		&cfg.PrysmRpcUrl,
	}
}

// Get the sections underneath this one
func (cfg *ExternalBeaconConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

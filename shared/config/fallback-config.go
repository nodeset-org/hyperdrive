package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	EcHttpUrl string = "ecHttpUrl"
	BnHttpUrl string = "bnHttpUrl"
)

// Fallback configuration
type FallbackConfig struct {
	Title string

	// The URL of the Execution Client HTTP endpoint
	EcHttpUrl types.Parameter[string]

	// The URL of the Beacon Node HTTP endpoint
	BnHttpUrl types.Parameter[string]

	// The URL of the Prysm gRPC endpoint (only needed if using Prysm VCs)
	PrysmRpcUrl types.Parameter[string]
}

// Generates a new FallbackConfig configuration
func NewFallbackConfig(cfg *HyperdriveConfig) *FallbackConfig {
	return &FallbackConfig{
		Title: "Fallback Client Settings",

		EcHttpUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 EcHttpUrl,
				Name:               "Execution Client URL",
				Description:        "The URL of the HTTP API endpoint for your fallback Execution client.\n\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		BnHttpUrl: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BnHttpUrl,
				Name:               "Beacon Node URL",
				Description:        "The URL of the HTTP Beacon API endpoint for your fallback Beacon Node.\n\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
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
				Description:        "**Only used if you have Prysm selected as a Validator Client in one of Hyperdrive's modules.**\n\nThe URL of Prysm's gRPC API endpoint for your fallback Beacon Node. Prysm's Validator Client will need this in order to connect to it.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
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

// Get the types.Parameters for this config
func (cfg *FallbackConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.EcHttpUrl,
		&cfg.BnHttpUrl,
		&cfg.PrysmRpcUrl,
	}
}

// The the title for the config
func (cfg *FallbackConfig) GetConfigTitle() string {
	return cfg.Title
}

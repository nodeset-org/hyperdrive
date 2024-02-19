package config

const (
	// Param IDs
	UseFallbackClientsID string = "useFallbackClients"
	EcHttpUrl            string = "ecHttpUrl"
	BnHttpUrl            string = "bnHttpUrl"
)

// Fallback configuration
type FallbackConfig struct {
	// Flag for enabling fallback clients
	UseFallbackClients Parameter[bool]

	// The URL of the Execution Client HTTP endpoint
	EcHttpUrl Parameter[string]

	// The URL of the Beacon Node HTTP endpoint
	BnHttpUrl Parameter[string]

	// The URL of the Prysm gRPC endpoint (only needed if using Prysm VCs)
	PrysmRpcUrl Parameter[string]
}

// Generates a new FallbackConfig configuration
func NewFallbackConfig() *FallbackConfig {
	return &FallbackConfig{
		UseFallbackClients: Parameter[bool]{
			ParameterCommon: &ParameterCommon{
				ID:                 UseFallbackClientsID,
				Name:               "Use Fallback Clients",
				Description:        "Enable this if you would like to specify a fallback Execution and Beacon Node, which will temporarily be used by Hyperdrive and your Validator Client if your primary Execution / Beacon Node pair ever go offline (e.g. if you switch, prune, or resync your clients).",
				AffectsContainers:  []ContainerID{ContainerID_Daemon, ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]bool{
				Network_All: false,
			},
		},

		EcHttpUrl: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 EcHttpUrl,
				Name:               "Execution Client URL",
				Description:        "The URL of the HTTP API endpoint for your fallback Execution client.\n\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []ContainerID{ContainerID_Daemon},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},

		BnHttpUrl: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 BnHttpUrl,
				Name:               "Beacon Node URL",
				Description:        "The URL of the HTTP Beacon API endpoint for your fallback Beacon Node.\n\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []ContainerID{ContainerID_Daemon, ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},

		PrysmRpcUrl: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 PrysmRpcUrlID,
				Name:               "RPC URL (Prysm Only)",
				Description:        "**Only used if you have Prysm selected as a Validator Client in one of Hyperdrive's modules.**\n\nThe URL of Prysm's gRPC API endpoint for your fallback Beacon Node. Prysm's Validator Client will need this in order to connect to it.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
				AffectsContainers:  []ContainerID{ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},
	}
}

// The title for the config
func (cfg *FallbackConfig) GetTitle() string {
	return "Fallback Clients"
}

// Get the Parameters for this config
func (cfg *FallbackConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.UseFallbackClients,
		&cfg.EcHttpUrl,
		&cfg.BnHttpUrl,
		&cfg.PrysmRpcUrl,
	}
}

// Get the sections underneath this one
func (cfg *FallbackConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

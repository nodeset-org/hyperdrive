package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Defaults
const (
	BitflySecretID      string = "bitflySecret"
	BitflyEndpointID    string = "bitflyEndpoint"
	BitflyMachineNameID string = "bitflyMachineName"
)

// Configuration for Bitfly Node Metrics
type BitflyNodeMetricsConfig struct {
	Secret types.Parameter[string]

	Endpoint types.Parameter[string]

	MachineName types.Parameter[string]

	// Internal Fields
	parent *MetricsConfig
}

// Generates a new Bitfly Node Metrics config
func NewBitflyNodeMetricsConfig(parent *MetricsConfig) *BitflyNodeMetricsConfig {
	return &BitflyNodeMetricsConfig{
		parent: parent,

		Secret: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                BitflySecretID,
				Name:              "Beaconcha.in API Key",
				Description:       "The API key used to authenticate your Beaconcha.in node metrics integration. Can be found in your Beaconcha.in account settings.\n\nPlease visit https://beaconcha.in/user/settings#api to access your account information.",
				AffectsContainers: []types.ContainerID{types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients},
				// ensures the string is 28 characters of Base64
				Regex:              "^[A-Za-z0-9+/]{28}$",
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},

		Endpoint: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BitflyEndpointID,
				Name:               "Node Metrics Endpoint",
				Description:        "The endpoint to send your Beaconcha.in Node Metrics data to. Should be left as the default.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "https://beaconcha.in/api/v1/client/metrics",
			},
		},

		MachineName: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BitflyMachineNameID,
				Name:               "Node Metrics Machine Name",
				Description:        "The name of the machine you are running on. This is used to identify your machine in the mobile app.\nChange this if you are running multiple Hyperdrives with the same Secret.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient, types.ContainerID_ValidatorClients},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "Hyperdrive Node",
			},
		},
	}
}

// The title for the config
func (cfg *BitflyNodeMetricsConfig) GetTitle() string {
	return "Bitfly Node Metrics"
}

// Get the parameters for this config
func (cfg *BitflyNodeMetricsConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Secret,
		&cfg.Endpoint,
		&cfg.MachineName,
	}
}

// Get the sections underneath this one
func (cfg *BitflyNodeMetricsConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

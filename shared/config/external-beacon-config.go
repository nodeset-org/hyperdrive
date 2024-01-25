package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	PrysmRpcUrlID string = "prysmRpcUrl"
)

// Configuration for external Beacon Nodes
type ExternalBeaconConfig struct {
	// The selected BN
	BeaconNode types.Parameter[types.BeaconNode]

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

		BeaconNode: types.Parameter[types.BeaconNode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 BnID,
				Name:               "Beacon Node",
				Description:        "Select which Beacon Node your external client is.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: []*types.ParameterOption[types.BeaconNode]{
				{
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Lighthouse",
						Description: "Select if your external client is Lighthouse.",
					},
					Value: types.BeaconNode_Lighthouse,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Lodestar",
						Description: "Select if your external client is Lodestar.",
					},
					Value: types.BeaconNode_Lodestar,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Nimbus",
						Description: "Select if your external client is Nimbus.",
					},
					Value: types.BeaconNode_Nimbus,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Prysm",
						Description: "Select if your external client is Prysm.",
					},
					Value: types.BeaconNode_Prysm,
				}, {
					ParameterOptionCommon: &types.ParameterOptionCommon{
						Name:        "Teku",
						Description: "Select if your external client is Teku.",
					},
					Value: types.BeaconNode_Teku,
				}},
			Default: map[types.Network]types.BeaconNode{
				types.Network_All: types.BeaconNode_Nimbus,
			},
		},

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
				Name:               "Prysm RPC URL",
				Description:        "The URL of Prysm's gRPC API endpoint for your external Beacon Node. Prysm's Validator Client will need this in order to connect to it.\nNOTE: If you are running it on the same machine as Hyperdrive, addresses like `localhost` and `127.0.0.1` will not work due to Docker limitations. Enter your machine's LAN IP address instead.",
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
		&cfg.BeaconNode,
		&cfg.HttpUrl,
		&cfg.PrysmRpcUrl,
	}
}

// Get the sections underneath this one
func (cfg *ExternalBeaconConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

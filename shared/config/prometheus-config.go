package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Tags
	prometheusTag string = "prom/prometheus:v2.47.1"
)

// Configuration for Prometheus
type PrometheusConfig struct {
	Title string

	// The port to serve metrics on
	Port types.Parameter[uint16]

	// Toggle for forwarding the API port outside of Docker
	OpenPort types.Parameter[types.RpcPortMode]

	// The Docker Hub tag for Prometheus
	ContainerTag types.Parameter[string]

	// Custom command line flags
	AdditionalFlags types.Parameter[string]
}

// Generates a new Prometheus config
func NewPrometheusConfig(cfg *HyperdriveConfig) *PrometheusConfig {
	return &PrometheusConfig{
		Title: "Prometheus Settings",

		Port: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 PortID,
				Name:               "Prometheus Port",
				Description:        "The port Prometheus should make its statistics available on.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Prometheus},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 9091,
			},
		},

		OpenPort: types.Parameter[types.RpcPortMode]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 OpenPortID,
				Name:               "Expose Prometheus Port",
				Description:        "Expose the Prometheus's port to other processes on your machine, or to your local network so other machines can access it too.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getPortModes(""),
			Default: map[types.Network]types.RpcPortMode{
				types.Network_All: types.RpcPortMode_Closed,
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ContainerTagID,
				Name:               "Prometheus Container Tag",
				Description:        "The tag name of the Prometheus container on Docker Hub you want to use.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_All: prometheusTag,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 AdditionalFlagsID,
				Name:               "Additional Prometheus Flags",
				Description:        "Additional custom command line flags you want to pass to Prometheus, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Grafana},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]string{
				types.Network_All: "",
			},
		},
	}
}

// Get the parameters for this config
func (cfg *PrometheusConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Port,
		&cfg.OpenPort,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// The the title for the config
func (cfg *PrometheusConfig) GetConfigTitle() string {
	return cfg.Title
}

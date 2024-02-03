package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Tags
	prometheusTag string = "prom/prometheus:v2.49.1"
)

// Configuration for Prometheus
type PrometheusConfig struct {
	// The port to serve metrics on
	Port types.Parameter[uint16]

	// Toggle for forwarding the API port outside of Docker
	OpenPort types.Parameter[types.RpcPortMode]

	// The Docker Hub tag for Prometheus
	ContainerTag types.Parameter[string]

	// Custom command line flags
	AdditionalFlags types.Parameter[string]

	// Internal Fields
	parent *MetricsConfig
}

// Generates a new Prometheus config
func NewPrometheusConfig(parent *MetricsConfig) *PrometheusConfig {
	return &PrometheusConfig{
		parent: parent,

		Port: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.PortID,
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
				ID:                 ids.OpenPortID,
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
				ID:                 ids.ContainerTagID,
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
				ID:                 ids.AdditionalFlagsID,
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

// The title for the config
func (cfg *PrometheusConfig) GetTitle() string {
	return "Prometheus"
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

// Get the sections underneath this one
func (cfg *PrometheusConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

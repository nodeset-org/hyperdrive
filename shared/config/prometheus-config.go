package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

const (
	// Tags
	prometheusTag string = "prom/prometheus:v2.49.1"
)

// Configuration for Prometheus
type PrometheusConfig struct {
	// The port to serve metrics on
	Port Parameter[uint16]

	// Toggle for forwarding the API port outside of Docker
	OpenPort Parameter[RpcPortMode]

	// The Docker Hub tag for Prometheus
	ContainerTag Parameter[string]

	// Custom command line flags
	AdditionalFlags Parameter[string]

	// Internal Fields
	parent *MetricsConfig
}

// Generates a new Prometheus config
func NewPrometheusConfig(parent *MetricsConfig) *PrometheusConfig {
	return &PrometheusConfig{
		parent: parent,

		Port: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.PortID,
				Name:               "Prometheus Port",
				Description:        "The port Prometheus should make its statistics available on.",
				AffectsContainers:  []ContainerID{ContainerID_Prometheus},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 9091,
			},
		},

		OpenPort: Parameter[RpcPortMode]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.OpenPortID,
				Name:               "Expose Prometheus Port",
				Description:        "Expose the Prometheus's port to other processes on your machine, or to your local network so other machines can access it too.",
				AffectsContainers:  []ContainerID{ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Options: getPortModes(""),
			Default: map[Network]RpcPortMode{
				Network_All: RpcPortMode_Closed,
			},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Prometheus Container Tag",
				Description:        "The tag name of the Prometheus container on Docker Hub you want to use.",
				AffectsContainers:  []ContainerID{ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_All: prometheusTag,
			},
		},

		AdditionalFlags: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Prometheus Flags",
				Description:        "Additional custom command line flags you want to pass to Prometheus, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
				AffectsContainers:  []ContainerID{ContainerID_Grafana},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},
	}
}

// The title for the config
func (cfg *PrometheusConfig) GetTitle() string {
	return "Prometheus"
}

// Get the parameters for this config
func (cfg *PrometheusConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.Port,
		&cfg.OpenPort,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *PrometheusConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
)

// Constants
const (
	// Tags
	grafanaTag string = "grafana/grafana:9.4.17"
)

// Configuration for Grafana
type GrafanaConfig struct {
	// The HTTP port to serve on
	Port Parameter[uint16]

	// The Docker Hub tag for Grafana
	ContainerTag Parameter[string]

	// Internal Fields
	parent *MetricsConfig
}

// Generates a new Grafana config
func NewGrafanaConfig(parent *MetricsConfig) *GrafanaConfig {
	return &GrafanaConfig{
		parent: parent,

		Port: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.PortID,
				Name:               "Grafana Port",
				Description:        "The port Grafana should run its HTTP server on - this is the port you will connect to in your browser.",
				AffectsContainers:  []ContainerID{ContainerID_Grafana},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 3100,
			},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Grafana Container Tag",
				Description:        "The tag name of the Grafana container you want to use on Docker Hub.",
				AffectsContainers:  []ContainerID{ContainerID_Grafana},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_All: grafanaTag,
			},
		},
	}
}

// The title for the config
func (cfg *GrafanaConfig) GetTitle() string {
	return "Grafana"
}

// Get the parameters for this config
func (cfg *GrafanaConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.Port,
		&cfg.ContainerTag,
	}
}

// Get the sections underneath this one
func (cfg *GrafanaConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

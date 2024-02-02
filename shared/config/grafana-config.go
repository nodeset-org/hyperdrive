package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Constants
const (
	// Tags
	grafanaTag string = "grafana/grafana:9.4.15"
)

// Configuration for Grafana
type GrafanaConfig struct {
	// The HTTP port to serve on
	Port types.Parameter[uint16]

	// The Docker Hub tag for Grafana
	ContainerTag types.Parameter[string]

	// Internal Fields
	parent *MetricsConfig
}

// Generates a new Grafana config
func NewGrafanaConfig(parent *MetricsConfig) *GrafanaConfig {
	return &GrafanaConfig{
		parent: parent,

		Port: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.PortID,
				Name:               "Grafana Port",
				Description:        "The port Grafana should run its HTTP server on - this is the port you will connect to in your browser.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Grafana},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 3100,
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Grafana Container Tag",
				Description:        "The tag name of the Grafana container you want to use on Docker Hub.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Grafana},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_All: grafanaTag,
			},
		},
	}
}

// The title for the config
func (cfg *GrafanaConfig) GetTitle() string {
	return "Grafana Settings"
}

// Get the parameters for this config
func (cfg *GrafanaConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Port,
		&cfg.ContainerTag,
	}
}

// Get the sections underneath this one
func (cfg *GrafanaConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

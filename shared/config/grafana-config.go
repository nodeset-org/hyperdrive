package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Constants
const (
	// Tags
	grafanaTag string = "grafana/grafana:9.4.15"
)

// Defaults
const defaultGrafanaPort uint16 = 3100

// Configuration for Grafana
type GrafanaConfig struct {
	Title string

	// The HTTP port to serve on
	Port types.Parameter[uint16]

	// The Docker Hub tag for Grafana
	ContainerTag types.Parameter[string]
}

// Generates a new Grafana config
func NewGrafanaConfig(cfg *HyperdriveConfig) *GrafanaConfig {
	return &GrafanaConfig{
		Title: "Grafana Settings",

		Port: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 PortID,
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
				ID:                 ContainerTagID,
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

// Get the parameters for this config
func (cfg *GrafanaConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.Port,
		&cfg.ContainerTag,
	}
}

// The the title for the config
func (cfg *GrafanaConfig) GetConfigTitle() string {
	return cfg.Title
}

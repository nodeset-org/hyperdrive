package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

const (
	// Param IDs
	ExporterEnableRootFsID string = "enableRootFs"

	// Tags
	exporterTag string = "prom/node-exporter:v1.7.0"
)

// Configuration for Exporter
type ExporterConfig struct {
	// Toggle for enabling access to the root filesystem (for multiple disk usage metrics)
	RootFs types.Parameter[bool]

	// The Docker Hub tag for the Exporter
	ContainerTag types.Parameter[string]

	// Custom command line flags
	AdditionalFlags types.Parameter[string]

	// Internal Fields
	parent *MetricsConfig
}

// Generates a new Exporter config
func NewExporterConfig(parent *MetricsConfig) *ExporterConfig {
	return &ExporterConfig{
		parent: parent,

		RootFs: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ExporterEnableRootFsID,
				Name:               "Allow Root Filesystem Access",
				Description:        "Give Prometheus's Node Exporter permission to view your root filesystem instead of being limited to its own Docker container.\nThis is needed if you want the Grafana dashboard to report the used disk space of a second SSD.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Exporter},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: false,
			},
		},

		ContainerTag: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Exporter Container Tag",
				Description:        "The tag name of the Prometheus Node Exporter container on Docker Hub you want to use.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Exporter},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[types.Network]string{
				types.Network_All: exporterTag,
			},
		},

		AdditionalFlags: types.Parameter[string]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Exporter Flags",
				Description:        "Additional custom command line flags you want to pass to the Node Exporter, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *ExporterConfig) GetTitle() string {
	return "Node Exporter"
}

// Get the parameters for this config
func (cfg *ExporterConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.RootFs,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *ExporterConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{}
}

package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
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
	RootFs Parameter[bool]

	// The Docker Hub tag for the Exporter
	ContainerTag Parameter[string]

	// Custom command line flags
	AdditionalFlags Parameter[string]

	// Internal Fields
	parent *MetricsConfig
}

// Generates a new Exporter config
func NewExporterConfig(parent *MetricsConfig) *ExporterConfig {
	return &ExporterConfig{
		parent: parent,

		RootFs: Parameter[bool]{
			ParameterCommon: &ParameterCommon{
				ID:                 ExporterEnableRootFsID,
				Name:               "Allow Root Filesystem Access",
				Description:        "Give Prometheus's Node Exporter permission to view your root filesystem instead of being limited to its own Docker container.\nThis is needed if you want the Grafana dashboard to report the used disk space of a second SSD.",
				AffectsContainers:  []ContainerID{ContainerID_Exporter},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]bool{
				Network_All: false,
			},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Exporter Container Tag",
				Description:        "The tag name of the Prometheus Node Exporter container on Docker Hub you want to use.",
				AffectsContainers:  []ContainerID{ContainerID_Exporter},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_All: exporterTag,
			},
		},

		AdditionalFlags: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Exporter Flags",
				Description:        "Additional custom command line flags you want to pass to the Node Exporter, to take advantage of other settings that Hyperdrive's configuration doesn't cover.",
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
func (cfg *ExporterConfig) GetTitle() string {
	return "Node Exporter"
}

// Get the parameters for this config
func (cfg *ExporterConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.RootFs,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *ExporterConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

package config

import (
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Constants
const (
	// Param IDs
	MetricsEnableID       string = "enableMetrics"
	MetricsEnableBitflyID string = "enableBitflyNodeMetrics"
	MetricsEcPortID       string = "ecMetricsPort"
	MetricsBnPortID       string = "bnMetricsPort"
	MetricsDaemonPortID   string = "daemonMetricsPort"
	MetricsExporterPortID string = "exporterMetricsPort"
)

// Configuration for Metrics
type MetricsConfig struct {
	EnableMetrics           types.Parameter[bool]
	EcMetricsPort           types.Parameter[uint16]
	BnMetricsPort           types.Parameter[uint16]
	DaemonMetricsPort       types.Parameter[uint16]
	ExporterMetricsPort     types.Parameter[uint16]
	EnableBitflyNodeMetrics types.Parameter[bool]

	// Subconfigs
	Grafana           *GrafanaConfig
	Prometheus        *PrometheusConfig
	Exporter          *ExporterConfig
	BitflyNodeMetrics *BitflyNodeMetricsConfig

	// Internal Fields
	parent *HyperdriveConfig
}

// Generates a new Besu configuration
func NewMetricsConfig(parent *HyperdriveConfig) *MetricsConfig {
	cfg := &MetricsConfig{
		parent: parent,

		EnableMetrics: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MetricsEnableID,
				Name:               "Enable Metrics",
				Description:        "Enable Hyperdrive's performance and status metrics system. This will provide you with the node operator's Grafana dashboard.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_ExecutionClient, types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients, types.ContainerID_Grafana, types.ContainerID_Prometheus, types.ContainerID_Exporter},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: true,
			},
		},

		EnableBitflyNodeMetrics: types.Parameter[bool]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MetricsEnableBitflyID,
				Name:               "Enable Beaconcha.in Node Metrics",
				Description:        "Enable the Beaconcha.in node metrics integration. This will allow you to track your node's metrics from your phone using the Beaconcha.in App.\n\nFor more information on setting up an account and the app, please visit https://beaconcha.in/mobile.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode, types.ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]bool{
				types.Network_All: false,
			},
		},

		EcMetricsPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MetricsEcPortID,
				Name:               "Execution Client Metrics Port",
				Description:        "The port your Execution client should expose its metrics on.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_ExecutionClient, types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 9105,
			},
		},

		BnMetricsPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MetricsBnPortID,
				Name:               "Beacon Node Metrics Port",
				Description:        "The port your Beacon Node's Beacon Node should expose its metrics on.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_BeaconNode, types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 9100,
			},
		},

		DaemonMetricsPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MetricsDaemonPortID,
				Name:               "Daemon Metrics Port",
				Description:        "The port your daemon container should expose its metrics on.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Daemon, types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 9102,
			},
		},

		ExporterMetricsPort: types.Parameter[uint16]{
			ParameterCommon: &types.ParameterCommon{
				ID:                 MetricsExporterPortID,
				Name:               "Exporter Metrics Port",
				Description:        "The port that Prometheus's Node Exporter should expose its metrics on.",
				AffectsContainers:  []types.ContainerID{types.ContainerID_Exporter, types.ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[types.Network]uint16{
				types.Network_All: 9103,
			},
		},
	}

	cfg.Grafana = NewGrafanaConfig(cfg)
	cfg.Prometheus = NewPrometheusConfig(cfg)
	cfg.Exporter = NewExporterConfig(cfg)
	cfg.BitflyNodeMetrics = NewBitflyNodeMetricsConfig(cfg)

	return cfg
}

// The title for the config
func (cfg *MetricsConfig) GetTitle() string {
	return "Metrics"
}

// Get the parameters for this config
func (cfg *MetricsConfig) GetParameters() []types.IParameter {
	return []types.IParameter{
		&cfg.EnableMetrics,
		&cfg.EcMetricsPort,
		&cfg.BnMetricsPort,
		&cfg.DaemonMetricsPort,
		&cfg.ExporterMetricsPort,
		&cfg.EnableBitflyNodeMetrics,
	}
}

// Get the sections underneath this one
func (cfg *MetricsConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{
		"grafana":    cfg.Grafana,
		"prometheus": cfg.Prometheus,
		"exporter":   cfg.Exporter,
		"bitfly":     cfg.BitflyNodeMetrics,
	}
}

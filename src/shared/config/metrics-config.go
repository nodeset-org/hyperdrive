package config

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
	EnableMetrics           Parameter[bool]
	EcMetricsPort           Parameter[uint16]
	BnMetricsPort           Parameter[uint16]
	DaemonMetricsPort       Parameter[uint16]
	ExporterMetricsPort     Parameter[uint16]
	EnableBitflyNodeMetrics Parameter[bool]

	// Subconfigs
	Grafana           *GrafanaConfig
	Prometheus        *PrometheusConfig
	Exporter          *ExporterConfig
	BitflyNodeMetrics *BitflyNodeMetricsConfig
}

// Generates a new Besu configuration
func NewMetricsConfig() *MetricsConfig {
	cfg := &MetricsConfig{
		EnableMetrics: Parameter[bool]{
			ParameterCommon: &ParameterCommon{
				ID:                 MetricsEnableID,
				Name:               "Enable Metrics",
				Description:        "Enable Hyperdrive's performance and status metrics system. This will provide you with the node operator's Grafana dashboard.",
				AffectsContainers:  []ContainerID{ContainerID_Daemon, ContainerID_ExecutionClient, ContainerID_BeaconNode, ContainerID_ValidatorClients, ContainerID_Grafana, ContainerID_Prometheus, ContainerID_Exporter},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]bool{
				Network_All: true,
			},
		},

		EnableBitflyNodeMetrics: Parameter[bool]{
			ParameterCommon: &ParameterCommon{
				ID:                 MetricsEnableBitflyID,
				Name:               "Enable Beaconcha.in Node Metrics",
				Description:        "Enable the Beaconcha.in node metrics integration. This will allow you to track your node's metrics from your phone using the Beaconcha.in App.\n\nFor more information on setting up an account and the app, please visit https://beaconcha.in/mobile.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode, ContainerID_ValidatorClients},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]bool{
				Network_All: false,
			},
		},

		EcMetricsPort: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 MetricsEcPortID,
				Name:               "Execution Client Metrics Port",
				Description:        "The port your Execution client should expose its metrics on.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient, ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 9105,
			},
		},

		BnMetricsPort: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 MetricsBnPortID,
				Name:               "Beacon Node Metrics Port",
				Description:        "The port your Beacon Node's Beacon Node should expose its metrics on.",
				AffectsContainers:  []ContainerID{ContainerID_BeaconNode, ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 9100,
			},
		},

		DaemonMetricsPort: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 MetricsDaemonPortID,
				Name:               "Daemon Metrics Port",
				Description:        "The port your daemon container should expose its metrics on.",
				AffectsContainers:  []ContainerID{ContainerID_Daemon, ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 9102,
			},
		},

		ExporterMetricsPort: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 MetricsExporterPortID,
				Name:               "Exporter Metrics Port",
				Description:        "The port that Prometheus's Node Exporter should expose its metrics on.",
				AffectsContainers:  []ContainerID{ContainerID_Exporter, ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 9103,
			},
		},
	}

	cfg.Grafana = NewGrafanaConfig()
	cfg.Prometheus = NewPrometheusConfig()
	cfg.Exporter = NewExporterConfig()
	cfg.BitflyNodeMetrics = NewBitflyNodeMetricsConfig()

	return cfg
}

// The title for the config
func (cfg *MetricsConfig) GetTitle() string {
	return "Metrics"
}

// Get the parameters for this config
func (cfg *MetricsConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.EnableMetrics,
		&cfg.EcMetricsPort,
		&cfg.BnMetricsPort,
		&cfg.DaemonMetricsPort,
		&cfg.ExporterMetricsPort,
		&cfg.EnableBitflyNodeMetrics,
	}
}

// Get the sections underneath this one
func (cfg *MetricsConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{
		"grafana":    cfg.Grafana,
		"prometheus": cfg.Prometheus,
		"exporter":   cfg.Exporter,
		"bitfly":     cfg.BitflyNodeMetrics,
	}
}

package config

// The master configuration struct
type RocketPoolConfig struct {
	Title string `yaml:"-"`

	Version string `yaml:"-"`

	RocketPoolDirectory string `yaml:"-"`

	IsNativeMode bool `yaml:"-"`

	// Execution client settings
	ExecutionClientMode Parameter `yaml:"executionClientMode,omitempty"`
	ExecutionClient     Parameter `yaml:"executionClient,omitempty"`

	// Fallback settings
	UseFallbackClients Parameter `yaml:"useFallbackClients,omitempty"`
	ReconnectDelay     Parameter `yaml:"reconnectDelay,omitempty"`

	// Consensus client settings
	ConsensusClientMode     Parameter `yaml:"consensusClientMode,omitempty"`
	ConsensusClient         Parameter `yaml:"consensusClient,omitempty"`
	ExternalConsensusClient Parameter `yaml:"externalConsensusClient,omitempty"`

	// Metrics settings
	EnableMetrics           Parameter `yaml:"enableMetrics,omitempty"`
	EnableODaoMetrics       Parameter `yaml:"enableODaoMetrics,omitempty"`
	EcMetricsPort           Parameter `yaml:"ecMetricsPort,omitempty"`
	BnMetricsPort           Parameter `yaml:"bnMetricsPort,omitempty"`
	VcMetricsPort           Parameter `yaml:"vcMetricsPort,omitempty"`
	NodeMetricsPort         Parameter `yaml:"nodeMetricsPort,omitempty"`
	ExporterMetricsPort     Parameter `yaml:"exporterMetricsPort,omitempty"`
	WatchtowerMetricsPort   Parameter `yaml:"watchtowerMetricsPort,omitempty"`
	EnableBitflyNodeMetrics Parameter `yaml:"enableBitflyNodeMetrics,omitempty"`

	// The Smartnode configuration
	Smartnode *SmartnodeConfig `yaml:"smartnode,omitempty"`

	// TODO: Import/Uncomment the rest later on as needed bases
	// Execution client configurations
	// ExecutionCommon   *ExecutionCommonConfig   `yaml:"executionCommon,omitempty"`
	// Geth              *GethConfig              `yaml:"geth,omitempty"`
	// Nethermind        *NethermindConfig        `yaml:"nethermind,omitempty"`
	// Besu              *BesuConfig              `yaml:"besu,omitempty"`
	// ExternalExecution *ExternalExecutionConfig `yaml:"externalExecution,omitempty"`

	// // Consensus client configurations
	// ConsensusCommon    *ConsensusCommonConfig    `yaml:"consensusCommon,omitempty"`
	// Lighthouse         *LighthouseConfig         `yaml:"lighthouse,omitempty"`
	// Lodestar           *LodestarConfig           `yaml:"lodestar,omitempty"`
	// Nimbus             *NimbusConfig             `yaml:"nimbus,omitempty"`
	// Prysm              *PrysmConfig              `yaml:"prysm,omitempty"`
	// Teku               *TekuConfig               `yaml:"teku,omitempty"`
	// ExternalLighthouse *ExternalLighthouseConfig `yaml:"externalLighthouse,omitempty"`
	// ExternalNimbus     *ExternalNimbusConfig     `yaml:"externalNimbus,omitempty"`
	// ExternalLodestar   *ExternalLodestarConfig   `yaml:"externalLodestar,omitempty"`
	// ExternalPrysm      *ExternalPrysmConfig      `yaml:"externalPrysm,omitempty"`
	// ExternalTeku       *ExternalTekuConfig       `yaml:"externalTeku,omitempty"`

	// // Fallback client configurations
	// FallbackNormal *FallbackNormalConfig `yaml:"fallbackNormal,omitempty"`
	// FallbackPrysm  *FallbackPrysmConfig  `yaml:"fallbackPrysm,omitempty"`

	// // Metrics
	// Grafana           *GrafanaConfig           `yaml:"grafana,omitempty"`
	// Prometheus        *PrometheusConfig        `yaml:"prometheus,omitempty"`
	// Exporter          *ExporterConfig          `yaml:"exporter,omitempty"`
	// BitflyNodeMetrics *BitflyNodeMetricsConfig `yaml:"bitflyNodeMetrics,omitempty"`

	// // Native mode
	// Native *NativeConfig `yaml:"native,omitempty"`

	// // MEV-Boost
	// EnableMevBoost config.Parameter `yaml:"enableMevBoost,omitempty"`
	// MevBoost       *MevBoostConfig  `yaml:"mevBoost,omitempty"`

	// // Addons
	// GraffitiWallWriter addontypes.SmartnodeAddon `yaml:"addon-gww,omitempty"`
}

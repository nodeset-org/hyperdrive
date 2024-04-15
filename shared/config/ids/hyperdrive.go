package ids

const (
	// Hyperdrive parameter IDs
	RootConfigID         string = "hyperdrive"
	VersionID            string = "version"
	NetworkID            string = "network"
	ClientModeID         string = "clientMode"
	UserDataPathID       string = "hdUserDataDir"
	ProjectNameID        string = "projectName"
	AutoTxMaxFeeID       string = "autoTxMaxFee"
	MaxPriorityFeeID     string = "maxPriorityFee"
	AutoTxGasThresholdID string = "autoTxGasThreshold"

	// Subconfig IDs
	LoggingID           string = "logging"
	FallbackID          string = "fallback"
	LocalExecutionID    string = "localExecution"
	ExternalExecutionID string = "externalExecution"
	LocalBeaconID       string = "localBeacon"
	ExternalBeaconID    string = "externalBeacon"
	MetricsID           string = "metrics"
)

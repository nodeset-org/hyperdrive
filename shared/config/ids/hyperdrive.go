package ids

const (
	// Hyperdrive parameter IDs
	RootConfigID               string = "hyperdrive"
	VersionID                  string = "version"
	UserDirID                  string = "hdUserDir"
	ApiPortID                  string = "apiPort"
	NetworkID                  string = "network"
	EnableIPv6ID               string = "enableIPv6"
	ClientModeID               string = "clientMode"
	UserDataPathID             string = "hdUserDataDir"
	ProjectNameID              string = "projectName"
	AutoTxMaxFeeID             string = "autoTxMaxFee"
	MaxPriorityFeeID           string = "maxPriorityFee"
	AutoTxGasThresholdID       string = "autoTxGasThreshold"
	AdditionalDockerNetworksID string = "additionalDockerNetworks"
	ContainerTagID             string = "containerTag"

	// Subconfig IDs
	LoggingID           string = "logging"
	FallbackID          string = "fallback"
	LocalExecutionID    string = "localExecution"
	ExternalExecutionID string = "externalExecution"
	LocalBeaconID       string = "localBeacon"
	ExternalBeaconID    string = "externalBeacon"
	MetricsID           string = "metrics"
	MevBoostID          string = "mevBoost"

	// MEV-Boost
	MevBoostEnableID             string = "enableMevBoost"
	MevBoostModeID               string = "mode"
	MevBoostSelectionModeID      string = "selectionMode"
	MevBoostOpenRpcPortID        string = "openRpcPort"
	MevBoostExternalUrlID        string = "externalUrl"
	MevBoostFlashbotsID          string = "flashbotsEnabled"
	MevBoostBloxRouteMaxProfitID string = "bloxRouteMaxProfitEnabled"
	MevBoostBloxRouteRegulatedID string = "bloxRouteRegulatedEnabled"
	MevBoostEdenID               string = "edenEnabled"
	MevBoostTitanRegionalID      string = "titanRegionaEnabled"
	MevBoostCustomRelaysID       string = "customRelays"
)

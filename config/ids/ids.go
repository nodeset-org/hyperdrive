package ids

const (
	// Base config
	ProjectNameID              string = "projectName"
	ApiPortID                  string = "apiPort"
	EnableIPv6ID               string = "enableIPv6"
	UserDataPathID             string = "userDataPath"
	AdditionalDockerNetworksID string = "additionalDockerNetworks"
	ClientTimeoutID            string = "clientTimeout"
	ContainerTagID             string = "containerTag"
	LoggingSectionID           string = "logging"

	// Logging
	LoggerLevelID      string = "level"
	LoggerFormatID     string = "format"
	LoggerAddSourceID  string = "addSource"
	LoggerMaxSizeID    string = "maxSize"
	LoggerMaxBackupsID string = "maxBackups"
	LoggerMaxAgeID     string = "maxAge"
	LoggerLocalTimeID  string = "localTime"
	LoggerCompressID   string = "compress"
)

package config

const (
	// The filename for Hyperdrive settings files
	ConfigFilename string = "user-settings.yml"

	// API base route for daemon requests
	HyperdriveDaemonRoute string = "hyperdrive"

	// API version for daemon requests
	HyperdriveApiVersion string = "1"

	// Complete API route for client requests
	HyperdriveApiClientRoute string = HyperdriveDaemonRoute + "/api/v" + HyperdriveApiVersion
)

package config

const (
	EventLogInterval         int    = 1000
	DockerApiVersion         string = "1.40"
	HyperdriveDaemonRoute    string = "hyperdrive"
	HyperdriveSocketFilename string = HyperdriveDaemonRoute + ".sock"
	ConfigFilename           string = "user-settings.yml"

	// Wallet
	UserAddressFilename    string = "address"
	UserWalletDataFilename string = "wallet"
	UserPasswordFilename   string = "password"
)

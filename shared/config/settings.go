package config

const (
	EventLogInterval      int    = 1000
	DockerApiVersion      string = "1.40"
	SocketFilename        string = "daemon.sock"
	HyperdriveDaemonRoute string = "hyperdrive"

	// Wallet
	UserAddressFilename      string = "address"
	UserWalletDataFilename   string = "wallet"
	UserLegacyWalletFilename string = "wallet-v3"
	UserPasswordFilename     string = "password"

	// Modules
	ModulesDir string = "modules"
)

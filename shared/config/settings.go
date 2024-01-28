package config

const (
	EventLogInterval int    = 1000
	DockerApiVersion string = "1.40"
	SocketFilename   string = "daemon.sock"

	// Wallet
	UserAddressFilename      = "address"
	UserWalletDataFilename   = "wallet.json"
	UserLegacyWalletFilename = "wallet-v3"
	UserPasswordFilename     = "password"
)

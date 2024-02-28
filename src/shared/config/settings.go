package config

import "time"

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

	// Scripts
	EcStartScript string = "start-ec.sh"
	BnStartScript string = "start-bn.sh"
	VcStartScript string = "start-vc.sh"

	// HTTP
	ClientTimeout time.Duration = 8 * time.Second

	// Volumes
	ExecutionClientDataVolume string = "ecdata"
	BeaconNodeDataVolume      string = "bndata"
)

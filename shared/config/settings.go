package config

import "time"

const (
	EventLogInterval            int    = 1000
	HyperdriveDaemonRoute       string = "hyperdrive"
	HyperdriveApiVersion        string = "1"
	HyperdriveApiClientRoute    string = HyperdriveDaemonRoute + "/api/v" + HyperdriveApiVersion
	HyperdriveCliSocketFilename string = HyperdriveDaemonRoute + "-cli.sock"
	HyperdriveNetSocketFilename string = HyperdriveDaemonRoute + "-net.sock"
	ConfigFilename              string = "user-settings.yml"

	// Wallet
	UserAddressFilename    string = "address"
	UserWalletDataFilename string = "wallet"
	UserPasswordFilename   string = "password"

	// Scripts
	EcStartScript string = "start-ec.sh"
	BnStartScript string = "start-bn.sh"
	VcStartScript string = "start-vc.sh"

	// HTTP
	ClientTimeout time.Duration = 1 * time.Minute

	// Volumes
	ExecutionClientDataVolume string = "ecdata"
	BeaconNodeDataVolume      string = "bndata"

	// Logging
	LogDir       string = "logs"
	ApiLogName   string = "api.log"
	TasksLogName string = "tasks.log"
)

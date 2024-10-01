package config

const (
	EventLogInterval         int    = 1000
	HyperdriveDaemonRoute    string = "hyperdrive"
	HyperdriveApiVersion     string = "1"
	HyperdriveApiClientRoute string = HyperdriveDaemonRoute + "/api/v" + HyperdriveApiVersion
	ConfigFilename           string = "user-settings.yml"
	DefaultApiPort           uint16 = 8080

	// Wallet
	UserAddressFilename    string = "address"
	UserWalletDataFilename string = "wallet"
	UserPasswordFilename   string = "password"

	// Scripts
	EcStartScript       string = "start-ec.sh"
	BnStartScript       string = "start-bn.sh"
	VcStartScript       string = "start-vc.sh"
	MevBoostStartScript string = "start-mev-boost.sh"

	// Volumes
	ExecutionClientDataVolume string = "ecdata"
	BeaconNodeDataVolume      string = "bndata"

	// Logging
	LogDir       string = "logs"
	ApiLogName   string = "api.log"
	TasksLogName string = "tasks.log"
)

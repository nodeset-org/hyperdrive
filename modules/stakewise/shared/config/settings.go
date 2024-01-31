package swconfig

const (
	DaemonRoute          string = "stakewise"
	SocketFilename       string = DaemonRoute + ".sock"
	WalletFilename       string = "wallet.json"
	PasswordFilename     string = "password.txt"
	KeystorePasswordFile string = "secret.txt"
	DepositDataFile      string = "deposit-data.json"

	// Container settings
	DaemonContainerSuffix   string = "sw_daemon"
	OperatorContainerSuffix string = "sw_operator"
	VcContainerSuffix       string = "sw_vc"
)

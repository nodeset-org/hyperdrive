package swconfig

const (
	ModuleName           string = "stakewise"
	DaemonBaseRoute      string = ModuleName
	ApiVersion           string = "1"
	ApiClientRoute       string = DaemonBaseRoute + "/api/v" + ApiVersion
	CliSocketFilename    string = ModuleName + "-cli.sock"
	NetSocketFilename    string = ModuleName + "-net.sock"
	WalletFilename       string = "wallet.json"
	PasswordFilename     string = "password.txt"
	KeystorePasswordFile string = "secret.txt"
	DepositDataFile      string = "deposit-data.json"
	ClientLogName        string = "hd-client.log"
)

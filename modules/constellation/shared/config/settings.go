package constconfig

const (
	ModuleName           string = "constellation"
	ShortModuleName      string = "cs"
	DaemonBaseRoute      string = ModuleName
	ApiVersion           string = "1"
	CliSocketFilename    string = ModuleName + "-cli.sock"
	NetSocketFilename    string = ModuleName + "-net.sock"
	SocketFilename       string = ModuleName + ".sock"
	WalletFilename       string = "wallet.json"
	PasswordFilename     string = "password.txt"
	KeystorePasswordFile string = "secret.txt"
	DepositDataFile      string = "deposit-data.json"
	ClientLogName        string = "hd-client.log"
)

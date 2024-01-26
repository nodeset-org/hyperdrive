package config

const (
	EventLogInterval int    = 1000
	DockerApiVersion string = "1.40"
	SocketFilename   string = "daemon.sock"

	// Daemon settings
	DaemonUserDirPath     string = "/hyperdrive/user/"
	DaemonConfigPath      string = DaemonUserDirPath + ConfigFilename
	DaemonSocketPath      string = DaemonUserDirPath + SocketFilename
	DaemonUserDataPath    string = "/hyperdrive/user/data"
	DaemonScriptsPath     string = "/hyperdrive/system/scripts"
	DaemonGlobalDataPath  string = "/hyperdrive/system/global"
	DaemonProjectDataPath string = "/hyperdrive/system/data"
)

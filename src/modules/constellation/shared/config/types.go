package constconfig

import (
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	// The Constellation Hyperdrive daemon
	ContainerID_ConstellationDaemon config.ContainerID = "constellation_daemon"
	// The constellation Validator client
	ContainerID_ConstellationValidator config.ContainerID = "constellation_vc"
)

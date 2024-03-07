package swconfig

import (
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	// The stakewise Hyperdrive daemon
	ContainerID_StakewiseDaemon config.ContainerID = "sw_daemon"

	// The stakewise operator container
	ContainerID_StakewiseOperator config.ContainerID = "sw_operator"

	// The stakewise Validator client
	ContainerID_StakewiseValidator config.ContainerID = "sw_vc"
)

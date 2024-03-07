package swconfig

import (
	nmc_config "github.com/rocket-pool/node-manager-core/config"
)

const (
	// The stakewise Hyperdrive daemon
	ContainerID_StakewiseDaemon nmc_config.ContainerID = "sw_daemon"

	// The stakewise operator container
	ContainerID_StakewiseOperator nmc_config.ContainerID = "sw_operator"

	// The stakewise Validator client
	ContainerID_StakewiseValidator nmc_config.ContainerID = "sw_vc"
)

package stakewise

import "github.com/nodeset-org/hyperdrive/shared/types"

const (
	// The stakewise Hyperdrive daemon
	ContainerID_StakewiseDaemon types.ContainerID = "sw-daemon"

	// The stakewise operator container
	ContainerID_StakewiseOperator types.ContainerID = "sw-operator"

	// The stakewise Validator client
	ContainerID_StakewiseValidator types.ContainerID = "sw-vc"
)

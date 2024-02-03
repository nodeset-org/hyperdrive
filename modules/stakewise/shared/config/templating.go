package swconfig

import (
	"fmt"

	swshared "github.com/nodeset-org/hyperdrive/modules/stakewise/shared"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

func (c *StakewiseConfig) WalletFilename() string {
	return WalletFilename
}

func (c *StakewiseConfig) PasswordFilename() string {
	return PasswordFilename
}

func (c *StakewiseConfig) KeystorePasswordFile() string {
	return KeystorePasswordFile
}

func (c *StakewiseConfig) DaemonContainerName() string {
	return string(ContainerID_StakewiseDaemon)
}

func (c *StakewiseConfig) OperatorContainerName() string {
	return string(ContainerID_StakewiseOperator)
}

func (c *StakewiseConfig) VcContainerName() string {
	return string(ContainerID_StakewiseValidator)
}

func (c *StakewiseConfig) DepositDataFile() string {
	return DepositDataFile
}

// The tag for the daemon container
func (cfg *StakewiseConfig) DaemonTag() string {
	return daemonTag
}

// Get the container tag of the selected VC
func (cfg *StakewiseConfig) GetVcContainerTag() string {
	bn := cfg.hdCfg.GetSelectedBeaconNode()
	switch bn {
	case types.BeaconNode_Lighthouse:
		return cfg.Lighthouse.ContainerTag.Value
	case types.BeaconNode_Lodestar:
		return cfg.Lodestar.ContainerTag.Value
	case types.BeaconNode_Nimbus:
		return cfg.Nimbus.ContainerTag.Value
	case types.BeaconNode_Prysm:
		return cfg.Prysm.ContainerTag.Value
	case types.BeaconNode_Teku:
		return cfg.Teku.ContainerTag.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", bn))
	}
}

// Gets the additional flags of the selected VC
func (cfg *StakewiseConfig) GetVcAdditionalFlags() string {
	bn := cfg.hdCfg.GetSelectedBeaconNode()
	switch bn {
	case types.BeaconNode_Lighthouse:
		return cfg.Lighthouse.AdditionalFlags.Value
	case types.BeaconNode_Lodestar:
		return cfg.Lodestar.AdditionalFlags.Value
	case types.BeaconNode_Nimbus:
		return cfg.Nimbus.AdditionalFlags.Value
	case types.BeaconNode_Prysm:
		return cfg.Prysm.AdditionalFlags.Value
	case types.BeaconNode_Teku:
		return cfg.Teku.AdditionalFlags.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", bn))
	}
}

// Check if any of the services have doppelganger detection enabled
// NOTE: update this with each new service that runs a VC!
func (cfg *StakewiseConfig) IsDoppelgangerEnabled() bool {
	return cfg.VcCommon.DoppelgangerDetection.Value
}

// Used by text/template to format validator.yml
func (cfg *StakewiseConfig) Graffiti() (string, error) {
	prefix := cfg.hdCfg.GraffitiPrefix()
	customGraffiti := cfg.VcCommon.Graffiti.Value
	if customGraffiti == "" {
		return prefix, nil
	}
	return fmt.Sprintf("%s (%s)", prefix, customGraffiti), nil
}

func (cfg *StakewiseConfig) FeeRecipient() string {
	res := swshared.NewStakewiseResources(cfg.hdCfg.Network.Value)
	return res.FeeRecipient.Hex()
}

func (cfg *StakewiseConfig) Vault() string {
	res := swshared.NewStakewiseResources(cfg.hdCfg.Network.Value)
	return res.Vault.Hex()
}

func (cfg *StakewiseConfig) Network() string {
	res := swshared.NewStakewiseResources(cfg.hdCfg.Network.Value)
	return res.NodesetNetwork
}

func (cfg *StakewiseConfig) IsEnabled() bool {
	return cfg.Enabled.Value
}

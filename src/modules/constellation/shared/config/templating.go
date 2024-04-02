package constconfig

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/shared/config"
)

func (c *ConstellationConfig) WalletFilename() string {
	return WalletFilename
}

func (c *ConstellationConfig) PasswordFilename() string {
	return PasswordFilename
}

func (c *ConstellationConfig) KeystorePasswordFile() string {
	return KeystorePasswordFile
}

func (c *ConstellationConfig) DaemonContainerName() string {
	return string(ContainerID_ConstellationDaemon)
}

func (c *ConstellationConfig) OperatorContainerName() string {
	return string(ContainerID_ConstellationOperator)
}

func (c *ConstellationConfig) VcContainerName() string {
	return string(ContainerID_ConstellationValidator)
}

func (c *ConstellationConfig) DepositDataFile() string {
	return DepositDataFile
}

// The tag for the daemon container
func (cfg *ConstellationConfig) DaemonTag() string {
	return daemonTag
}

// Get the container tag of the selected VC
func (cfg *ConstellationConfig) GetVcContainerTag() string {
	bn := cfg.hdCfg.GetSelectedBeaconNode()
	switch bn {
	case config.BeaconNode_Lighthouse:
		return cfg.Lighthouse.ContainerTag.Value
	case config.BeaconNode_Lodestar:
		return cfg.Lodestar.ContainerTag.Value
	case config.BeaconNode_Nimbus:
		return cfg.Nimbus.ContainerTag.Value
	case config.BeaconNode_Prysm:
		return cfg.Prysm.ContainerTag.Value
	case config.BeaconNode_Teku:
		return cfg.Teku.ContainerTag.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", bn))
	}
}

// Gets the additional flags of the selected VC
func (cfg *ConstellationConfig) GetVcAdditionalFlags() string {
	bn := cfg.hdCfg.GetSelectedBeaconNode()
	switch bn {
	case config.BeaconNode_Lighthouse:
		return cfg.Lighthouse.AdditionalFlags.Value
	case config.BeaconNode_Lodestar:
		return cfg.Lodestar.AdditionalFlags.Value
	case config.BeaconNode_Nimbus:
		return cfg.Nimbus.AdditionalFlags.Value
	case config.BeaconNode_Prysm:
		return cfg.Prysm.AdditionalFlags.Value
	case config.BeaconNode_Teku:
		return cfg.Teku.AdditionalFlags.Value
	default:
		panic(fmt.Sprintf("Unknown Beacon Node %s", bn))
	}
}

// Check if any of the services have doppelganger detection enabled
// NOTE: update this with each new service that runs a VC!
func (cfg *ConstellationConfig) IsDoppelgangerEnabled() bool {
	return cfg.VcCommon.DoppelgangerDetection.Value
}

// Used by text/template to format validator.yml
func (cfg *ConstellationConfig) Graffiti() (string, error) {
	prefix := cfg.hdCfg.GraffitiPrefix()
	customGraffiti := cfg.VcCommon.Graffiti.Value
	if customGraffiti == "" {
		return prefix, nil
	}
	return fmt.Sprintf("%s (%s)", prefix, customGraffiti), nil
}

func (cfg *ConstellationConfig) FeeRecipient() string {
	res := swshared.NewConstellationResources(cfg.hdCfg.Network.Value)
	return res.FeeRecipient.Hex()
}

func (cfg *ConstellationConfig) Vault() string {
	res := swshared.NewConstellationResources(cfg.hdCfg.Network.Value)
	return res.Vault.Hex()
}

func (cfg *ConstellationConfig) Network() string {
	res := swshared.NewConstellationResources(cfg.hdCfg.Network.Value)
	return res.NodesetNetwork
}

func (cfg *ConstellationConfig) IsEnabled() bool {
	return cfg.Enabled.Value
}

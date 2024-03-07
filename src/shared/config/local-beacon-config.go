package config

import (
	nmc_config "github.com/rocket-pool/node-manager-core/config"
)

// Create a new LocalBeaconConfig struct
func NewLocalBeaconConfig() *nmc_config.LocalBeaconConfig {
	cfg := nmc_config.NewLocalBeaconConfig()
	cfg.Lighthouse.ContainerTag.Default[Network_HoleskyDev] = cfg.Lighthouse.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Lodestar.ContainerTag.Default[Network_HoleskyDev] = cfg.Lodestar.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Nimbus.ContainerTag.Default[Network_HoleskyDev] = cfg.Nimbus.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Prysm.ContainerTag.Default[Network_HoleskyDev] = cfg.Prysm.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Teku.ContainerTag.Default[Network_HoleskyDev] = cfg.Teku.ContainerTag.Default[nmc_config.Network_Holesky]
	return cfg
}

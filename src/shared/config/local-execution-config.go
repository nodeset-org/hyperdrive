package config

import (
	nmc_config "github.com/rocket-pool/node-manager-core/config"
)

// Create a new LocalExecutionConfig struct
func NewLocalExecutionConfig() *nmc_config.LocalExecutionConfig {
	cfg := nmc_config.NewLocalExecutionConfig()
	cfg.Besu.ContainerTag.Default[Network_HoleskyDev] = cfg.Besu.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Geth.ContainerTag.Default[Network_HoleskyDev] = cfg.Geth.ContainerTag.Default[nmc_config.Network_Holesky]
	cfg.Nethermind.ContainerTag.Default[Network_HoleskyDev] = cfg.Nethermind.ContainerTag.Default[nmc_config.Network_Holesky]
	return cfg
}

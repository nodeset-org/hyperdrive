package config

import (
	"github.com/rocket-pool/node-manager-core/config"
)

// Create a new LocalExecutionClient struct
func NewLocalExecutionClient() *config.LocalExecutionConfig {
	cfg := config.NewLocalExecutionConfig()
	cfg.Besu.ContainerTag.Default[Network_HoleskyDev] = cfg.Besu.ContainerTag.Default[config.Network_Holesky]
	cfg.Geth.ContainerTag.Default[Network_HoleskyDev] = cfg.Geth.ContainerTag.Default[config.Network_Holesky]
	cfg.Nethermind.ContainerTag.Default[Network_HoleskyDev] = cfg.Nethermind.ContainerTag.Default[config.Network_Holesky]
	cfg.Nethermind.FullPruningThresholdMb.Default[Network_HoleskyDev] = cfg.Nethermind.FullPruningThresholdMb.Default[config.Network_Holesky]
	cfg.Reth.ContainerTag.Default[Network_HoleskyDev] = cfg.Reth.ContainerTag.Default[config.Network_Holesky]
	return cfg
}

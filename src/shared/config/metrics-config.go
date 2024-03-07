package config

import (
	nmc_config "github.com/rocket-pool/node-manager-core/config"
)

// Generates a new metrics config
func NewMetricsConfig() *nmc_config.MetricsConfig {
	cfg := nmc_config.NewMetricsConfig()
	cfg.BitflyNodeMetrics.MachineName.Default = map[nmc_config.Network]string{
		nmc_config.Network_All: "Hyperdrive Node",
	}
	return cfg
}

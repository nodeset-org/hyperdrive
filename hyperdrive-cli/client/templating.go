package client

import "github.com/nodeset-org/hyperdrive-daemon/shared/config"

// Get the configs for all of the modules in the system that are enabled
func (c *GlobalConfig) GetEnabledModuleConfigNames() []string {
	names := []string{}
	for _, cfg := range c.GetAllModuleConfigs() {
		if cfg.IsEnabled() {
			names = append(names, cfg.GetModuleName())
		}
	}
	return names
}

func (c *GlobalConfig) ModulesDirectory() string {
	return config.ModulesName
}

func (c *GlobalConfig) ValidatorsDirectory() string {
	return config.ValidatorsDirectory
}

func (c *GlobalConfig) HyperdriveApiKeyPath() string {
	return hdApiKeyRelPath
}

func (c *GlobalConfig) StakeWiseApiKeyPath() string {
	return swApiKeyRelPath
}

func (c *GlobalConfig) ConstellationApiKeyPath() string {
	return csApiKeyRelPath
}

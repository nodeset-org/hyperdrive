package client

import modconfig "github.com/nodeset-org/hyperdrive/shared/config/modules"

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
	return modconfig.ModulesName
}

func (c *GlobalConfig) ValidatorsDirectory() string {
	return modconfig.ValidatorsDirectory
}

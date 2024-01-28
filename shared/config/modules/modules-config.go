package modules

import (
	"github.com/nodeset-org/hyperdrive/shared/config/modules/stakewise"
	"github.com/nodeset-org/hyperdrive/shared/types"
)

// Configuration for Hyperdrive modules
type ModulesConfig struct {
	Stakewise *stakewise.StakewiseConfig
}

// Generates a new Modules config
func NewModulesConfig() *ModulesConfig {
	cfg := &ModulesConfig{}

	cfg.Stakewise = stakewise.NewStakewiseConfig()

	return cfg
}

// The the title for the config
func (cfg *ModulesConfig) GetTitle() string {
	return "Hyperdrive Modules"
}

// Get the parameters for this config
func (cfg *ModulesConfig) GetParameters() []types.IParameter {
	return []types.IParameter{}
}

// Get the sections underneath this one
func (cfg *ModulesConfig) GetSubconfigs() map[string]types.IConfigSection {
	return map[string]types.IConfigSection{
		"stakewise": cfg.Stakewise,
	}
}

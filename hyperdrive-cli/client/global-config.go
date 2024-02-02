package client

import (
	"fmt"
	"reflect"

	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config"
	modconfig "github.com/nodeset-org/hyperdrive/shared/config/modules"
)

// Wrapper for global configuration
type GlobalConfig struct {
	Hyperdrive *config.HyperdriveConfig
	Stakewise  *swconfig.StakewiseConfig
}

// Make a new global config
func NewGlobalConfig(hdCfg *config.HyperdriveConfig) *GlobalConfig {
	return &GlobalConfig{
		Hyperdrive: hdCfg,
		Stakewise:  swconfig.NewStakewiseConfig(hdCfg),
	}
}

// Serialize the config and all modules
func (c *GlobalConfig) Serialize() map[string]any {
	return c.Hyperdrive.Serialize([]modconfig.IModuleConfig{
		c.Stakewise,
	})
}

// Deserialize the config's modules (assumes the Hyperdrive config itself has already been deserialized)
func (c *GlobalConfig) DeserializeModules() error {
	// Load Stakewise
	stakewiseName := c.Stakewise.GetModuleName()
	section, exists := c.Hyperdrive.Modules[stakewiseName]
	if exists {
		configMap, ok := section.(map[string]any)
		if !ok {
			return fmt.Errorf("config module section [%s] is not a map, it's a %s", stakewiseName, reflect.TypeOf(section))
		}
		err := config.Deserialize(c.Stakewise, configMap, c.Hyperdrive.Network.Value)
		if err != nil {
			return fmt.Errorf("error deserializing stakewise configuration: %w", err)
		}
	}
	return nil
}

package constcommon

import (
	"fmt"
	"reflect"

	"github.com/nodeset-org/hyperdrive/daemon-utils/services"
	swconfig "github.com/nodeset-org/hyperdrive/modules/constellation/shared/config"
)

type ConstellationServiceProvider struct {
	*services.ServiceProvider
	swCfg     *swconfig.ConstellationConfig
	resources *swconfig.ConstellationResources
}

// Create a new service provider with Constellation daemon-specific features
func NewConstellationServiceProvider(sp *services.ServiceProvider) (*ConstellationServiceProvider, error) {
	// Create the resources
	cfg := sp.GetHyperdriveConfig()
	res := swconfig.NewConstellationResources(cfg.Network.Value)
	swCfg, ok := sp.GetModuleConfig().(*swconfig.ConstellationConfig)
	if !ok {
		return nil, fmt.Errorf("constellation config is not the correct type, it's a %s", reflect.TypeOf(swCfg))
	}

	// Make the provider
	constellationSp := &ConstellationServiceProvider{
		ServiceProvider: sp,
		swCfg:           swCfg,
		resources:       res,
	}

	return constellationSp, nil
}

func (s *ConstellationServiceProvider) GetModuleConfig() *swconfig.ConstellationConfig {
	return s.swCfg
}

func (s *ConstellationServiceProvider) GetResources() *swconfig.ConstellationResources {
	return s.resources
}

package constcommon

import (
	"fmt"
	"reflect"

	"github.com/nodeset-org/hyperdrive/daemon-utils/services"

	swconfig "github.com/nodeset-org/hyperdrive/modules/constellation/shared/config"
)

type ConstellationServiceProvider struct {
	ServiceProvider *services.ServiceProvider
	// RpServiceProvider *rpservices.ServiceProvider

	swCfg     *swconfig.ConstellationConfig
	resources *swconfig.ConstellationResources
}

// Create a new service provider with Constellation daemon-specific features
func NewConstellationServiceProvider(sp *services.ServiceProvider) (*ConstellationServiceProvider, error) {
	// Create the resources
	cfg := sp.GetHyperdriveConfig()
	fmt.Printf("!!! cfg: %v\n", cfg)
	res := swconfig.NewConstellationResources(cfg.Network.Value)
	fmt.Printf("!!! res: %v\n", res)
	swCfg, ok := sp.GetModuleConfig().(*swconfig.ConstellationConfig)
	fmt.Printf("!!! swCfg: %v\n", swCfg)
	if !ok {
		return nil, fmt.Errorf("constellation config is not the correct type, it's a %s", reflect.TypeOf(swCfg))
	}
	// rpServiceProvider, err := rpservices.NewServiceProvider(sp.GetUserDir())
	// fmt.Printf("!!! rpServiceProvider: %v\n", rpServiceProvider)
	// if err != nil {
	// 	return nil, err
	// }
	// Make the provider
	constellationSp := &ConstellationServiceProvider{
		ServiceProvider: sp,
		// RpServiceProvider: rpServiceProvider,
		swCfg:     swCfg,
		resources: res,
	}
	fmt.Printf("!!! constellationSp: %v\n", constellationSp)
	return constellationSp, nil
}

func (s *ConstellationServiceProvider) GetModuleConfig() *swconfig.ConstellationConfig {
	return s.swCfg
}

func (s *ConstellationServiceProvider) GetResources() *swconfig.ConstellationResources {
	return s.resources
}

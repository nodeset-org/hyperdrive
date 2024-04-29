package swcommon

import (
	"fmt"
	"reflect"

	"github.com/nodeset-org/hyperdrive/daemon-utils/services"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
)

type StakewiseServiceProvider struct {
	*services.ServiceProvider
	swCfg              *swconfig.StakewiseConfig
	wallet             *Wallet
	resources          *swconfig.StakewiseResources
	depositDataManager *DepositDataManager
	nodesetClient      *NodesetClient
}

// Create a new service provider with Stakewise daemon-specific features
func NewStakewiseServiceProvider(sp *services.ServiceProvider) (*StakewiseServiceProvider, error) {
	// Create the wallet
	wallet, err := NewWallet(sp)
	if err != nil {
		return nil, fmt.Errorf("error initializing wallet: %w", err)
	}

	// Create the resources
	cfg := sp.GetHyperdriveConfig()
	res := swconfig.NewStakewiseResources(cfg.Network.Value)
	swCfg, ok := sp.GetModuleConfig().(*swconfig.StakewiseConfig)
	if !ok {
		return nil, fmt.Errorf("stakewise config is not the correct type, it's a %s", reflect.TypeOf(swCfg))
	}

	// Make the provider
	stakewiseSp := &StakewiseServiceProvider{
		ServiceProvider: sp,
		swCfg:           swCfg,
		wallet:          wallet,
		resources:       res,
	}

	// Create the deposit data manager
	ddMgr, err := NewDepositDataManager(stakewiseSp)
	if err != nil {
		return nil, fmt.Errorf("error initializing deposit data manager: %w", err)
	}
	stakewiseSp.depositDataManager = ddMgr

	// Create the nodeset client
	nc := NewNodesetClient(stakewiseSp)
	stakewiseSp.nodesetClient = nc
	return stakewiseSp, nil
}

func (s *StakewiseServiceProvider) GetModuleConfig() *swconfig.StakewiseConfig {
	return s.swCfg
}

func (s *StakewiseServiceProvider) GetWallet() *Wallet {
	return s.wallet
}

func (s *StakewiseServiceProvider) GetResources() *swconfig.StakewiseResources {
	return s.resources
}

func (s *StakewiseServiceProvider) GetDepositDataManager() *DepositDataManager {
	return s.depositDataManager
}

func (s *StakewiseServiceProvider) GetNodesetClient() *NodesetClient {
	return s.nodesetClient
}

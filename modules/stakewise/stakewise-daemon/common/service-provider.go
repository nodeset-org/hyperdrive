package swcommon

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/daemon-utils/services"
	swshared "github.com/nodeset-org/hyperdrive/modules/stakewise/shared"
)

type StakewiseServiceProvider struct {
	*services.ServiceProvider
	wallet             *Wallet
	resources          *swshared.StakewiseResources
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
	cfg := sp.GetConfig()
	res := swshared.NewStakewiseResources(cfg.Network.Value)

	// Make the provider
	stakewiseSp := &StakewiseServiceProvider{
		ServiceProvider: sp,
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

func (s *StakewiseServiceProvider) GetWallet() *Wallet {
	return s.wallet
}

func (s *StakewiseServiceProvider) GetResources() *swshared.StakewiseResources {
	return s.resources
}

func (s *StakewiseServiceProvider) GetDepositDataManager() *DepositDataManager {
	return s.depositDataManager
}

func (s *StakewiseServiceProvider) GetNodesetClient() *NodesetClient {
	return s.nodesetClient
}

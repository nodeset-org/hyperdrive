package common

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/modules/common/services"
	"github.com/nodeset-org/hyperdrive/modules/stakewise/wallet"
	swutils "github.com/nodeset-org/hyperdrive/shared/utils/modules/stakewise"
)

type StakewiseServiceProvider struct {
	*services.ServiceProvider
	wallet             *wallet.Wallet
	resources          *swutils.StakewiseResources
	depositDataManager *DepositDataManager
	nodesetClient      *NodesetClient
}

// Create a new service provider with Stakewise daemon-specific features
func NewStakewiseServiceProvider(sp *services.ServiceProvider) (*StakewiseServiceProvider, error) {
	// Create the wallet
	wallet, err := wallet.NewWallet(sp)
	if err != nil {
		return nil, fmt.Errorf("error initializing wallet: %w", err)
	}

	// Create the resources
	cfg := sp.GetConfig()
	res := swutils.NewStakewiseResources(cfg.Network.Value)

	// Make the provider
	stakewiseSp := &StakewiseServiceProvider{
		ServiceProvider: sp,
		wallet:          wallet,
		resources:       res,
	}

	// Create the deposit data manager
	ddMgr := NewDepositDataManager(stakewiseSp)
	stakewiseSp.depositDataManager = ddMgr

	// Create the nodeset client
	nc := NewNodesetClient(stakewiseSp)
	stakewiseSp.nodesetClient = nc
	return stakewiseSp, nil
}

func (s *StakewiseServiceProvider) GetWallet() *wallet.Wallet {
	return s.wallet
}

func (s *StakewiseServiceProvider) GetResources() *swutils.StakewiseResources {
	return s.resources
}

func (s *StakewiseServiceProvider) GetDepositDataManager() *DepositDataManager {
	return s.depositDataManager
}

func (s *StakewiseServiceProvider) GetNodesetClient() *NodesetClient {
	return s.nodesetClient
}

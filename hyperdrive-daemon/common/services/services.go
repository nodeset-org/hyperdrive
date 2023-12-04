package services

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
)

// A container for all of the various services used by the Smartnode
type ServiceProvider struct {
	cfg *config.RocketPoolConfig
	// nodeWallet *wallet.LocalWallet
	// ecManager  *ExecutionClientManager
	// bcManager  *BeaconClientManager
	// rocketPool *rocketpool.RocketPool
	// rplFaucet          *contracts.RplFaucet
	// snapshotDelegation *contracts.SnapshotDelegation
	// docker             *client.Client
}

// ===============
// === Getters ===
// ===============

func (p *ServiceProvider) GetConfig() *config.RocketPoolConfig {
	return p.cfg
}

// func (p *ServiceProvider) GetWallet() *wallet.LocalWallet {
// 	return p.nodeWallet
// }

// func (p *ServiceProvider) GetEthClient() *ExecutionClientManager {
// 	return p.ecManager
// }

// func (p *ServiceProvider) GetRocketPool() *rocketpool.RocketPool {
// 	return p.rocketPool
// }

// func (p *ServiceProvider) GetRplFaucet() *contracts.RplFaucet {
// 	return p.rplFaucet
// }

// func (p *ServiceProvider) GetSnapshotDelegation() *contracts.SnapshotDelegation {
// 	return p.snapshotDelegation
// }

// func (p *ServiceProvider) GetBeaconClient() *BeaconClientManager {
// 	return p.bcManager
// }

// func (p *ServiceProvider) GetDocker() *client.Client {
// 	return p.docker
// }

package services

import "github.com/rocket-pool/rocketpool-go/rocketpool"

// A container for all of the various services used by the Smartnode
type ServiceProvider struct {
	// cfg *config.RocketPoolConfig
	// nodeWallet *wallet.LocalWallet
	// ecManager  *ExecutionClientManager
	bcManager  *BeaconClientManager
	rocketPool *rocketpool.RocketPool
	// rplFaucet          *contracts.RplFaucet
	// snapshotDelegation *contracts.SnapshotDelegation
	// docker             *client.Client
}

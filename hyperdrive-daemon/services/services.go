package services

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"

	"github.com/rocket-pool/smartnode/shared/services/config"
)

// A container for all of the various services used by the Smartnode
type ServiceProvider struct {
	cfg        *config.RocketPoolConfig
	rocketPool *rocketpool.RocketPool
	ecManager  *ExecutionClientManager
	bcManager  *BeaconClientManager
}

// Creates a new ServiceProvider instance
func NewServiceProvider(settingsPath string) (*ServiceProvider, error) {
	// Config
	settingsFile := os.ExpandEnv(settingsPath)
	cfg, err := loadRpConfigFromFile(settingsFile)
	if err != nil {
		return nil, fmt.Errorf("error loading Smartnode config: %w", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("Smartnode config settings file [%s] not found", settingsFile)
	}

	// EC Manager
	ecManager, err := NewExecutionClientManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating executon client manager: %w", err)
	}

	// Rocket Pool
	rp, err := rocketpool.NewRocketPool(
		ecManager,
		common.HexToAddress(cfg.Smartnode.GetStorageAddress()),
		common.HexToAddress(cfg.Smartnode.GetMulticallAddress()),
		common.HexToAddress(cfg.Smartnode.GetBalanceBatcherAddress()),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Rocket Pool binding: %w", err)
	}

	// Beacon manager
	bcManager, err := NewBeaconClientManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Beacon client manager: %w", err)
	}

	// Create the provider
	provider := &ServiceProvider{
		cfg:        cfg,
		ecManager:  ecManager,
		bcManager:  bcManager,
		rocketPool: rp,
	}
	return provider, nil
}

// ===============
// === Getters ===
// ===============

func (p *ServiceProvider) GetSmartnodeConfig() *config.RocketPoolConfig {
	return p.cfg
}

func (p *ServiceProvider) GetRocketPool() *rocketpool.RocketPool {
	return p.rocketPool
}

func (p *ServiceProvider) GetEthClient() *ExecutionClientManager {
	return p.ecManager
}

func (p *ServiceProvider) GetBeaconClient() *BeaconClientManager {
	return p.bcManager
}

// =============
// === Utils ===
// =============

// Loads a Smartnode config without updating it if it exists
func loadRpConfigFromFile(path string) (*config.RocketPoolConfig, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, nil
	}

	cfg, err := config.LoadFromFile(path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

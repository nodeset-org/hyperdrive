package services

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// A container for all of the various services used by the Smartnode
type ServiceProvider struct {
	cfg        *config.HyperdriveConfig
	rocketPool *rocketpool.RocketPool
	ecManager  *ExecutionClientManager
	bcManager  *BeaconClientManager
}

// Creates a new ServiceProvider instance
func NewServiceProvider(cfgPath string) (*ServiceProvider, error) {
	// Config
	cfg, err := loadConfigFromFile(os.ExpandEnv(cfgPath))
	if err != nil {
		return nil, fmt.Errorf("error loading hyperdrive config: %w", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("hyperdrive config settings file [%s] not found", cfgPath)
	}

	// Return an "empty" config if the Smartnode config doesn't exist
	smartnodeCfg := cfg.SmartnodeConfig
	if smartnodeCfg == nil {
		var err error
		switch cfg.SmartnodeStatus {
		case config.SmartnodeStatus_EmptyDir:
			err = fmt.Errorf("smartnode config directory has not been set yet")
		case config.SmartnodeStatus_InvalidConfig:
			err = fmt.Errorf("invalid smartnode config file: %s", cfg.SmartnodeConfigLoadErrorMessage)
		case config.SmartnodeStatus_InvalidDir:
			err = fmt.Errorf("invalid smartnode path: %s", cfg.SmartnodeConfigLoadErrorMessage)
		case config.SmartnodeStatus_MissingCfg:
			err = fmt.Errorf("the smartnode config file does not exist in the provided path")
		case config.SmartnodeStatus_Unknown:
			err = fmt.Errorf("unknown error")
		}
		return nil, fmt.Errorf("smartnode could not be loaded: %w", err)
	}

	// EC Manager
	ecManager, err := NewExecutionClientManager(cfg.SmartnodeConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating executon client manager: %w", err)
	}

	// Rocket Pool
	rp, err := rocketpool.NewRocketPool(
		ecManager,
		common.HexToAddress(smartnodeCfg.Smartnode.GetStorageAddress()),
		common.HexToAddress(smartnodeCfg.Smartnode.GetMulticallAddress()),
		common.HexToAddress(smartnodeCfg.Smartnode.GetBalanceBatcherAddress()),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Rocket Pool binding: %w", err)
	}

	// Beacon manager
	bcManager, err := NewBeaconClientManager(smartnodeCfg)
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

func (p *ServiceProvider) GetConfig() *config.HyperdriveConfig {
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

// Loads a Hyperdrive config without updating it if it exists
func loadConfigFromFile(path string) (*config.HyperdriveConfig, error) {
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

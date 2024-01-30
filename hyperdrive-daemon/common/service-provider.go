package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/nodeset-org/eth-utils/eth"
	"github.com/nodeset-org/hyperdrive/daemon-utils/services"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/wallet"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

const (
	apiLogColor color.Attribute = color.FgHiCyan
)

// A container for all of the various services used by Hyperdrive
type ServiceProvider struct {
	// Services
	cfg        *config.HyperdriveConfig
	nodeWallet *wallet.Wallet
	ecManager  *services.ExecutionClientManager
	bcManager  *services.BeaconClientManager
	docker     *client.Client
	txMgr      *eth.TransactionManager
	queryMgr   *eth.QueryManager
	resources  *utils.Resources

	// TODO: find a better place for this than the common service provider
	apiLogger *log.ColorLogger

	// Path info
	userDir string
}

// Creates a new ServiceProvider instance
func NewServiceProvider(userDir string) (*ServiceProvider, error) {
	// Config
	cfgPath := filepath.Join(userDir, config.ConfigFilename)
	cfg, err := loadConfigFromFile(os.ExpandEnv(cfgPath))
	if err != nil {
		return nil, fmt.Errorf("error loading hyperdrive config: %w", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("hyperdrive config settings file [%s] not found", cfgPath)
	}

	// Logger
	apiLogger := log.NewColorLogger(apiLogColor)

	// Resources
	resources := utils.NewResources(cfg.Network.Value)

	// Wallet
	userDataPath, err := homedir.Expand(cfg.UserDataPath.Value)
	if err != nil {
		return nil, fmt.Errorf("error expanding user data path [%s]: %w", cfg.UserDataPath.Value, err)
	}
	nodeAddressPath := filepath.Join(userDataPath, config.UserAddressFilename)
	walletDataPath := filepath.Join(userDataPath, config.UserWalletDataFilename)
	passwordPath := filepath.Join(userDataPath, config.UserPasswordFilename)
	nodeWallet, err := wallet.NewWallet(walletDataPath, nodeAddressPath, passwordPath, resources.ChainID)
	if err != nil {
		return nil, fmt.Errorf("error creating node wallet: %w", err)
	}

	// EC Manager
	ecManager, err := services.NewExecutionClientManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating executon client manager: %w", err)
	}

	// Beacon manager
	bcManager, err := services.NewBeaconClientManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Beacon client manager: %w", err)
	}

	// Docker client
	dockerClient, err := client.NewClientWithOpts(client.WithVersion(config.DockerApiVersion))
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %w", err)
	}

	// TX Manager
	txMgr, err := eth.NewTransactionManager(ecManager, eth.DefaultSafeGasBuffer, eth.DefaultSafeGasMultiplier)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction manager: %w", err)
	}

	// Query Manager - set the default concurrent run limit to half the CPUs so the EC doesn't get overwhelmed
	concurrentCallLimit := runtime.NumCPU()
	if concurrentCallLimit < 1 {
		concurrentCallLimit = 1
	}
	queryMgr := eth.NewQueryManager(ecManager, resources.MulticallAddress, concurrentCallLimit)

	// Create the provider
	provider := &ServiceProvider{
		userDir:    userDir,
		cfg:        cfg,
		nodeWallet: nodeWallet,
		ecManager:  ecManager,
		bcManager:  bcManager,
		docker:     dockerClient,
		resources:  resources,
		txMgr:      txMgr,
		queryMgr:   queryMgr,
		apiLogger:  &apiLogger,
	}
	return provider, nil
}

// ===============
// === Getters ===
// ===============

func (p *ServiceProvider) GetUserDir() string {
	return p.userDir
}

func (p *ServiceProvider) GetConfig() *config.HyperdriveConfig {
	return p.cfg
}

func (p *ServiceProvider) GetWallet() *wallet.Wallet {
	return p.nodeWallet
}

func (p *ServiceProvider) GetEthClient() *services.ExecutionClientManager {
	return p.ecManager
}

func (p *ServiceProvider) GetBeaconClient() *services.BeaconClientManager {
	return p.bcManager
}

func (p *ServiceProvider) GetDocker() *client.Client {
	return p.docker
}

func (p *ServiceProvider) GetResources() *utils.Resources {
	return p.resources
}

func (p *ServiceProvider) GetTransactionManager() *eth.TransactionManager {
	return p.txMgr
}

func (p *ServiceProvider) GetQueryManager() *eth.QueryManager {
	return p.queryMgr
}

func (p *ServiceProvider) GetApiLogger() *log.ColorLogger {
	return p.apiLogger
}

func (p *ServiceProvider) IsDebugMode() bool {
	return p.cfg.DebugMode.Value
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

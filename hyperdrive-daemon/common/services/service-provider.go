package services

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/common/wallet"
	nmkeystore "github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/common/wallet/keystore/nimbus"
	tkkeystore "github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/common/wallet/keystore/teku"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/config"
)

// A container for all of the various services used by the Smartnode
type ServiceProvider struct {
	cfg        *config.HyperdriveConfig
	nodeWallet *wallet.LocalWallet
	ecManager  *ExecutionClientManager
	bcManager  *BeaconClientManager
	docker     *client.Client
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

	// Wallet
	chainID := cfg.GetChainID()
	nodeAddressPath := "" // os.ExpandEnv(cfg.Smartnode.GetNodeAddressPath())
	keystorePath := ""    //os.ExpandEnv(cfg.Smartnode.GetWalletPath())
	passwordPath := ""    //os.ExpandEnv(cfg.Smartnode.GetPasswordPath())
	nodeWallet, err := wallet.NewLocalWallet(nodeAddressPath, keystorePath, passwordPath, chainID, true)
	if err != nil {
		return nil, fmt.Errorf("error creating node wallet: %w", err)
	}

	// Keystores
	validatorKeychainPath := "" //os.ExpandEnv(cfg.Smartnode.GetValidatorKeychainPath())
	nimbusKeystore := nmkeystore.NewKeystore(validatorKeychainPath)
	tekuKeystore := tkkeystore.NewKeystore(validatorKeychainPath)
	nodeWallet.AddValidatorKeystore("nimbus", nimbusKeystore)
	nodeWallet.AddValidatorKeystore("teku", tekuKeystore)

	// EC Manager
	ecManager, err := NewExecutionClientManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating executon client manager: %w", err)
	}

	// Beacon manager
	bcManager, err := NewBeaconClientManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Beacon client manager: %w", err)
	}

	// Docker client
	dockerClient, err := client.NewClientWithOpts(client.WithVersion(config.DockerApiVersion))
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %w", err)
	}

	// Create the provider
	provider := &ServiceProvider{
		cfg:        cfg,
		nodeWallet: nodeWallet,
		ecManager:  ecManager,
		bcManager:  bcManager,
		docker:     dockerClient,
	}
	return provider, nil
}

// ===============
// === Getters ===
// ===============

func (p *ServiceProvider) GetConfig() *config.HyperdriveConfig {
	return p.cfg
}

func (p *ServiceProvider) GetWallet() *wallet.LocalWallet {
	return p.nodeWallet
}

func (p *ServiceProvider) GetEthClient() *ExecutionClientManager {
	return p.ecManager
}

func (p *ServiceProvider) GetBeaconClient() *BeaconClientManager {
	return p.bcManager
}

func (p *ServiceProvider) GetDocker() *client.Client {
	return p.docker
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

	cfg, err := config.LoadFromFile(path, true)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

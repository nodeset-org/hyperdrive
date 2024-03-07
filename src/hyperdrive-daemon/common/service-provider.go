package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive/shared/config"
	nmc_config "github.com/rocket-pool/node-manager-core/config"
	nmc_services "github.com/rocket-pool/node-manager-core/node/services"
)

// A container for all of the various services used by Hyperdrive
type ServiceProvider struct {
	*nmc_services.ServiceProvider

	// Services
	cfg       *config.HyperdriveConfig
	resources *nmc_config.NetworkResources

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
	resources := cfg.GetNetworkResources()

	// Core provider
	sp, err := nmc_services.NewServiceProvider(cfg, config.ClientTimeout, cfg.DebugMode.Value)
	if err != nil {
		return nil, fmt.Errorf("error creating core service provider: %w", err)
	}

	// Create the provider
	provider := &ServiceProvider{
		ServiceProvider: sp,
		userDir:         userDir,
		cfg:             cfg,
		resources:       resources,
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

func (p *ServiceProvider) GetResources() *nmc_config.NetworkResources {
	return p.resources
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

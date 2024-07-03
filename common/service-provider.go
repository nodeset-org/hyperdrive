package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/rocket-pool/node-manager-core/node/services"
)

// A container for all of the various services used by Hyperdrive
type ServiceProvider struct {
	*services.ServiceProvider

	// Services
	cfg *hdconfig.HyperdriveConfig
	res *hdconfig.HyperdriveResources
	ns  *NodeSetServiceManager

	// Path info
	userDir string
}

// Creates a new ServiceProvider instance by loading the Hyperdrive config in the provided directory
func NewServiceProvider(userDir string) (*ServiceProvider, error) {
	// Config
	cfgPath := filepath.Join(userDir, hdconfig.ConfigFilename)
	cfg, err := loadConfigFromFile(os.ExpandEnv(cfgPath))
	if err != nil {
		return nil, fmt.Errorf("error loading hyperdrive config: %w", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("hyperdrive config settings file [%s] not found", cfgPath)
	}

	return NewServiceProviderFromConfig(cfg)
}

// Creates a new ServiceProvider instance directly from a Hyperdrive config instead of loading it from the filesystem
func NewServiceProviderFromConfig(cfg *hdconfig.HyperdriveConfig) (*ServiceProvider, error) {
	// Core provider
	sp, err := services.NewServiceProvider(cfg, hdconfig.ClientTimeout)
	if err != nil {
		return nil, fmt.Errorf("error creating core service provider: %w", err)
	}

	// Create the provider
	provider := &ServiceProvider{
		ServiceProvider: sp,
		userDir:         cfg.GetUserDirectory(),
		cfg:             cfg,
		res:             cfg.GetResources(),
	}
	ns := NewNodeSetServiceManager(provider)
	provider.ns = ns
	return provider, nil
}

// Creates a new ServiceProvider instance from custom services and artifacts
func NewServiceProviderFromCustomServices(cfg *hdconfig.HyperdriveConfig, resources *hdconfig.HyperdriveResources, ecManager *services.ExecutionClientManager, bnManager *services.BeaconClientManager, docker client.APIClient) (*ServiceProvider, error) {
	// Core provider
	sp, err := services.NewServiceProviderWithCustomServices(cfg, resources.NetworkResources, ecManager, bnManager, docker)
	if err != nil {
		return nil, fmt.Errorf("error creating core service provider: %w", err)
	}

	// Create the provider
	provider := &ServiceProvider{
		ServiceProvider: sp,
		userDir:         cfg.GetUserDirectory(),
		cfg:             cfg,
		res:             resources,
	}
	ns := NewNodeSetServiceManager(provider)
	provider.ns = ns
	return provider, nil
}

// ===============
// === Getters ===
// ===============

func (p *ServiceProvider) GetUserDir() string {
	return p.userDir
}

func (p *ServiceProvider) GetConfig() *hdconfig.HyperdriveConfig {
	return p.cfg
}

func (p *ServiceProvider) GetResources() *hdconfig.HyperdriveResources {
	return p.res
}

func (p *ServiceProvider) GetNodeSetServiceManager() *NodeSetServiceManager {
	return p.ns
}

// =============
// === Utils ===
// =============

// Loads a Hyperdrive config without updating it if it exists
func loadConfigFromFile(path string) (*hdconfig.HyperdriveConfig, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, nil
	}

	cfg, err := hdconfig.LoadFromFile(path)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

package common

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/rocket-pool/node-manager-core/node/services"
)

// Provides Hyperdrive's configuration
type IHyperdriveConfigProvider interface {
	// Gets Hyperdrive's configuration
	GetConfig() *hdconfig.HyperdriveConfig

	// Gets Hyperdrive's list of resources
	GetResources() *hdconfig.HyperdriveResources
}

// Provides a manager for nodeset.io communications
type INodeSetManagerProvider interface {
	// Gets the NodeSetServiceManager
	GetNodeSetServiceManager() *NodeSetServiceManager
}

// Provides methods for requiring or waiting for various conditions to be met
type IRequirementsProvider interface {
	// Require Hyperdrive has a node address set
	RequireNodeAddress() error

	// Require Hyperdrive has a wallet that's loaded and ready for transactions
	RequireWalletReady() error

	// Require that the Ethereum client is synced
	RequireEthClientSynced(ctx context.Context) error

	// Require that the Beacon chain client is synced
	RequireBeaconClientSynced(ctx context.Context) error

	// Require the node has been registered with a nodeset.io account
	RequireRegisteredWithNodeSet(ctx context.Context) error

	// Wait for the Ethereum client to be synced
	WaitEthClientSynced(ctx context.Context, verbose bool) error

	// Wait for the Beacon chain client to be synced
	WaitBeaconClientSynced(ctx context.Context, verbose bool) error

	// Wait for the node to have a wallet loaded and ready for transactions
	WaitForWallet(ctx context.Context) error

	// Wait for the node to be registered with a nodeset.io account
	WaitForNodeSetRegistration(ctx context.Context) bool
}

// Provides access to all of Hyperdrive's services
type IHyperdriveServiceProvider interface {
	IHyperdriveConfigProvider
	INodeSetManagerProvider
	IRequirementsProvider
	services.IServiceProvider
}

// A container for all of the various services used by Hyperdrive
type serviceProvider struct {
	services.IServiceProvider

	// Services
	cfg *hdconfig.HyperdriveConfig
	res *hdconfig.HyperdriveResources
	ns  *NodeSetServiceManager

	// Path info
	userDir string
}

// Creates a new IHyperdriveServiceProvider instance by loading the Hyperdrive config in the provided directory
func NewHyperdriveServiceProvider(userDir string) (IHyperdriveServiceProvider, error) {
	// Config
	cfgPath := filepath.Join(userDir, hdconfig.ConfigFilename)
	cfg, err := loadConfigFromFile(os.ExpandEnv(cfgPath))
	if err != nil {
		return nil, fmt.Errorf("error loading hyperdrive config: %w", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("hyperdrive config settings file [%s] not found", cfgPath)
	}

	// Make the resources
	resources := hdconfig.NewHyperdriveResources(cfg.Network.Value)

	return NewHyperdriveServiceProviderFromConfig(cfg, resources)
}

// Creates a new IHyperdriveServiceProvider instance directly from a Hyperdrive config and resources list instead of loading them from the filesystem
func NewHyperdriveServiceProviderFromConfig(cfg *hdconfig.HyperdriveConfig, resources *hdconfig.HyperdriveResources) (IHyperdriveServiceProvider, error) {
	// Core provider
	sp, err := services.NewServiceProvider(cfg, resources.NetworkResources, hdconfig.ClientTimeout)
	if err != nil {
		return nil, fmt.Errorf("error creating core service provider: %w", err)
	}

	// Create the provider
	provider := &serviceProvider{
		IServiceProvider: sp,
		userDir:          cfg.GetUserDirectory(),
		cfg:              cfg,
		res:              resources,
	}
	ns := NewNodeSetServiceManager(provider)
	provider.ns = ns
	return provider, nil
}

// Creates a new IHyperdriveServiceProvider instance from custom services and artifacts
func NewHyperdriveServiceProviderFromCustomServices(cfg *hdconfig.HyperdriveConfig, resources *hdconfig.HyperdriveResources, ecManager *services.ExecutionClientManager, bnManager *services.BeaconClientManager, docker client.APIClient) (IHyperdriveServiceProvider, error) {
	// Core provider
	sp, err := services.NewServiceProviderWithCustomServices(cfg, resources.NetworkResources, ecManager, bnManager, docker)
	if err != nil {
		return nil, fmt.Errorf("error creating core service provider: %w", err)
	}

	// Create the provider
	provider := &serviceProvider{
		IServiceProvider: sp,
		userDir:          cfg.GetUserDirectory(),
		cfg:              cfg,
		res:              resources,
	}
	ns := NewNodeSetServiceManager(provider)
	provider.ns = ns
	return provider, nil
}

// ===============
// === Getters ===
// ===============

func (p *serviceProvider) GetConfig() *hdconfig.HyperdriveConfig {
	return p.cfg
}

func (p *serviceProvider) GetResources() *hdconfig.HyperdriveResources {
	return p.res
}

func (p *serviceProvider) GetNodeSetServiceManager() *NodeSetServiceManager {
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

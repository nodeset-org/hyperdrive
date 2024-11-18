package common

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/node/services"
)

// ==================
// === Interfaces ===
// ==================

// Provides Hyperdrive's configuration
type IHyperdriveConfigProvider interface {
	// Gets Hyperdrive's configuration
	GetConfig() *hdconfig.HyperdriveConfig

	// Gets Hyperdrive's list of resources
	GetResources() *hdconfig.MergedResources
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

// ========================
// === Service Provider ===
// ========================

// A container for all of the various services used by Hyperdrive
type serviceProvider struct {
	services.IServiceProvider

	// Services
	cfg *hdconfig.HyperdriveConfig
	res *hdconfig.MergedResources
	ns  *NodeSetServiceManager

	// Path info
	userDir string
}

// Creates a new IHyperdriveServiceProvider instance by loading the Hyperdrive config in the provided directory
func NewHyperdriveServiceProvider(userDir string, resourcesDir string) (IHyperdriveServiceProvider, error) {
	// Load the network settings
	settingsList, err := hdconfig.LoadSettingsFiles(resourcesDir)
	if err != nil {
		return nil, fmt.Errorf("error loading network settings: %w", err)
	}

	// Create the config
	cfgPath := filepath.Join(userDir, hdconfig.ConfigFilename)
	cfg, err := loadConfigFromFile(os.ExpandEnv(cfgPath), settingsList)
	if err != nil {
		return nil, fmt.Errorf("error loading hyperdrive config: %w", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("hyperdrive config settings file [%s] not found", cfgPath)
	}

	// Get the resources from the selected network
	var selectedResources *hdconfig.MergedResources
	for _, network := range settingsList {
		if network.Key == cfg.Network.Value {
			selectedResources = &hdconfig.MergedResources{
				NetworkResources:    network.NetworkResources,
				HyperdriveResources: network.HyperdriveResources,
			}
			break
		}
	}
	if selectedResources == nil {
		return nil, fmt.Errorf("no resources found for selected network [%s]", cfg.Network.Value)
	}

	return NewHyperdriveServiceProviderFromConfig(cfg, selectedResources)
}

// Creates a new IHyperdriveServiceProvider instance directly from a Hyperdrive config and resources list instead of loading them from the filesystem
func NewHyperdriveServiceProviderFromConfig(cfg *hdconfig.HyperdriveConfig, resources *hdconfig.MergedResources) (IHyperdriveServiceProvider, error) {
	// Core provider
	sp, err := services.NewServiceProvider(cfg, resources.NetworkResources, time.Duration(cfg.ClientTimeout.Value)*time.Second)
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
func NewHyperdriveServiceProviderFromCustomServices(cfg *hdconfig.HyperdriveConfig, resources *hdconfig.MergedResources, ecManager *services.ExecutionClientManager, bnManager *services.BeaconClientManager, docker client.APIClient) (IHyperdriveServiceProvider, error) {
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

func (p *serviceProvider) GetResources() *hdconfig.MergedResources {
	return p.res
}

func (p *serviceProvider) GetNodeSetServiceManager() *NodeSetServiceManager {
	return p.ns
}

// =============
// === Utils ===
// =============

// Loads a Hyperdrive config without updating it if it exists
func loadConfigFromFile(configPath string, networks []*hdconfig.HyperdriveSettings) (*hdconfig.HyperdriveConfig, error) {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	cfg, err := hdconfig.LoadFromFile(configPath, networks)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

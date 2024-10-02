package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nodeset-org/hyperdrive-daemon/client"
	"github.com/nodeset-org/hyperdrive-daemon/shared/auth"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	bclient "github.com/rocket-pool/node-manager-core/beacon/client"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/services"
	"github.com/rocket-pool/node-manager-core/wallet"
)

// ==================
// === Interfaces ===
// ==================

// Provides the configurations for Hyperdrive and the module
type IModuleConfigProvider interface {
	// Gets Hyperdrive's configuration
	GetHyperdriveConfig() *hdconfig.HyperdriveConfig

	// Gets Hyperdrive's list of resources
	GetHyperdriveResources() *hdconfig.MergedResources

	// Gets the module's configuration
	GetModuleConfig() hdconfig.IModuleConfig

	// Gets the path to the module's data directory
	GetModuleDir() string
}

// Provides a Hyperdrive API client
type IHyperdriveClientProvider interface {
	// Gets the Hyperdrive client
	GetHyperdriveClient() *client.ApiClient
}

// Provides a signer that can sign messages from the node's wallet
type ISignerProvider interface {
	// Gets the module's signer
	GetSigner() *ModuleSigner
}

// Provides access to the module's loggers
type ILoggerProvider interface {
	services.ILoggerProvider

	// Gets the logger for the Hyperdrive client
	GetClientLogger() *log.Logger
}

// Provides methods for requiring or waiting for various conditions to be met
type IRequirementsProvider interface {
	// Require Hyperdrive has a node address set
	RequireNodeAddress(status wallet.WalletStatus) error

	// Require Hyperdrive has a wallet that's loaded and ready for transactions
	RequireWalletReady(status wallet.WalletStatus) error

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

	// Wait for Hyperdrive to have a node address assigned
	WaitForNodeAddress(ctx context.Context) (*wallet.WalletStatus, error)

	// Wait for the node to have a wallet loaded and ready for transactions
	WaitForWallet(ctx context.Context) (*wallet.WalletStatus, error)

	// Wait for the node to be registered with a nodeset.io account
	WaitForNodeSetRegistration(ctx context.Context) bool
}

// Provides access to all of the standard services the module can use
type IModuleServiceProvider interface {
	IModuleConfigProvider
	IHyperdriveClientProvider
	ISignerProvider
	ILoggerProvider
	IRequirementsProvider

	// Standard NMC interfaces
	services.IEthClientProvider
	services.IBeaconClientProvider
	services.IContextProvider
	io.Closer
}

// ========================
// === Service Provider ===
// ========================

// A container for all of the various services used by Hyperdrive
type moduleServiceProvider struct {
	// Services
	hdCfg        *hdconfig.HyperdriveConfig
	moduleConfig hdconfig.IModuleConfig
	hdClient     *client.ApiClient
	ecManager    *services.ExecutionClientManager
	bcManager    *services.BeaconClientManager
	resources    *hdconfig.MergedResources
	signer       *ModuleSigner
	txMgr        *eth.TransactionManager
	queryMgr     *eth.QueryManager
	ctx          context.Context
	cancel       context.CancelFunc

	// Logging
	clientLogger *log.Logger
	apiLogger    *log.Logger
	tasksLogger  *log.Logger

	// Path info
	moduleDir string
	userDir   string
}

// Creates a new IModuleServiceProvider instance
func NewModuleServiceProvider[ConfigType hdconfig.IModuleConfig](hyperdriveUrl *url.URL, moduleDir string, moduleName string, clientLogName string, factory func(*hdconfig.HyperdriveConfig) (ConfigType, error), authMgr *auth.AuthorizationManager) (IModuleServiceProvider, error) {
	hdCfg, resources, hdClient, err := getHdConfig(hyperdriveUrl, authMgr)
	if err != nil {
		return nil, fmt.Errorf("error getting Hyperdrive config: %w", err)
	}
	clientTimeout := time.Duration(hdCfg.ClientTimeout.Value) * time.Second

	// EC Manager
	var ecManager *services.ExecutionClientManager
	primaryEcUrl, fallbackEcUrl := hdCfg.GetExecutionClientUrls()
	primaryEc, err := ethclient.Dial(primaryEcUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to primary EC at [%s]: %w", primaryEcUrl, err)
	}
	if fallbackEcUrl != "" {
		// Get the fallback EC url, if applicable
		fallbackEc, err := ethclient.Dial(fallbackEcUrl)
		if err != nil {
			return nil, fmt.Errorf("error connecting to fallback EC at [%s]: %w", fallbackEcUrl, err)
		}
		ecManager = services.NewExecutionClientManagerWithFallback(primaryEc, fallbackEc, resources.ChainID, clientTimeout)
	} else {
		ecManager = services.NewExecutionClientManager(primaryEc, resources.ChainID, clientTimeout)
	}

	// Beacon manager
	var bcManager *services.BeaconClientManager
	primaryBnUrl, fallbackBnUrl := hdCfg.GetBeaconNodeUrls()
	primaryBc := bclient.NewStandardHttpClient(primaryBnUrl, clientTimeout)
	if fallbackBnUrl != "" {
		fallbackBc := bclient.NewStandardHttpClient(fallbackBnUrl, clientTimeout)
		bcManager = services.NewBeaconClientManagerWithFallback(primaryBc, fallbackBc, resources.ChainID, clientTimeout)
	} else {
		bcManager = services.NewBeaconClientManager(primaryBc, resources.ChainID, clientTimeout)
	}

	// Get the module config
	moduleCfg, err := factory(hdCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating module config: %w", err)
	}
	modCfgEnrty, exists := hdCfg.Modules[moduleName]
	if !exists {
		return nil, fmt.Errorf("config section for module [%s] not found", moduleName)
	}
	modCfgMap, ok := modCfgEnrty.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("config section for module [%s] is not a map, it's a %s", moduleName, reflect.TypeOf(modCfgMap))
	}
	err = config.Deserialize(moduleCfg, modCfgMap, hdCfg.Network.Value)
	if err != nil {
		return nil, fmt.Errorf("error deserialzing config for module [%s]: %w", moduleName, err)
	}

	return NewModuleServiceProviderFromArtifacts(hdClient, hdCfg, moduleCfg, resources, moduleDir, moduleName, clientLogName, ecManager, bcManager)
}

// Creates a new IModuleServiceProvider instance, using the given artifacts instead of creating ones based on the config parameters
func NewModuleServiceProviderFromArtifacts(hdClient *client.ApiClient, hdCfg *hdconfig.HyperdriveConfig, moduleCfg hdconfig.IModuleConfig, resources *hdconfig.MergedResources, moduleDir string, moduleName string, clientLogName string, ecManager *services.ExecutionClientManager, bcManager *services.BeaconClientManager) (IModuleServiceProvider, error) {
	// Set up the client logger
	logPath := hdCfg.GetModuleLogFilePath(moduleName, clientLogName)
	clientLogger, err := log.NewLogger(logPath, hdCfg.GetLoggerOptions())
	if err != nil {
		return nil, fmt.Errorf("error creating HD Client logger: %w", err)
	}
	hdClient.SetLogger(clientLogger.Logger)

	// Make the API logger
	apiLogPath := hdCfg.GetModuleLogFilePath(moduleName, moduleCfg.GetApiLogFileName())
	apiLogger, err := log.NewLogger(apiLogPath, hdCfg.GetLoggerOptions())
	if err != nil {
		return nil, fmt.Errorf("error creating API logger: %w", err)
	}

	// Make the tasks logger
	tasksLogPath := hdCfg.GetModuleLogFilePath(moduleName, moduleCfg.GetTasksLogFileName())
	tasksLogger, err := log.NewLogger(tasksLogPath, hdCfg.GetLoggerOptions())
	if err != nil {
		return nil, fmt.Errorf("error creating tasks logger: %w", err)
	}

	// Signer
	signer := NewModuleSigner(hdClient)

	// TX Manager
	txMgr, err := eth.NewTransactionManager(ecManager, eth.DefaultSafeGasBuffer, eth.DefaultSafeGasMultiplier)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction manager: %w", err)
	}

	// Query Manager - set the default concurrent run limit to half the CPUs so the EC doesn't get overwhelmed
	concurrentCallLimit := runtime.NumCPU() / 2
	if concurrentCallLimit < 1 {
		concurrentCallLimit = 1
	}
	queryMgr := eth.NewQueryManager(ecManager, resources.MulticallAddress, concurrentCallLimit)

	// Context for handling task cancellation during shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Log startup
	clientLogger.Info("Starting Hyperdrive Client logger.")
	apiLogger.Info("Starting API logger.")
	tasksLogger.Info("Starting Tasks logger.")

	// Create the provider
	provider := &moduleServiceProvider{
		moduleDir:    moduleDir,
		userDir:      hdCfg.GetUserDirectory(),
		hdCfg:        hdCfg,
		moduleConfig: moduleCfg,
		ecManager:    ecManager,
		bcManager:    bcManager,
		hdClient:     hdClient,
		resources:    resources,
		signer:       signer,
		txMgr:        txMgr,
		queryMgr:     queryMgr,
		clientLogger: clientLogger,
		apiLogger:    apiLogger,
		tasksLogger:  tasksLogger,
		ctx:          ctx,
		cancel:       cancel,
	}
	return provider, nil
}

// Closes the service provider and its underlying services
func (p *moduleServiceProvider) Close() error {
	p.clientLogger.Close()
	p.apiLogger.Close()
	p.tasksLogger.Close()
	return nil
}

// ===============
// === Getters ===
// ===============

func (p *moduleServiceProvider) GetModuleDir() string {
	return p.moduleDir
}

func (p *moduleServiceProvider) GetHyperdriveConfig() *hdconfig.HyperdriveConfig {
	return p.hdCfg
}

func (p *moduleServiceProvider) GetHyperdriveResources() *hdconfig.MergedResources {
	return p.resources
}

func (p *moduleServiceProvider) GetModuleConfig() hdconfig.IModuleConfig {
	return p.moduleConfig
}

func (p *moduleServiceProvider) GetEthClient() *services.ExecutionClientManager {
	return p.ecManager
}

func (p *moduleServiceProvider) GetBeaconClient() *services.BeaconClientManager {
	return p.bcManager
}

func (p *moduleServiceProvider) GetHyperdriveClient() *client.ApiClient {
	return p.hdClient
}

func (p *moduleServiceProvider) GetSigner() *ModuleSigner {
	return p.signer
}

func (p *moduleServiceProvider) GetTransactionManager() *eth.TransactionManager {
	return p.txMgr
}

func (p *moduleServiceProvider) GetQueryManager() *eth.QueryManager {
	return p.queryMgr
}

func (p *moduleServiceProvider) GetClientLogger() *log.Logger {
	return p.clientLogger
}

func (p *moduleServiceProvider) GetApiLogger() *log.Logger {
	return p.apiLogger
}

func (p *moduleServiceProvider) GetTasksLogger() *log.Logger {
	return p.tasksLogger
}

func (p *moduleServiceProvider) GetBaseContext() context.Context {
	return p.ctx
}

func (p *moduleServiceProvider) CancelContextOnShutdown() {
	p.cancel()
}

// ==========================
// === Internal Functions ===
// ==========================

func getHdConfig(hyperdriveUrl *url.URL, authMgr *auth.AuthorizationManager) (*hdconfig.HyperdriveConfig, *hdconfig.MergedResources, *client.ApiClient, error) {
	// Add the API client route if missing
	hyperdriveUrl.Path = strings.TrimSuffix(hyperdriveUrl.Path, "/")
	if hyperdriveUrl.Path == "" {
		hyperdriveUrl.Path = fmt.Sprintf("%s/%s", hyperdriveUrl.Path, hdconfig.HyperdriveApiClientRoute)
	}

	// Create a client for the Hyperdrive daemon
	defaultLogger := slog.Default()
	hdClient := client.NewApiClient(hyperdriveUrl, defaultLogger, nil, authMgr)

	// Get the Hyperdrive settings for the selected network
	settingsResponse, err := hdClient.Service.GetNetworkSettings()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting resources from Hyperdrive server: %w", err)
	}
	settings := settingsResponse.Data.Settings

	// Create the Hyperdrive network
	hdCfg, err := hdconfig.NewHyperdriveConfig("", []*hdconfig.HyperdriveSettings{settings})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating Hyperdrive config: %w", err)
	}
	cfgResponse, err := hdClient.Service.GetConfig()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting config from Hyperdrive server: %w", err)
	}
	err = hdCfg.Deserialize(cfgResponse.Data.Config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error deserializing Hyperdrive config: %w", err)
	}

	res := &hdconfig.MergedResources{
		NetworkResources:    settings.NetworkResources,
		HyperdriveResources: settings.HyperdriveResources,
	}
	return hdCfg, res, hdClient, nil
}

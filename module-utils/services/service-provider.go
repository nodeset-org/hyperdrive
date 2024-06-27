package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nodeset-org/hyperdrive-daemon/client"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	bclient "github.com/rocket-pool/node-manager-core/beacon/client"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/services"
)

// A container for all of the various services used by Hyperdrive
type ServiceProvider struct {
	// Services
	hdCfg        *hdconfig.HyperdriveConfig
	moduleConfig hdconfig.IModuleConfig
	hdClient     *client.ApiClient
	ecManager    *services.ExecutionClientManager
	bcManager    *services.BeaconClientManager
	resources    *hdconfig.HyperdriveResources
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

// Creates a new ServiceProvider instance
func NewServiceProvider[ConfigType hdconfig.IModuleConfig](hyperdriveUrl *url.URL, moduleDir string, moduleName string, clientLogName string, factory func(*hdconfig.HyperdriveConfig) ConfigType, clientTimeout time.Duration) (*ServiceProvider, error) {
	hdCfg, hdClient, err := getHdConfig(hyperdriveUrl)
	if err != nil {
		return nil, fmt.Errorf("error getting Hyperdrive config: %w", err)
	}

	// Resources
	resources := hdCfg.GetResources()

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
	moduleCfg := factory(hdCfg)
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

	return NewServiceProviderFromArtifacts(hdClient, hdCfg, moduleCfg, resources, moduleDir, moduleName, clientLogName, ecManager, bcManager)
}

// Creates a new ServiceProvider instance, using the given artifacts instead of creating ones based on the config parameters
func NewServiceProviderFromArtifacts(hdClient *client.ApiClient, hdCfg *hdconfig.HyperdriveConfig, moduleCfg hdconfig.IModuleConfig, resources *hdconfig.HyperdriveResources, moduleDir string, moduleName string, clientLogName string, ecManager *services.ExecutionClientManager, bcManager *services.BeaconClientManager) (*ServiceProvider, error) {
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
	provider := &ServiceProvider{
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
func (p *ServiceProvider) Close() {
	p.clientLogger.Close()
	p.apiLogger.Close()
	p.tasksLogger.Close()
}

// ===============
// === Getters ===
// ===============

func (p *ServiceProvider) GetModuleDir() string {
	return p.moduleDir
}

func (p *ServiceProvider) GetUserDir() string {
	return p.userDir
}

func (p *ServiceProvider) GetHyperdriveConfig() *hdconfig.HyperdriveConfig {
	return p.hdCfg
}

func (p *ServiceProvider) GetModuleConfig() hdconfig.IModuleConfig {
	return p.moduleConfig
}

func (p *ServiceProvider) GetEthClient() *services.ExecutionClientManager {
	return p.ecManager
}

func (p *ServiceProvider) GetBeaconClient() *services.BeaconClientManager {
	return p.bcManager
}

func (p *ServiceProvider) GetHyperdriveClient() *client.ApiClient {
	return p.hdClient
}

func (p *ServiceProvider) GetResources() *hdconfig.HyperdriveResources {
	return p.resources
}

func (p *ServiceProvider) GetSigner() *ModuleSigner {
	return p.signer
}

func (p *ServiceProvider) GetTransactionManager() *eth.TransactionManager {
	return p.txMgr
}

func (p *ServiceProvider) GetQueryManager() *eth.QueryManager {
	return p.queryMgr
}

func (p *ServiceProvider) GetClientLogger() *log.Logger {
	return p.clientLogger
}

func (p *ServiceProvider) GetApiLogger() *log.Logger {
	return p.apiLogger
}

func (p *ServiceProvider) GetTasksLogger() *log.Logger {
	return p.tasksLogger
}

func (p *ServiceProvider) GetBaseContext() context.Context {
	return p.ctx
}

func (p *ServiceProvider) CancelContextOnShutdown() {
	p.cancel()
}

// ==========================
// === Internal Functions ===
// ==========================

func getHdConfig(hyperdriveUrl *url.URL) (*hdconfig.HyperdriveConfig, *client.ApiClient, error) {
	// Add the API client route if missing
	hyperdriveUrl.Path = strings.TrimSuffix(hyperdriveUrl.Path, "/")
	if hyperdriveUrl.Path == "" {
		hyperdriveUrl.Path = fmt.Sprintf("%s/%s", hyperdriveUrl.Path, hdconfig.HyperdriveApiClientRoute)
	}

	// Create a client for the Hyperdrive daemon
	defaultLogger := slog.Default()
	hdClient := client.NewApiClient(hyperdriveUrl, defaultLogger, nil)

	// Get the Hyperdrive config
	hdCfg := hdconfig.NewHyperdriveConfig("")
	cfgResponse, err := hdClient.Service.GetConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting config from Hyperdrive server: %w", err)
	}
	err = hdCfg.Deserialize(cfgResponse.Data.Config)
	if err != nil {
		return nil, nil, fmt.Errorf("error deserializing Hyperdrive config: %w", err)
	}

	return hdCfg, hdClient, nil
}

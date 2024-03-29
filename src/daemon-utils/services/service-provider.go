package services

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/nodeset-org/hyperdrive/client"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
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
	resources    *config.NetworkResources
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
func NewServiceProvider[ConfigType hdconfig.IModuleConfig](moduleDir string, moduleName string, clientLogName string, factory func(*hdconfig.HyperdriveConfig) ConfigType, clientTimeout time.Duration) (*ServiceProvider, error) {
	// Create a client for the Hyperdrive daemon
	defaultLogger := slog.Default()
	hyperdriveSocket := filepath.Join(moduleDir, hdconfig.HyperdriveCliSocketFilename)
	hdClient := client.NewApiClient(hdconfig.HyperdriveApiClientRoute, hyperdriveSocket, defaultLogger)

	// Get the Hyperdrive config
	hdCfg := hdconfig.NewHyperdriveConfig("")
	cfgResponse, err := hdClient.Service.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting config from Hyperdrive server: %w", err)
	}
	err = hdCfg.Deserialize(cfgResponse.Data.Config)
	if err != nil {
		return nil, fmt.Errorf("error deserializing Hyperdrive config: %w", err)
	}

	// Set up the client logger
	debugMode := hdCfg.DebugMode.Value
	moduleLogDir := filepath.Join(hdCfg.HyperdriveUserDirectory, hdconfig.LogDir, moduleName)
	logPath := filepath.Join(moduleLogDir, clientLogName)
	clientLogger, err := log.NewLogger(logPath, debugMode, false)
	hdClient.SetLogger(clientLogger.Logger)

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

	// Make the API logger
	apiLogPath := filepath.Join(moduleLogDir, moduleCfg.GetApiLogFileName())
	apiLogger, err := log.NewLogger(apiLogPath, debugMode, false)
	if err != nil {
		return nil, fmt.Errorf("error creating API logger: %w", err)
	}

	// Make the tasks logger
	tasksLogPath := filepath.Join(moduleLogDir, moduleCfg.GetTasksLogFileName())
	tasksLogger, err := log.NewLogger(tasksLogPath, debugMode, false)
	if err != nil {
		return nil, fmt.Errorf("error creating tasks logger: %w", err)
	}

	// Resources
	resources := hdCfg.GetNetworkResources()

	// Signer
	signer := NewModuleSigner(hdClient)

	// EC Manager
	primaryEcUrl, fallbackEcUrl := hdCfg.GetExecutionClientUrls()
	ecManager, err := services.NewExecutionClientManager(primaryEcUrl, fallbackEcUrl, resources.ChainID, clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("error creating executon client manager: %w", err)
	}

	// Beacon manager
	primaryBnUrl, fallbackBnUrl := hdCfg.GetBeaconNodeUrls()
	bcManager, err := services.NewBeaconClientManager(primaryBnUrl, fallbackBnUrl, clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("error creating Beacon client manager: %w", err)
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

	// Context for handling task cancellation during shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create the provider
	provider := &ServiceProvider{
		moduleDir:    moduleDir,
		userDir:      hdCfg.HyperdriveUserDirectory,
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

func (p *ServiceProvider) GetResources() *config.NetworkResources {
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

func (p *ServiceProvider) GetApiLogger() *log.Logger {
	return p.apiLogger
}

func (p *ServiceProvider) GetTasksLogger() *log.Logger {
	return p.tasksLogger
}

func (p *ServiceProvider) IsDebugMode() bool {
	return p.hdCfg.DebugMode.Value
}

func (p *ServiceProvider) GetBaseContext() context.Context {
	return p.ctx
}

func (p *ServiceProvider) CancelContextOnShutdown() {
	p.cancel()
}

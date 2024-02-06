package services

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	docker "github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/nodeset-org/eth-utils/eth"
	"github.com/nodeset-org/hyperdrive/client"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

const (
	apiLogColor color.Attribute = color.FgHiCyan
)

// A container for all of the various services used by Hyperdrive
type ServiceProvider[ConfigType config.IModuleConfig] struct {
	// Services
	hdCfg        *config.HyperdriveConfig
	moduleConfig ConfigType
	hdClient     *client.ApiClient
	ecManager    *ExecutionClientManager
	bcManager    *BeaconClientManager
	docker       *docker.Client
	txMgr        *eth.TransactionManager
	queryMgr     *eth.QueryManager
	resources    *utils.Resources

	// TODO: find a better place for this than the common service provider
	apiLogger *log.ColorLogger

	// Path info
	moduleDir string
	userDir   string
}

// Creates a new ServiceProvider instance
func NewServiceProvider[ConfigType config.IModuleConfig](moduleDir string, factory func(*config.HyperdriveConfig) ConfigType) (*ServiceProvider[ConfigType], error) {
	// Create a client for the Hyperdrive daemon
	hyperdriveSocket := filepath.Join(moduleDir, config.HyperdriveSocketFilename)
	hdClient := client.NewApiClient(config.HyperdriveDaemonRoute, hyperdriveSocket, false)

	// Get the Hyperdrive config
	hdCfg := config.NewHyperdriveConfig("")
	cfgResponse, err := hdClient.Service.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting config from Hyperdrive server: %w", err)
	}
	err = hdCfg.Deserialize(cfgResponse.Data.Config)
	if err != nil {
		return nil, fmt.Errorf("error deserializing Hyperdrive config: %w", err)
	}
	hdClient.SetDebug(hdCfg.DebugMode.Value)

	// Get the module config
	moduleCfg := factory(hdCfg)
	moduleName := moduleCfg.GetModuleName()
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

	// Logger
	apiLogger := log.NewColorLogger(apiLogColor)

	// Resources
	resources := utils.NewResources(hdCfg.Network.Value)

	// EC Manager
	ecManager, err := NewExecutionClientManager(hdCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating executon client manager: %w", err)
	}

	// Beacon manager
	bcManager, err := NewBeaconClientManager(hdCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Beacon client manager: %w", err)
	}

	// Docker client
	dockerClient, err := docker.NewClientWithOpts(docker.WithVersion(config.DockerApiVersion))
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
	provider := &ServiceProvider[ConfigType]{
		moduleDir: moduleDir,
		userDir:   hdCfg.HyperdriveUserDirectory,
		hdCfg:     hdCfg,
		hdClient:  hdClient,
		ecManager: ecManager,
		bcManager: bcManager,
		docker:    dockerClient,
		resources: resources,
		txMgr:     txMgr,
		queryMgr:  queryMgr,
		apiLogger: &apiLogger,
	}
	return provider, nil
}

// ===============
// === Getters ===
// ===============

func (p *ServiceProvider[_]) GetModuleDir() string {
	return p.moduleDir
}

func (p *ServiceProvider[_]) GetUserDir() string {
	return p.userDir
}

func (p *ServiceProvider[_]) GetHyperdriveConfig() *config.HyperdriveConfig {
	return p.hdCfg
}

func (p *ServiceProvider[ConfigType]) GetModuleConfig() ConfigType {
	return p.moduleConfig
}

func (p *ServiceProvider[_]) GetHyperdriveClient() *client.ApiClient {
	return p.hdClient
}

func (p *ServiceProvider[_]) GetEthClient() *ExecutionClientManager {
	return p.ecManager
}

func (p *ServiceProvider[_]) GetBeaconClient() *BeaconClientManager {
	return p.bcManager
}

func (p *ServiceProvider[_]) GetDocker() *docker.Client {
	return p.docker
}

func (p *ServiceProvider[_]) GetResources() *utils.Resources {
	return p.resources
}

func (p *ServiceProvider[_]) GetTransactionManager() *eth.TransactionManager {
	return p.txMgr
}

func (p *ServiceProvider[_]) GetQueryManager() *eth.QueryManager {
	return p.queryMgr
}

func (p *ServiceProvider[_]) GetApiLogger() *log.ColorLogger {
	return p.apiLogger
}

func (p *ServiceProvider[_]) IsDebugMode() bool {
	return p.hdCfg.DebugMode.Value
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

package services

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/fatih/color"
	"github.com/nodeset-org/hyperdrive/client"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/rocket-pool/node-manager-core/node/services"
)

const (
	apiLogColor color.Attribute = color.FgHiCyan
)

// A container for all of the various services used by Hyperdrive
type ServiceProvider[ConfigType hdconfig.IModuleConfig] struct {
	*services.ServiceProvider

	// Services
	hdCfg        *hdconfig.HyperdriveConfig
	moduleConfig ConfigType
	hdClient     *client.ApiClient
	resources    *config.NetworkResources
	signer       *ModuleSigner

	// Path info
	moduleDir string
	userDir   string
}

// Creates a new ServiceProvider instance
func NewServiceProvider[ConfigType hdconfig.IModuleConfig](moduleDir string, factory func(*hdconfig.HyperdriveConfig) ConfigType) (*ServiceProvider[ConfigType], error) {
	// Create a client for the Hyperdrive daemon
	hyperdriveSocket := filepath.Join(moduleDir, hdconfig.HyperdriveSocketFilename)
	hdClient := client.NewApiClient(hdconfig.HyperdriveDaemonRoute, hyperdriveSocket, false)

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

	// Resources
	resources := hdCfg.GetNetworkResources()

	// Signer
	signer := NewModuleSigner(hdClient)

	// Create the provider
	provider := &ServiceProvider[ConfigType]{
		moduleDir:    moduleDir,
		userDir:      hdCfg.HyperdriveUserDirectory,
		hdCfg:        hdCfg,
		moduleConfig: moduleCfg,
		hdClient:     hdClient,
		resources:    resources,
		signer:       signer,
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

func (p *ServiceProvider[_]) GetHyperdriveConfig() *hdconfig.HyperdriveConfig {
	return p.hdCfg
}

func (p *ServiceProvider[ConfigType]) GetModuleConfig() ConfigType {
	return p.moduleConfig
}

func (p *ServiceProvider[_]) GetHyperdriveClient() *client.ApiClient {
	return p.hdClient
}

func (p *ServiceProvider[_]) GetResources() *config.NetworkResources {
	return p.resources
}

func (p *ServiceProvider[_]) GetSigner() *ModuleSigner {
	return p.signer
}

func (p *ServiceProvider[_]) IsDebugMode() bool {
	return p.hdCfg.DebugMode.Value
}

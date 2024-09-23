package client

import (
	"fmt"
	"reflect"
	"time"

	csconfig "github.com/nodeset-org/hyperdrive-constellation/shared/config"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	swconfig "github.com/nodeset-org/hyperdrive-stakewise/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	externalIPTimeout time.Duration = 3 * time.Second
)

// Wrapper for global configuration
type GlobalConfig struct {
	ExternalIP string

	// Hyperdrive
	Hyperdrive          *hdconfig.HyperdriveConfig
	HyperdriveResources *hdconfig.MergedResources

	// StakeWise
	StakeWise          *swconfig.StakeWiseConfig
	StakeWiseResources *swconfig.StakeWiseResources

	// Constellation
	Constellation          *csconfig.ConstellationConfig
	ConstellationResources *csconfig.ConstellationResources
}

// Make a new global config
func NewGlobalConfig(hdCfg *hdconfig.HyperdriveConfig, hdSettings []*hdconfig.HyperdriveSettings, swCfg *swconfig.StakeWiseConfig, swSettings []*swconfig.StakeWiseSettings, csCfg *csconfig.ConstellationConfig, csSettings []*csconfig.ConstellationSettings) (*GlobalConfig, error) {
	// Make the config
	cfg := &GlobalConfig{
		Hyperdrive:    hdCfg,
		StakeWise:     swCfg,
		Constellation: csCfg,
	}

	// Get the HD resources
	network := hdCfg.Network.Value
	for _, setting := range hdSettings {
		if setting.Key == network {
			cfg.HyperdriveResources = &hdconfig.MergedResources{
				NetworkResources:    setting.NetworkResources,
				HyperdriveResources: setting.HyperdriveResources,
			}
			break
		}
	}
	if cfg.HyperdriveResources == nil {
		return nil, fmt.Errorf("could not find hyperdrive resources for network [%s]", network)
	}

	// Get the StakeWise resources
	for _, setting := range swSettings {
		if setting.Key == network {
			cfg.StakeWiseResources = setting.StakeWiseResources
			break
		}
	}
	if cfg.StakeWiseResources == nil {
		return nil, fmt.Errorf("could not find stakewise resources for network [%s]", network)
	}

	// Get the Constellation resources
	for _, setting := range csSettings {
		if setting.Key == network {
			cfg.ConstellationResources = setting.ConstellationResources
			break
		}
	}
	if cfg.ConstellationResources == nil {
		return nil, fmt.Errorf("could not find constellation resources for network [%s]", network)
	}

	/*
		for _, module := range cfg.GetAllModuleConfigs() {
			config.ApplyDefaults(module, hdCfg.Network.Value)
		}
	*/
	return cfg, nil
}

// Get the configs for all of the modules in the system
func (c *GlobalConfig) GetAllModuleConfigs() []hdconfig.IModuleConfig {
	return []hdconfig.IModuleConfig{
		c.StakeWise,
		c.Constellation,
	}
}

// Serialize the config and all modules
func (c *GlobalConfig) Serialize() map[string]any {
	return c.Hyperdrive.Serialize(c.GetAllModuleConfigs(), false)
}

// Deserialize the config's modules (assumes the Hyperdrive config itself has already been deserialized)
func (c *GlobalConfig) DeserializeModules() error {
	// Load Stakewise
	stakewiseName := c.StakeWise.GetModuleName()
	section, exists := c.Hyperdrive.Modules[stakewiseName]
	if exists {
		configMap, ok := section.(map[string]any)
		if !ok {
			return fmt.Errorf("config module section [%s] is not a map, it's a %s", stakewiseName, reflect.TypeOf(section))
		}
		err := c.StakeWise.Deserialize(configMap, c.Hyperdrive.Network.Value)
		if err != nil {
			return fmt.Errorf("error deserializing stakewise configuration: %w", err)
		}
	}

	// Load Constellation
	constellationName := c.Constellation.GetModuleName()
	section, exists = c.Hyperdrive.Modules[constellationName]
	if exists {
		configMap, ok := section.(map[string]any)
		if !ok {
			return fmt.Errorf("config module section [%s] is not a map, it's a %s", constellationName, reflect.TypeOf(section))
		}
		err := c.Constellation.Deserialize(configMap, c.Hyperdrive.Network.Value)
		if err != nil {
			return fmt.Errorf("error deserializing constellation configuration: %w", err)
		}
	}
	return nil
}

// Creates a copy of the configuration
func (c *GlobalConfig) CreateCopy() *GlobalConfig {
	hdCopy := c.Hyperdrive.Clone()
	swCopy := c.StakeWise.Clone().(*swconfig.StakeWiseConfig)
	csCopy := c.Constellation.Clone().(*csconfig.ConstellationConfig)

	return &GlobalConfig{
		Hyperdrive:    hdCopy,
		StakeWise:     swCopy,
		Constellation: csCopy,
	}
}

// Changes the current network, propagating new parameter settings if they are affected
func (c *GlobalConfig) ChangeNetwork(newNetwork config.Network) {
	// Get the current network
	oldNetwork := c.Hyperdrive.Network.Value
	if oldNetwork == newNetwork {
		return
	}

	// Run the changes
	c.Hyperdrive.ChangeNetwork(newNetwork)
	for _, module := range c.GetAllModuleConfigs() {
		module.ChangeNetwork(oldNetwork, newNetwork)
	}
}

// Updates the default parameters based on the current network value
func (c *GlobalConfig) UpdateDefaults() {
	network := c.Hyperdrive.Network.Value
	config.UpdateDefaults(c.Hyperdrive, network)
	for _, module := range c.GetAllModuleConfigs() {
		module.UpdateDefaults(network)
	}
}

// Checks to see if the current configuration is valid; if not, returns a list of errors
func (c *GlobalConfig) Validate() []string {
	errors := []string{}

	// Check for illegal blank strings
	/* TODO - this needs to be smarter and ignore irrelevant settings
	for _, param := range config.GetParameters() {
		if param.Type == ParameterType_String && !param.CanBeBlank && param.Value == "" {
			errors = append(errors, fmt.Sprintf("[%s] cannot be blank.", param.Name))
		}
	}

	for name, subconfig := range config.GetSubconfigs() {
		for _, param := range subconfig.GetParameters() {
			if param.Type == ParameterType_String && !param.CanBeBlank && param.Value == "" {
				errors = append(errors, fmt.Sprintf("[%s - %s] cannot be blank.", name, param.Name))
			}
		}
	}
	*/

	// Ensure there's a MEV-boost URL
	if c.Hyperdrive.MevBoost.Enable.Value {
		switch c.Hyperdrive.MevBoost.Mode.Value {
		case config.ClientMode_Local:
			// In local MEV-boost mode, the user has to have at least one relay
			relays := c.Hyperdrive.MevBoost.GetEnabledMevRelays()
			if len(relays) == 0 {
				errors = append(errors, "You have MEV-boost enabled in local mode but don't have any profiles or relays enabled. Please select at least one profile or relay to use MEV-boost.")
			}
		case config.ClientMode_External:
			// In external MEV-boost mode, the user has to have an external URL if they're running Docker mode
			if c.Hyperdrive.IsLocalMode() && c.Hyperdrive.MevBoost.ExternalUrl.Value == "" {
				errors = append(errors, "You have MEV-boost enabled in external mode but don't have a URL set. Please enter the external MEV-boost server URL to use it.")
			}
		default:
			errors = append(errors, "You do not have a MEV-Boost mode configured. You must either select a mode in the `hyperdrive service config` UI, or disable MEV-Boost.")
		}
	}

	// Ensure the selected port numbers are unique. Keeps track of all the errors
	portMap := make(map[uint16]bool)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.ApiPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.StakeWise.ApiPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Constellation.ApiPort, errors)
	if c.Hyperdrive.ClientMode.Value == config.ClientMode_Local {
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalExecutionClient.HttpPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalExecutionClient.WebsocketPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalExecutionClient.EnginePort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalExecutionClient.P2pPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalBeaconClient.HttpPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalBeaconClient.P2pPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalBeaconClient.Lighthouse.P2pQuicPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalBeaconClient.Prysm.RpcPort, errors)
	}
	if c.Hyperdrive.Metrics.EnableMetrics.Value {
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.EcMetricsPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.BnMetricsPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.Prometheus.Port, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.ExporterMetricsPort, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.Grafana.Port, errors)
		portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.DaemonMetricsPort, errors)
		if c.StakeWise.Enabled.Value {
			portMap, errors = addAndCheckForDuplicate(portMap, c.StakeWise.VcCommon.MetricsPort, errors)
		}
		if c.Constellation.Enabled.Value {
			portMap, errors = addAndCheckForDuplicate(portMap, c.Constellation.VcCommon.MetricsPort, errors)
		}
	}
	if c.Hyperdrive.MevBoost.Enable.Value && c.Hyperdrive.MevBoost.Mode.Value == config.ClientMode_Local {
		_, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.MevBoost.Port, errors)
	}

	return errors
}

// Get all of the settings that have changed between an old config and this config, and get all of the containers that are affected by those changes - also returns whether or not the selected network was changed
func (c *GlobalConfig) GetChanges(oldConfig *GlobalConfig) ([]*config.ChangedSection, map[config.ContainerID]bool, bool) {
	sectionList := []*config.ChangedSection{}
	changedContainers := map[config.ContainerID]bool{}

	// Process all configs for changes
	sectionList = getChanges(oldConfig.Hyperdrive, c.Hyperdrive, sectionList, changedContainers)
	sectionList = getChanges(oldConfig.StakeWise, c.StakeWise, sectionList, changedContainers)
	sectionList = getChanges(oldConfig.Constellation, c.Constellation, sectionList, changedContainers)

	// Add all VCs to the list of changed containers if any change requires a VC change
	if changedContainers[config.ContainerID_ValidatorClient] {
		delete(changedContainers, config.ContainerID_ValidatorClient)
		for _, module := range c.GetAllModuleConfigs() {
			vcInfo := module.GetValidatorContainerTagInfo()
			for name := range vcInfo {
				changedContainers[name] = true
			}
		}
	}

	// Check if the network has changed
	changeNetworks := false
	if oldConfig.Hyperdrive.Network.Value != c.Hyperdrive.Network.Value {
		changeNetworks = true
	}

	return sectionList, changedContainers, changeNetworks
}

// Attempts to load the system's external IP address
func (c *GlobalConfig) LoadExternalIP() {
	if c.ExternalIP != "" {
		return
	}

	// Get the external IP address
	ip, err := config.GetExternalIP(externalIPTimeout)
	if err != nil {
		fmt.Println("Warning: couldn't get external IP address; if you're using Nimbus, Besu, or Teku, it may have trouble finding peers:")
		fmt.Println(err.Error())
		return
	}

	if ip.To4() == nil {
		fmt.Println("Warning: external IP address is v6; if you're using Nimbus, Besu, or Teku, it may have trouble finding peers.")
	}

	c.ExternalIP = ip.String()
}

// Compare two config sections and see what's changed between them, generating a ChangedSection for the results.
func getChanges(
	oldConfig config.IConfigSection,
	newConfig config.IConfigSection,
	sectionList []*config.ChangedSection,
	changedContainers map[config.ContainerID]bool,
) []*config.ChangedSection {
	section, changeCount := config.GetChangedSettings(oldConfig, newConfig)
	section.Name = newConfig.GetTitle()
	if changeCount > 0 {
		config.GetAffectedContainers(section, changedContainers)
		sectionList = append(sectionList, section)
	}
	return sectionList
}

// Check a port setting to see if it's already used elsewhere
func addAndCheckForDuplicate(portMap map[uint16]bool, param config.Parameter[uint16], errors []string) (map[uint16]bool, []string) {
	port := param.Value
	if portMap[port] {
		return portMap, append(errors, fmt.Sprintf("Port %d for %s is already in use", port, param.GetCommon().Name))
	} else {
		portMap[port] = true
	}
	return portMap, errors
}

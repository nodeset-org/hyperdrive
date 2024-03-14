package client

import (
	"fmt"
	"reflect"

	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
)

// Wrapper for global configuration
type GlobalConfig struct {
	Hyperdrive *hdconfig.HyperdriveConfig
	Stakewise  *swconfig.StakewiseConfig
}

// Make a new global config
func NewGlobalConfig(hdCfg *hdconfig.HyperdriveConfig) *GlobalConfig {
	cfg := &GlobalConfig{
		Hyperdrive: hdCfg,
		Stakewise:  swconfig.NewStakewiseConfig(hdCfg),
	}

	for _, module := range cfg.GetAllModuleConfigs() {
		config.ApplyDefaults(module, hdCfg.Network.Value)
	}
	return cfg
}

// Get the configs for all of the modules in the system
func (c *GlobalConfig) GetAllModuleConfigs() []hdconfig.IModuleConfig {
	return []hdconfig.IModuleConfig{
		c.Stakewise,
	}
}

// Serialize the config and all modules
func (c *GlobalConfig) Serialize() map[string]any {
	return c.Hyperdrive.Serialize(c.GetAllModuleConfigs())
}

// Deserialize the config's modules (assumes the Hyperdrive config itself has already been deserialized)
func (c *GlobalConfig) DeserializeModules() error {
	// Load Stakewise
	stakewiseName := c.Stakewise.GetModuleName()
	section, exists := c.Hyperdrive.Modules[stakewiseName]
	if exists {
		configMap, ok := section.(map[string]any)
		if !ok {
			return fmt.Errorf("config module section [%s] is not a map, it's a %s", stakewiseName, reflect.TypeOf(section))
		}
		err := config.Deserialize(c.Stakewise, configMap, c.Hyperdrive.Network.Value)
		if err != nil {
			return fmt.Errorf("error deserializing stakewise configuration: %w", err)
		}
	}
	return nil
}

// Creates a copy of the configuration
func (c *GlobalConfig) CreateCopy() *GlobalConfig {
	// Hyperdrive
	network := c.Hyperdrive.Network.Value
	hdCopy := hdconfig.NewHyperdriveConfig(c.Hyperdrive.HyperdriveUserDirectory)
	config.Clone(c.Hyperdrive, hdCopy, network)

	// Stakewise
	swCopy := swconfig.NewStakewiseConfig(hdCopy)
	config.Clone(c.Stakewise, swCopy, network)

	return &GlobalConfig{
		Hyperdrive: hdCopy,
		Stakewise:  swCopy,
	}
}

// Changes the current network, propagating new parameter settings if they are affected
func (c *GlobalConfig) ChangeNetwork(newNetwork config.Network) {
	// Get the current network
	oldNetwork := c.Hyperdrive.Network.Value
	if oldNetwork == newNetwork {
		return
	}
	c.Hyperdrive.Network.Value = newNetwork

	// Run the changes
	c.Hyperdrive.ChangeNetwork(newNetwork)
	for _, module := range c.GetAllModuleConfigs() {
		config.ChangeNetwork(module, oldNetwork, newNetwork)
	}
}

// Updates the default parameters based on the current network value
func (c *GlobalConfig) UpdateDefaults() {
	network := c.Hyperdrive.Network.Value
	config.UpdateDefaults(c.Hyperdrive, network)
	for _, module := range c.GetAllModuleConfigs() {
		config.UpdateDefaults(module, network)
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

	// Ensure the selected port numbers are unique. Keeps track of all the errors
	portMap := make(map[uint16]bool)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalBeaconConfig.HttpPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalBeaconConfig.P2pPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalExecutionConfig.HttpPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalExecutionConfig.WebsocketPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalExecutionConfig.EnginePort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalExecutionConfig.P2pPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.EcMetricsPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.BnMetricsPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.Prometheus.Port, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.ExporterMetricsPort, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.Grafana.Port, errors)
	portMap, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.Metrics.DaemonMetricsPort, errors)
	_, errors = addAndCheckForDuplicate(portMap, c.Hyperdrive.LocalBeaconConfig.Lighthouse.P2pQuicPort, errors)

	return errors
}

// Get all of the settings that have changed between an old config and this config, and get all of the containers that are affected by those changes - also returns whether or not the selected network was changed
func (c *GlobalConfig) GetChanges(oldConfig *GlobalConfig) ([]*config.ChangedSection, map[config.ContainerID]bool, bool) {
	sectionList := []*config.ChangedSection{}
	changedContainers := map[config.ContainerID]bool{}

	// Process all configs for changes
	sectionList = getChanges(oldConfig.Hyperdrive, c.Hyperdrive, sectionList, changedContainers)
	sectionList = getChanges(oldConfig.Stakewise, c.Stakewise, sectionList, changedContainers)

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

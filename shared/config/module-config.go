package config

import (
	"github.com/rocket-pool/node-manager-core/config"
)

const (
	ModulesName         string = "modules"
	ValidatorsDirectory string = "validators"
)

type IModuleConfig interface {
	config.IConfigSection

	// Get the name of the module
	GetModuleName() string

	// Get the short name of the module, for things like prefixing
	GetShortName() string

	// The name to use for the Hyperdrive Client log file
	GetHdClientLogFileName() string

	// The name to use for the API log file
	GetApiLogFileName() string

	// The name to use for the tasks log file
	GetTasksLogFileName() string

	// Get the list of all log file names used by the module
	GetLogNames() []string

	// A map of the Validator Client IDs to their container tags
	GetValidatorContainerTagInfo() map[config.ContainerID]string

	// Return if doppelganger detection is enabled for any of the VCs
	IsDoppelgangerEnabled() bool

	// True if the module is enabled
	IsEnabled() bool

	// Get the list of containers that should be deployed
	GetContainersToDeploy() []config.ContainerID

	// Handle a network change, for example by updating the default parameters
	ChangeNetwork(oldNetwork config.Network, newNetwork config.Network)

	// Create a copy of this config
	Clone() IModuleConfig

	// Update the default parameters based on the Hyperdrive config's current network value
	UpdateDefaults(network config.Network)

	// Verify the current settings and publish a list of errors that must be resolved before saving
	Validate() []string

	// Serialize the module config to a map
	Serialize() map[string]any

	// Deserialize the module config from a map
	Deserialize(configMap map[string]any, network config.Network) error

	// Get the version of the module config
	GetVersion() string
}

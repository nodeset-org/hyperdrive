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

	// The name to use for the API log file
	GetApiLogFileName() string

	// The name to use for the tasks log file
	GetTasksLogFileName() string

	// A map of the Validator Client IDs to their container tags
	GetValidatorContainerTagInfo() map[config.ContainerID]string

	// Return if doppelganger detection is enabled for any of the VCs
	IsDoppelgangerEnabled() bool

	// True if the module is enabled
	IsEnabled() bool

	// Get the list of containers that should be deployed
	GetContainersToDeploy() []config.ContainerID
}

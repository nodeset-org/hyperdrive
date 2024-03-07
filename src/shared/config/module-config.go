package config

import (
	nmc_config "github.com/rocket-pool/node-manager-core/config"
)

const (
	ModulesName         string = "modules"
	ValidatorsDirectory string = "validators"
)

type IModuleConfig interface {
	nmc_config.IConfigSection

	GetModuleName() string

	// A map of the Validator Client IDs to their container tags
	GetValidatorContainerTagInfo() map[nmc_config.ContainerID]string

	// Return if doppelganger detection is enabled for any of the VCs
	IsDoppelgangerEnabled() bool

	// True if the module is enabled
	IsEnabled() bool

	// Get the list of containers that should be deployed
	GetContainersToDeploy() []nmc_config.ContainerID
}

package modconfig

import "github.com/nodeset-org/hyperdrive/shared/types"

const (
	ModulesName         string = "modules"
	ValidatorsDirectory string = "validators"
)

type IModuleConfig interface {
	types.IConfigSection

	GetModuleName() string

	// A map of the Validator Client suffixes to their container tags
	GetValidatorContainerTagInfo() map[string]string

	// Return if doppelganger detection is enabled for any of the VCs
	IsDoppelgangerEnabled() bool

	// True if the module is enabled
	IsEnabled() bool
}

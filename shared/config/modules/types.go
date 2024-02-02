package modconfig

import "github.com/nodeset-org/hyperdrive/shared/types"

type IModuleConfig interface {
	types.IConfigSection

	GetModuleName() string

	// A map of the Validator Client suffixes to their container tags
	GetValidatorContainerTagInfo() map[string]string
}
